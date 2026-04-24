package service

import (
	"Authentication_Service/internal/config"
	"Authentication_Service/internal/constant"
	"Authentication_Service/internal/dto/request"
	response "Authentication_Service/internal/dto/response"
	"Authentication_Service/internal/mapper"
	"Authentication_Service/internal/model"
	_interface "Authentication_Service/internal/repository/interface"
	serviceInterface "Authentication_Service/internal/service/interface"
	"Authentication_Service/internal/utils"
	"Authentication_Service/pkg/email"
	model2 "Authentication_Service/pkg/email/model"
	"Authentication_Service/pkg/google"
	"Authentication_Service/pkg/kafka"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthenticationService struct {
	AuthenticationRepository _interface.IAuthenticationRepository
	UserRepository           _interface.IUserRepository
	OtpRepository            _interface.IOtpRepository
	TokenGenerator           serviceInterface.ITokenGenerator
	RefreshTokenRepository   _interface.IRefreshTokenRepository
	RegisterEventProducer    kafka.RegisterEventProducer
	GoogleOAuth2             google.OAuth2Client
}

func (authentication *AuthenticationService) Logout(ctx context.Context, req *request.LogoutReqDto) (*response.LogoutResDto, *model.BusinessError) {
	tokenHash := utils.HashRefreshToken(req.RefreshToken)
	refreshToken, errGet := authentication.RefreshTokenRepository.GetByTokenHash(tokenHash)

	if errGet != nil {
		if errors.Is(errGet, gorm.ErrRecordNotFound) {
			return nil, config.SessionNotFound
		}
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to find session")
	}

	if refreshToken == nil {
		return nil, config.SessionNotFound
	}

	// Soft revoke in DB (audit trail; matches is_revoked in refresh_tokens)
	refreshToken.IsRevoked = true
	if errUpdate := authentication.RefreshTokenRepository.Update(refreshToken); errUpdate != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to revoke session")
	}

	// Blacklist access token in Redis until expiry so it cannot be reused
	if req.AccessToken != nil && *req.AccessToken != "" {
		expAt, errExp := authentication.TokenGenerator.GetExpiresAt(*req.AccessToken)
		if errExp == nil {
			ttl := time.Until(expAt)
			if ttl <= 0 {
				ttl = time.Second
			}
			if errBl := authentication.AuthenticationRepository.AddAccessTokenToBlacklist(ctx, *req.AccessToken, ttl); errBl != nil {
				logger.Warn("logout: failed to blacklist access token: %v", errBl)
			}
		}
	}

	return &response.LogoutResDto{
		IsLogout: true,
	}, nil
}

func (authentication *AuthenticationService) ValidateMail(ctx context.Context, validateMailReq *request.ValidateMailReqDto) (*response.ValidateMailResDto, *model.BusinessError) {
	otp, errGet := authentication.OtpRepository.GetOtp(ctx, validateMailReq.Email)

	if errGet != nil {
		return nil, config.GetOtpFailedError
	}

	if otp != validateMailReq.Otp {
		return nil, config.InvalidOtpError
	}

	isValid, errCheckTime := authentication.OtpRepository.CheckTimeLife(ctx, validateMailReq.Email)

	if errCheckTime != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to check OTP time life")
	}

	if !isValid {
		return nil, model.NewBusinessError(http.StatusBadRequest, "OTP has expired")
	}

	user, errFindUser := authentication.UserRepository.FindByEmail(validateMailReq.Email)

	if errFindUser != nil {
		if errors.Is(errFindUser, gorm.ErrRecordNotFound) {
			return nil, config.UserNotFound
		}
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to find user")
	}

	user.IsActive = constant.ACTIVE

	_, errUpdate := authentication.UserRepository.Update(user)

	if errUpdate != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to activate user")
	}

	errorDel := authentication.OtpRepository.DeleteOtp(ctx, validateMailReq.Email)

	if errorDel != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to delete OTP")
	}

	return &response.ValidateMailResDto{
		Valid: true,
	}, nil
}

func (authentication *AuthenticationService) CheckEmail(email string) (bool, *model.User) {
	res, err := authentication.AuthenticationRepository.FindByEmail(email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if res != nil {
		return true, res
	}

	return false, nil
}

func (authentication *AuthenticationService) Login(c context.Context, user *request.LoginReqDto) (*response.LoginResDto, *model.BusinessError) {
	userEmail, err := authentication.AuthenticationRepository.FindByEmail(user.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, config.UserNotFound
		}
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to find user")
	}

	if userEmail == nil {
		return nil, config.UserNotFound
	}

	if userEmail.IsActive == constant.NO_ACTIVE {
		return nil, config.UserNotActivated
	}

	if !utils.CheckPasswordHash(user.Password, userEmail.Password) {
		return nil, config.InvalidCredentials
	}

	userResDto := mapper.UserFromUserResDto(userEmail)
	if userResDto == nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to map user data")
	}

	token, err := authentication.TokenGenerator.GenerateToken(*userResDto, constant.ACCESS_TOKEN_EXPIRE)
	if err != nil {
		logger.Error("failed to generate token for user %s: %v", user.Email, err)
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to generate token")
	}

	refreshToken, errRefreshToken := authentication.TokenGenerator.GenerateToken(*userResDto, constant.REFRESH_TOKEN_EXPIRE)

	if errRefreshToken != nil {
		logger.Error("failed to generate refresh token for user %s: %v", user.Email, errRefreshToken)
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to generate refresh token")
	}

	ipAddress := ""
	if user.IpAddress != nil {
		ipAddress = *user.IpAddress
	}
	deviceInfo := ""
	if user.DeviceInfo != nil {
		deviceInfo = *user.DeviceInfo
	}
	expiresAt := time.Now().Add(time.Duration(constant.REFRESH_TOKEN_EXPIRE) * time.Minute)

	errSave := authentication.AuthenticationRepository.SaveRefreshToken(c, &model.RefreshToken{
		UserId:     userEmail.Id,
		TokenHash:  utils.HashRefreshToken(refreshToken),
		IpAddress:  ipAddress,
		DeviceInfo: deviceInfo,
		ExpiresAt:  expiresAt,
	})
	if errSave != nil {
		logger.Error("failed to save refresh token for user %s: %v", user.Email, errSave)
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to save session")
	}

	// Log refresh token for Postman/testing (remove or guard in production)
	logger.Info("[Postman] refresh_token for %s: %s", user.Email, refreshToken)

	return &response.LoginResDto{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (authentication *AuthenticationService) Register(c context.Context, user *request.RegisterReqDto, email email.IEmailStrategy) (*response.RegisterResDto, *model.BusinessError) {
	currentUser, businessErr := authentication.validateUserForRegistration(c, user.Email)
	if businessErr != nil {
		return nil, businessErr
	}

	otp, businessErr := authentication.generateAndSaveOtp(c, user.Email)
	if businessErr != nil {
		return nil, businessErr
	}

	businessErr = authentication.enqueueOtpEmail(c, user, otp, email)
	if businessErr != nil {
		return nil, businessErr
	}

	hashedPassword := utils.HashPassword(user.Password)
	userAfterSave, businessErr := authentication.createOrUpdateUser(c, user, hashedPassword, currentUser)
	if businessErr != nil {
		return nil, businessErr
	}

	return mapper.FromUserModelToRegisterRes(userAfterSave), nil
}

func (authentication *AuthenticationService) enqueueOtpEmail(ctx context.Context, user *request.RegisterReqDto, otp string, emailStrategy email.IEmailStrategy) *model.BusinessError {
	if authentication.RegisterEventProducer == nil {
		return authentication.sendOtpEmail(user, otp, emailStrategy)
	}

	event := kafka.RegisterOTPEvent{
		Email:       user.Email,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Otp:         otp,
		RequestedAt: time.Now().UTC(),
	}

	if err := authentication.RegisterEventProducer.PublishRegisterOTP(ctx, event); err != nil {
		logger.Warn("register: failed to publish OTP event for %s, fallback to sync email: %v", user.Email, err)
		return authentication.sendOtpEmail(user, otp, emailStrategy)
	}

	return nil
}

func (authentication *AuthenticationService) validateUserForRegistration(ctx context.Context, email string) (*model.User, *model.BusinessError) {
	currentUser, err := authentication.UserRepository.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to check user existence")
	}

	if currentUser == nil {
		return nil, nil
	}

	if currentUser.IsActive == constant.ACTIVE {
		return nil, config.UserAlreadyExists
	}

	if currentUser.IsActive == constant.NO_ACTIVE {
		businessErr := authentication.handleExistingOtp(ctx, email)
		if businessErr != nil {
			return nil, businessErr
		}
		return currentUser, nil
	}

	return nil, model.NewBusinessError(http.StatusInternalServerError, "invalid user status")
}

func (authentication *AuthenticationService) handleExistingOtp(ctx context.Context, email string) *model.BusinessError {
	otp, errGetOtp := authentication.OtpRepository.GetOtp(ctx, email)
	if errGetOtp != nil || otp == "" {
		return nil
	}

	isValid, errCheckTime := authentication.OtpRepository.CheckTimeLife(ctx, email)
	if errCheckTime != nil {
		return model.NewBusinessError(http.StatusInternalServerError, "failed to check OTP time life")
	}

	if isValid {
		return model.NewBusinessError(http.StatusBadRequest, "OTP has been sent to your email. Please check and validate it.")
	}

	if errDel := authentication.OtpRepository.DeleteOtp(ctx, email); errDel != nil {
		return model.NewBusinessError(http.StatusInternalServerError, "failed to delete expired OTP")
	}

	return nil
}

func (authentication *AuthenticationService) generateAndSaveOtp(ctx context.Context, email string) (string, *model.BusinessError) {
	otp, errOtp := utils.GenerateRandomOtp(6)
	if errOtp != nil {
		return "", config.GenOtpFailed
	}

	if errSave := authentication.OtpRepository.SaveOtp(ctx, email, otp, 5*time.Minute); errSave != nil {
		return "", config.SaveUserRedisError
	}

	return otp, nil
}

func (authentication *AuthenticationService) sendOtpEmail(user *request.RegisterReqDto, otp string, email email.IEmailStrategy) *model.BusinessError {
	to := model2.To{
		Email: user.Email,
		Name:  fmt.Sprintf("%s %s", utils.StringValue(user.Lastname), utils.StringValue(user.Firstname)),
	}

	if err := email.SendEmail(to, "Your OTP Code", "Your OTP code is: "+otp+". It is valid for 5 minutes."); err != nil {
		logger.Errorf("failed to send OTP email: %v", err)
		return model.NewBusinessError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (authentication *AuthenticationService) createOrUpdateUser(ctx context.Context, user *request.RegisterReqDto, hashedPassword string, currentUser *model.User) (*model.User, *model.BusinessError) {
	if currentUser != nil {
		return authentication.updateExistingUser(ctx, user, hashedPassword, currentUser)
	}

	return authentication.createNewUser(ctx, user, hashedPassword)
}

func (authentication *AuthenticationService) updateExistingUser(ctx context.Context, user *request.RegisterReqDto, hashedPassword string, currentUser *model.User) (*model.User, *model.BusinessError) {
	existingUser, errRecheck := authentication.UserRepository.FindByEmail(user.Email)
	if errRecheck != nil {
		if errors.Is(errRecheck, gorm.ErrRecordNotFound) {
			return nil, model.NewBusinessError(http.StatusInternalServerError, "user was deleted during registration")
		}
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to recheck user before update")
	}

	if existingUser == nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "user was deleted during registration")
	}

	if existingUser.IsActive == constant.ACTIVE {
		return nil, config.UserAlreadyExists
	}

	if existingUser.IsActive != constant.NO_ACTIVE {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "invalid user status for update")
	}

	existingUser.Password = hashedPassword
	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.IsActive = constant.NO_ACTIVE

	updatedUser, errUpdate := authentication.UserRepository.Update(existingUser)
	if errUpdate != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, errUpdate.Error())
	}

	return updatedUser, nil
}

func (authentication *AuthenticationService) createNewUser(ctx context.Context, user *request.RegisterReqDto, hashedPassword string) (*model.User, *model.BusinessError) {
	userModel := &model.User{
		Id:        uuid.New().String(),
		Email:     user.Email,
		Password:  hashedPassword,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		IsActive:  constant.NO_ACTIVE,
	}

	_, errCreate := authentication.UserRepository.Create(userModel)
	if errCreate == nil {
		return userModel, nil
	}

	if !errors.Is(errCreate, gorm.ErrDuplicatedKey) {
		return nil, model.NewBusinessError(http.StatusInternalServerError, errCreate.Error())
	}

	return authentication.handleDuplicateKeyError(ctx, user, hashedPassword)
}

func (authentication *AuthenticationService) ResendOtp(ctx context.Context, req *request.ResendOtpReqDto, emailStrategy email.IEmailStrategy) (*response.ResendOtpResDto, *model.BusinessError) {
	user, err := authentication.UserRepository.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to find user")
	}
	if user == nil {
		return nil, config.UserNotFound
	}
	if user.IsActive == constant.ACTIVE {
		return nil, config.UserAlreadyExists
	}
	// Only allow resend for pending verification (NO_ACTIVE)
	otp, businessErr := authentication.generateAndSaveOtp(ctx, req.Email)
	if businessErr != nil {
		return nil, businessErr
	}
	to := model2.To{Email: req.Email, Name: req.Email}
	if err := emailStrategy.SendEmail(to, "Your OTP Code", "Your OTP code is: "+otp+". It is valid for 5 minutes."); err != nil {
		logger.Errorf("failed to send OTP email: %v", err)
		return nil, model.NewBusinessError(http.StatusInternalServerError, err.Error())
	}
	return &response.ResendOtpResDto{IsResent: true}, nil
}

func (authentication *AuthenticationService) ForgotPassword(ctx context.Context, req *request.ForgotPasswordReqDto, emailStrategy email.IEmailStrategy) (*response.ForgotPasswordResDto, *model.BusinessError) {
	user, err := authentication.UserRepository.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to find user")
	}
	if user == nil {
		return nil, config.UserNotFound
	}
	if user.IsActive != constant.ACTIVE {
		return nil, config.UserNotActivated
	}
	otp, errOtp := utils.GenerateRandomOtp(6)
	if errOtp != nil {
		return nil, config.GenOtpFailed
	}
	if errSave := authentication.OtpRepository.SaveForgotPasswordOtp(ctx, req.Email, otp, 10*time.Minute); errSave != nil {
		return nil, config.SaveUserRedisError
	}
	to := model2.To{Email: req.Email, Name: req.Email}
	if err := emailStrategy.SendEmail(to, "Password Reset OTP", "Your password reset OTP is: "+otp+". It is valid for 10 minutes."); err != nil {
		logger.Errorf("failed to send forgot password email: %v", err)
		return nil, model.NewBusinessError(http.StatusInternalServerError, err.Error())
	}
	return &response.ForgotPasswordResDto{Sent: true}, nil
}

func (authentication *AuthenticationService) ResetPassword(ctx context.Context, req *request.ResetPasswordReqDto) (*response.ResetPasswordResDto, *model.BusinessError) {
	storedOtp, errGet := authentication.OtpRepository.GetForgotPasswordOtp(ctx, req.Email)
	if errGet != nil {
		return nil, config.GetOtpFailedError
	}
	if storedOtp != req.Otp {
		return nil, config.InvalidOtpError
	}
	valid, errCheck := authentication.OtpRepository.CheckForgotPasswordOtpTimeLife(ctx, req.Email)
	if errCheck != nil || !valid {
		return nil, model.NewBusinessError(http.StatusBadRequest, "OTP has expired")
	}
	user, errFind := authentication.UserRepository.FindByEmail(req.Email)
	if errFind != nil || user == nil {
		return nil, config.UserNotFound
	}
	user.Password = utils.HashPassword(req.NewPassword)
	if _, errUpdate := authentication.UserRepository.Update(user); errUpdate != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to update password")
	}
	_ = authentication.OtpRepository.DeleteForgotPasswordOtp(ctx, req.Email)
	return &response.ResetPasswordResDto{Success: true}, nil
}

func (authentication *AuthenticationService) ChangePassword(ctx context.Context, userID string, req *request.ChangePasswordReqDto) (*response.ChangePasswordResDto, *model.BusinessError) {
	user, err := authentication.UserRepository.GetById(userID)
	if err != nil || user == nil {
		return nil, config.UserNotFound
	}
	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		return nil, config.InvalidOldPassword
	}
	user.Password = utils.HashPassword(req.NewPassword)
	if _, errUpdate := authentication.UserRepository.Update(user); errUpdate != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to update password")
	}
	return &response.ChangePasswordResDto{Success: true}, nil
}

func (authentication *AuthenticationService) handleDuplicateKeyError(ctx context.Context, user *request.RegisterReqDto, hashedPassword string) (*model.User, *model.BusinessError) {
	existingUser, errFind := authentication.UserRepository.FindByEmail(user.Email)
	if errFind != nil {
		if errors.Is(errFind, gorm.ErrRecordNotFound) {
			return nil, model.NewBusinessError(http.StatusInternalServerError, "user creation failed but user not found")
		}
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to find user after duplicate key error")
	}

	if existingUser == nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "user creation failed but user is nil")
	}

	if existingUser.IsActive == constant.ACTIVE {
		return nil, config.UserAlreadyExists
	}

	existingUser.Password = hashedPassword
	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.IsActive = constant.NO_ACTIVE

	updatedUser, errUpdate := authentication.UserRepository.Update(existingUser)
	if errUpdate != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, errUpdate.Error())
	}

	return updatedUser, nil
}

func (authentication *AuthenticationService) Introspect(ctx context.Context, req *request.IntrospectReqDto) (*response.IntrospectResDto, *model.BusinessError) {
	result, err := authentication.TokenGenerator.Introspect(req.Token)
	if err != nil {
		return &response.IntrospectResDto{Active: false}, nil
	}
	return &response.IntrospectResDto{
		Active:    true,
		Sub:       result.User.Id,
		Exp:       result.Exp.Unix(),
		Iat:       result.Iat.Unix(),
		Email:     result.User.Email,
		Firstname:  result.User.Firstname,
		Lastname:   result.User.Lastname,
		Avatar:     result.User.Avatar,
		CoverImage: result.User.CoverImage,
	}, nil
}

// GoogleSSO exchanges the Google OAuth2 code for user info, finds or creates the user, and returns tokens.
func (authentication *AuthenticationService) GoogleSSO(ctx context.Context, code string) (*response.LoginResDto, *model.BusinessError) {
	if authentication.GoogleOAuth2 == nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "Google SSO is not configured")
	}
	info, err := authentication.GoogleOAuth2.ExchangeCodeAndGetUserInfo(ctx, code)
	if err != nil {
		return nil, model.NewBusinessError(http.StatusBadRequest, "invalid or expired Google authorization code")
	}
	if info.Email == "" {
		return nil, model.NewBusinessError(http.StatusBadRequest, "Google did not return an email")
	}

	user, errFind := authentication.UserRepository.FindByEmail(info.Email)
	if errFind != nil && !errors.Is(errFind, gorm.ErrRecordNotFound) {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to find user")
	}

	if user == nil {
		// New user: create with Google placeholder password and ACTIVE (Google verified email).
		firstname, lastname := info.GivenName, info.FamilyName
		user = &model.User{
			Id:        uuid.New().String(),
			Email:     info.Email,
			Password:  constant.GoogleSSOPasswordPlaceholder,
			Firstname: &firstname,
			Lastname:  &lastname,
			IsActive:  constant.ACTIVE,
		}
		if info.Picture != "" {
			user.Avatar = &info.Picture
		}
		if _, errCreate := authentication.UserRepository.Create(user); errCreate != nil {
			return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to create user")
		}
	} else {
		// Existing user: ensure active so they can sign in (e.g. previously registered with email, now using Google).
		if user.IsActive == constant.NO_ACTIVE {
			user.IsActive = constant.ACTIVE
			if _, errUpdate := authentication.UserRepository.Update(user); errUpdate != nil {
				return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to activate user")
			}
		}
	}

	userResDto := mapper.UserFromUserResDto(user)
	if userResDto == nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to map user data")
	}

	accessToken, err := authentication.TokenGenerator.GenerateToken(*userResDto, constant.ACCESS_TOKEN_EXPIRE)
	if err != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to generate token")
	}
	refreshToken, errRefresh := authentication.TokenGenerator.GenerateToken(*userResDto, constant.REFRESH_TOKEN_EXPIRE)
	if errRefresh != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to generate refresh token")
	}

	expiresAt := time.Now().Add(time.Duration(constant.REFRESH_TOKEN_EXPIRE) * time.Minute)
	if errSave := authentication.AuthenticationRepository.SaveRefreshToken(ctx, &model.RefreshToken{
		UserId:     user.Id,
		TokenHash:  utils.HashRefreshToken(refreshToken),
		IpAddress:  "",
		DeviceInfo: "google_sso",
		ExpiresAt:  expiresAt,
	}); errSave != nil {
		return nil, model.NewBusinessError(http.StatusInternalServerError, "failed to save session")
	}

	return &response.LoginResDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetGoogleAuthURL returns the Google OAuth2 consent URL for the given state. Empty if SSO is not configured.
func (authentication *AuthenticationService) GetGoogleAuthURL(state string) string {
	if authentication.GoogleOAuth2 == nil {
		return ""
	}
	return authentication.GoogleOAuth2.AuthCodeURL(state)
}

func (authentication *AuthenticationService) SetSSOState(ctx context.Context, state string) error {
	return authentication.AuthenticationRepository.SetSSOState(ctx, state, 10*time.Minute)
}

func (authentication *AuthenticationService) ValidateAndConsumeSSOState(ctx context.Context, state string) (bool, error) {
	return authentication.AuthenticationRepository.ValidateAndConsumeSSOState(ctx, state)
}
