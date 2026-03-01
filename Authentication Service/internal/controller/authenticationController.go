package controller

import (
	"Authentication_Service/internal/constant"
	"Authentication_Service/internal/dto/common"
	"Authentication_Service/internal/dto/request"
	response "Authentication_Service/internal/dto/response"
	"Authentication_Service/internal/middleware"
	_interface "Authentication_Service/internal/service/interface"
	"Authentication_Service/pkg/email"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthenticationController struct {
	AuthenticationService _interface.IAuthenticationService
	Email                 email.IEmailStrategy
}

// ResendOtpController godoc
// @Summary      Resend OTP
// @Description  Resend OTP to the given email address.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  request.ResendOtpReqDto  true  "Email"
// @Success      200   {object}  common.ApiResponse[response.ResendOtpResDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Failure      500   {object}  common.ApiErrorResponse
// @Router       /auth/resend-otp [post]
func (authentication *AuthenticationController) ResendOtpController(c *gin.Context) {
	var ResendOtpReqDto request.ResendOtpReqDto

	err := c.ShouldBindJSON(&ResendOtpReqDto)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	res, businessErr := authentication.AuthenticationService.ResendOtp(c, &ResendOtpReqDto, authentication.Email)

	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{
			Code:    businessErr.StatusCode,
			Message: businessErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, common.ApiResponse[*response.ResendOtpResDto]{
		Code:    http.StatusOK,
		Message: "OTP resent successfully",
		Data:    res,
	})
}

// ValidateMailController godoc
// @Summary      Validate email with OTP
// @Description  Validate user email using the OTP sent to the email.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  request.ValidateMailReqDto  true  "Email and OTP"
// @Success      200   {object}  common.ApiResponse[response.ValidateMailResDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Failure      500   {object}  common.ApiErrorResponse
// @Router       /auth/validate-mail [post]
func (authentication *AuthenticationController) ValidateMailController(c *gin.Context) {
	var validateMailReq request.ValidateMailReqDto

	err := c.ShouldBindJSON(&validateMailReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	res, businessErr := authentication.AuthenticationService.ValidateMail(c, &validateMailReq)

	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{
			Code:    businessErr.StatusCode,
			Message: businessErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, common.ApiResponse[*response.ValidateMailResDto]{
		Code:    http.StatusOK,
		Message: "Validate mail successful",
		Data:    res,
	})
}

// RegisterController godoc
// @Summary      Register
// @Description  Register a new user with email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  request.RegisterReqDto  true  "Registration payload"
// @Success      200   {object}  common.ApiResponse[response.RegisterResDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Failure      500   {object}  common.ApiErrorResponse
// @Router       /auth/register [post]
func (authentication *AuthenticationController) RegisterController(c *gin.Context) {

	var registerReq request.RegisterReqDto

	err := c.ShouldBindJSON(&registerReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	res, businessErr := authentication.AuthenticationService.Register(c, &registerReq, authentication.Email)

	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{
			Code:    businessErr.StatusCode,
			Message: businessErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, common.ApiResponse[*response.RegisterResDto]{
		Code:    http.StatusOK,
		Message: "Register successful",
		Data:    res,
	})
}

// LoginController godoc
// @Summary      Login
// @Description  Login with email and password. Returns access token and sets refresh token in HTTP-only cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  request.LoginReqDto  true  "Email and password"
// @Success      200   {object}  common.ApiResponse[response.AccessTokenDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Failure      401   {object}  common.ApiErrorResponse
// @Failure      500   {object}  common.ApiErrorResponse
// @Router       /auth/login [post]
func (authentication *AuthenticationController) LoginController(c *gin.Context) {
	var loginReq request.LoginReqDto

	err := c.ShouldBindJSON(&loginReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Set IP and device info from request (server-derived, not client-supplied)
	ipAddress := c.ClientIP()
	loginReq.IpAddress = &ipAddress
	if deviceInfo := c.GetHeader("User-Agent"); deviceInfo != "" {
		loginReq.DeviceInfo = &deviceInfo
	}
	if v := c.GetHeader("X-Device-Info"); v != "" && (loginReq.DeviceInfo == nil || *loginReq.DeviceInfo == "") {
		loginReq.DeviceInfo = &v
	}

	res, businessErr := authentication.AuthenticationService.Login(c, &loginReq)

	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{
			Code:    businessErr.StatusCode,
			Message: businessErr.Message,
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("refreshToken", res.RefreshToken, constant.REFRESH_TOKEN_EXPIRE*60, "/", "", false, true)

	c.JSON(http.StatusOK, common.ApiResponse[*response.AccessTokenDto]{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data: &response.AccessTokenDto{
			AccessToken: res.AccessToken,
		},
	})
}

// LogoutController godoc
// @Summary      Logout
// @Description  Logout and invalidate refresh token. Clears refresh token cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  request.LogoutReqDto  true  "Refresh token (and optional access token)"
// @Success      200   {object}  common.ApiResponse[response.LogoutResDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Failure      401   {object}  common.ApiErrorResponse
// @Failure      500   {object}  common.ApiErrorResponse
// @Router       /auth/logout [post]
func (authentication *AuthenticationController) LogoutController(c *gin.Context) {
	var req request.LogoutReqDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	res, businessErr := authentication.AuthenticationService.Logout(c.Request.Context(), &req)

	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{
			Code:    businessErr.StatusCode,
			Message: businessErr.Message,
		})
		return
	}

	// Clear refresh token cookie
	c.SetCookie("refreshToken", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, common.ApiResponse[*response.LogoutResDto]{
		Code:    http.StatusOK,
		Message: "Logout successful",
		Data:    res,
	})
}

// ForgotPasswordController godoc
// @Summary      Forgot password
// @Description  Request a password reset OTP to be sent to the given email.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  request.ForgotPasswordReqDto  true  "Email"
// @Success      200   {object}  common.ApiResponse[response.ForgotPasswordResDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Failure      500   {object}  common.ApiErrorResponse
// @Router       /auth/forgot-password [post]
func (authentication *AuthenticationController) ForgotPasswordController(c *gin.Context) {
	var req request.ForgotPasswordReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	res, businessErr := authentication.AuthenticationService.ForgotPassword(c, &req, authentication.Email)
	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{Code: businessErr.StatusCode, Message: businessErr.Message})
		return
	}
	c.JSON(http.StatusOK, common.ApiResponse[*response.ForgotPasswordResDto]{
		Code: http.StatusOK, Message: "Password reset OTP sent",
		Data: res,
	})
}

// ResetPasswordController godoc
// @Summary      Reset password
// @Description  Reset password using email, OTP, and new password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  request.ResetPasswordReqDto  true  "Email, OTP, new password"
// @Success      200   {object}  common.ApiResponse[response.ResetPasswordResDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Failure      500   {object}  common.ApiErrorResponse
// @Router       /auth/reset-password [post]
func (authentication *AuthenticationController) ResetPasswordController(c *gin.Context) {
	var req request.ResetPasswordReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	res, businessErr := authentication.AuthenticationService.ResetPassword(c, &req)
	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{Code: businessErr.StatusCode, Message: businessErr.Message})
		return
	}
	c.JSON(http.StatusOK, common.ApiResponse[*response.ResetPasswordResDto]{
		Code: http.StatusOK, Message: "Password reset successful",
		Data: res,
	})
}

// ChangePasswordController godoc
// @Summary      Change password
// @Description  Change password for the authenticated user. Requires Bearer token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {accessToken}"
// @Param        body           body    request.ChangePasswordReqDto  true  "Old and new password"
// @Success      200   {object}  common.ApiResponse[response.ChangePasswordResDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Failure      401   {object}  common.ApiErrorResponse
// @Failure      500   {object}  common.ApiErrorResponse
// @Router       /auth/change-password [post]
func (authentication *AuthenticationController) ChangePasswordController(c *gin.Context) {
	var req request.ChangePasswordReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, common.ApiResponse[any]{Code: http.StatusUnauthorized, Message: "unauthorized"})
		return
	}
	res, businessErr := authentication.AuthenticationService.ChangePassword(c, user.Id, &req)
	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{Code: businessErr.StatusCode, Message: businessErr.Message})
		return
	}
	c.JSON(http.StatusOK, common.ApiResponse[*response.ChangePasswordResDto]{
		Code: http.StatusOK, Message: "Password changed successfully",
		Data: res,
	})
}

// IntrospectController godoc
// @Summary      Introspect token
// @Description  Validate an access token and return its claims (RFC 7662 style). Returns active=true with sub, exp, iat, email when valid; active=false otherwise.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  request.IntrospectReqDto  true  "Token to introspect"
// @Success      200   {object}  common.ApiResponse[response.IntrospectResDto]
// @Failure      400   {object}  common.ApiErrorResponse
// @Router       /auth/introspect [post]
func (authentication *AuthenticationController) IntrospectController(c *gin.Context) {
	var req request.IntrospectReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	res, _ := authentication.AuthenticationService.Introspect(c, &req)
	c.JSON(http.StatusOK, common.ApiResponse[*response.IntrospectResDto]{
		Code: http.StatusOK, Message: "OK", Data: res,
	})
}

// GoogleSSORedirectController godoc
// @Summary      Start Google SSO
// @Description  Redirects the user to Google's OAuth2 consent page. Browser only; follow redirect to sign in with Google.
// @Tags         auth
// @Produce      json
// @Success      307  "Redirect to Google"
// @Failure      500  {object}  common.ApiErrorResponse  "Failed to set SSO state"
// @Failure      503  {object}  common.ApiErrorResponse  "Google SSO not configured"
// @Router       /auth/google [get]
func (authentication *AuthenticationController) GoogleSSORedirectController(c *gin.Context) {
	state := uuid.New().String()
	if err := authentication.AuthenticationService.SetSSOState(c.Request.Context(), state); err != nil {
		c.JSON(http.StatusInternalServerError, common.ApiResponse[any]{Code: http.StatusInternalServerError, Message: "failed to set SSO state"})
		return
	}
	authURL := authentication.AuthenticationService.GetGoogleAuthURL(state)
	if authURL == "" {
		c.JSON(http.StatusServiceUnavailable, common.ApiResponse[any]{Code: http.StatusServiceUnavailable, Message: "Google SSO is not configured"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleSSOCallbackController godoc
// @Summary      Google SSO callback
// @Description  OAuth2 callback: validates state, exchanges code with Google, finds or creates user, returns access and refresh tokens.
// @Tags         auth
// @Param        state  query  string  true   "State from redirect (CSRF)"
// @Param        code  query  string  true   "Authorization code from Google"
// @Produce      json
// @Success      200  {object}  common.ApiResponse[response.LoginResDto]
// @Failure      400  {object}  common.ApiErrorResponse  "Missing or invalid state/code"
// @Failure      500  {object}  common.ApiErrorResponse
// @Router       /auth/google/callback [get]
func (authentication *AuthenticationController) GoogleSSOCallbackController(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{Code: http.StatusBadRequest, Message: "missing state or code"})
		return
	}
	ok, err := authentication.AuthenticationService.ValidateAndConsumeSSOState(c.Request.Context(), state)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{Code: http.StatusBadRequest, Message: "invalid or expired state"})
		return
	}
	res, businessErr := authentication.AuthenticationService.GoogleSSO(c.Request.Context(), code)
	if businessErr != nil {
		c.JSON(businessErr.StatusCode, common.ApiResponse[any]{Code: businessErr.StatusCode, Message: businessErr.Message})
		return
	}
	c.JSON(http.StatusOK, common.ApiResponse[*response.LoginResDto]{
		Code: http.StatusOK, Message: "Signed in with Google", Data: res,
	})
}
