package service

import (
	"Authentication_Service/internal/constant"
	"Authentication_Service/internal/model"
	response "Authentication_Service/internal/dto/response"
	_interface "Authentication_Service/internal/service/interface"
	"Authentication_Service/pkg/google"
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"
)

// mockGoogleOAuth2 is a test double for google.OAuth2Client.
type mockGoogleOAuth2 struct {
	authURL string
	userInfo *google.UserInfo
	exchangeErr error
}

func (m *mockGoogleOAuth2) AuthCodeURL(state string) string {
	if m.authURL != "" {
		return m.authURL + "?state=" + state
	}
	return ""
}

func (m *mockGoogleOAuth2) ExchangeCodeAndGetUserInfo(ctx context.Context, code string) (*google.UserInfo, error) {
	if m.exchangeErr != nil {
		return nil, m.exchangeErr
	}
	return m.userInfo, nil
}

// mockAuthRepoSSO mocks IAuthenticationRepository for SSO (SaveRefreshToken, SetSSOState, ValidateAndConsumeSSOState).
type mockAuthRepoSSO struct {
	user              *model.User
	findByEmailErr    error
	saveRefreshErr    error
	setSSOStateErr    error
	validateStateOk   bool
	validateStateErr  error
}

func (m *mockAuthRepoSSO) FindByEmail(email string) (*model.User, error) {
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	return m.user, nil
}

func (m *mockAuthRepoSSO) SaveRefreshToken(ctx context.Context, record *model.RefreshToken) error {
	return m.saveRefreshErr
}

func (m *mockAuthRepoSSO) AddAccessTokenToBlacklist(ctx context.Context, accessToken string, ttl time.Duration) error {
	return nil
}

func (m *mockAuthRepoSSO) SetSSOState(ctx context.Context, state string, ttl time.Duration) error {
	return m.setSSOStateErr
}

func (m *mockAuthRepoSSO) ValidateAndConsumeSSOState(ctx context.Context, state string) (bool, error) {
	return m.validateStateOk, m.validateStateErr
}

// mockUserRepoSSO mocks IUserRepository for SSO.
type mockUserRepoSSO struct {
	user         *model.User
	findErr     error
	createErr   error
	updateErr   error
	createCalled bool
	updateCalled bool
}

func (m *mockUserRepoSSO) GetById(id string) (*model.User, error) { return nil, nil }
func (m *mockUserRepoSSO) FindByEmail(email string) (*model.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.user, nil
}
func (m *mockUserRepoSSO) Create(user *model.User) (*model.User, error) {
	m.createCalled = true
	if m.createErr != nil {
		return nil, m.createErr
	}
	return user, nil
}
func (m *mockUserRepoSSO) Update(user *model.User) (*model.User, error) {
	m.updateCalled = true
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return user, nil
}

// mockTokenGenSSO mocks ITokenGenerator for SSO.
type mockTokenGenSSO struct {
	accessToken  string
	refreshToken string
	genErr       error
}

func (m *mockTokenGenSSO) GenerateToken(user response.UserResDto, expireMinutes int64) (string, error) {
	if m.genErr != nil {
		return "", m.genErr
	}
	if expireMinutes == constant.ACCESS_TOKEN_EXPIRE {
		return m.accessToken, nil
	}
	return m.refreshToken, nil
}
func (m *mockTokenGenSSO) ValidateToken(tokenString string) (*response.UserResDto, error) { return nil, nil }
func (m *mockTokenGenSSO) GetExpiresAt(tokenString string) (time.Time, error)             { return time.Time{}, nil }
func (m *mockTokenGenSSO) Introspect(tokenString string) (*_interface.IntrospectResult, error) {
	return nil, nil
}

func TestAuthenticationService_GoogleSSO_NotConfigured(t *testing.T) {
	svc := &AuthenticationService{GoogleOAuth2: nil}
	res, err := svc.GoogleSSO(context.Background(), "any-code")
	if res != nil {
		t.Fatal("expected nil response")
	}
	if err == nil || err.StatusCode != 500 || err.Message != "Google SSO is not configured" {
		t.Fatalf("expected 500 Google SSO is not configured, got %v", err)
	}
}

func TestAuthenticationService_GoogleSSO_ExchangeError(t *testing.T) {
	svc := &AuthenticationService{
		GoogleOAuth2: &mockGoogleOAuth2{exchangeErr: errors.New("invalid code")},
	}
	res, err := svc.GoogleSSO(context.Background(), "bad-code")
	if res != nil {
		t.Fatal("expected nil response")
	}
	if err == nil || err.StatusCode != 400 || err.Message != "invalid or expired Google authorization code" {
		t.Fatalf("expected 400 invalid code, got %v", err)
	}
}

func TestAuthenticationService_GoogleSSO_EmptyEmail(t *testing.T) {
	svc := &AuthenticationService{
		GoogleOAuth2: &mockGoogleOAuth2{userInfo: &google.UserInfo{Email: ""}},
		UserRepository: &mockUserRepoSSO{},
	}
	res, err := svc.GoogleSSO(context.Background(), "code")
	if res != nil {
		t.Fatal("expected nil response")
	}
	if err == nil || err.StatusCode != 400 || err.Message != "Google did not return an email" {
		t.Fatalf("expected 400 no email, got %v", err)
	}
}

func TestAuthenticationService_GoogleSSO_NewUser(t *testing.T) {
	authRepo := &mockAuthRepoSSO{user: nil, findByEmailErr: gorm.ErrRecordNotFound}
	userRepo := &mockUserRepoSSO{user: nil, findErr: gorm.ErrRecordNotFound}
	tokenGen := &mockTokenGenSSO{accessToken: "at", refreshToken: "rt"}
	googleMock := &mockGoogleOAuth2{
		userInfo: &google.UserInfo{
			Email:     "new@example.com",
			GivenName: "New",
			FamilyName: "User",
		},
	}
	svc := &AuthenticationService{
		AuthenticationRepository: authRepo,
		UserRepository:          userRepo,
		TokenGenerator:          tokenGen,
		GoogleOAuth2:            googleMock,
	}
	// Auth repo FindByEmail is used by service; user repo FindByEmail is used
	// So we need authRepo to return nil (no user) - but the service uses UserRepository.FindByEmail
	authRepo.user = nil
	userRepo.user = nil
	userRepo.findErr = gorm.ErrRecordNotFound

	res, err := svc.GoogleSSO(context.Background(), "code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.AccessToken != "at" || res.RefreshToken != "rt" {
		t.Fatalf("expected tokens, got %+v", res)
	}
	if !userRepo.createCalled {
		t.Fatal("expected Create to be called for new user")
	}
}

func TestAuthenticationService_GoogleSSO_ExistingUserActive(t *testing.T) {
	first, last := "Existing", "User"
	existing := &model.User{
		Id:        "user-1",
		Email:     "existing@example.com",
		Password:  constant.GoogleSSOPasswordPlaceholder,
		Firstname: &first,
		Lastname:  &last,
		IsActive:  constant.ACTIVE,
	}
	authRepo := &mockAuthRepoSSO{}
	userRepo := &mockUserRepoSSO{user: existing}
	tokenGen := &mockTokenGenSSO{accessToken: "at2", refreshToken: "rt2"}
	googleMock := &mockGoogleOAuth2{
		userInfo: &google.UserInfo{Email: "existing@example.com", GivenName: "E", FamilyName: "U"},
	}
	svc := &AuthenticationService{
		AuthenticationRepository: authRepo,
		UserRepository:          userRepo,
		TokenGenerator:          tokenGen,
		GoogleOAuth2:            googleMock,
	}

	res, err := svc.GoogleSSO(context.Background(), "code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.AccessToken != "at2" || res.RefreshToken != "rt2" {
		t.Fatalf("expected tokens, got %+v", res)
	}
	if userRepo.createCalled {
		t.Fatal("expected Create not to be called for existing user")
	}
}

func TestAuthenticationService_GoogleSSO_ExistingUserInactive_Activates(t *testing.T) {
	first, last := "In", "Active"
	existing := &model.User{
		Id:        "user-2",
		Email:     "inactive@example.com",
		Firstname: &first,
		Lastname:  &last,
		IsActive:  constant.NO_ACTIVE,
	}
	userRepo := &mockUserRepoSSO{user: existing}
	tokenGen := &mockTokenGenSSO{accessToken: "at3", refreshToken: "rt3"}
	googleMock := &mockGoogleOAuth2{
		userInfo: &google.UserInfo{Email: "inactive@example.com", GivenName: "In", FamilyName: "Active"},
	}
	svc := &AuthenticationService{
		AuthenticationRepository: &mockAuthRepoSSO{},
		UserRepository:          userRepo,
		TokenGenerator:          tokenGen,
		GoogleOAuth2:            googleMock,
	}

	res, err := svc.GoogleSSO(context.Background(), "code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected success response")
	}
	if !userRepo.updateCalled {
		t.Fatal("expected Update to be called to activate user")
	}
	if userRepo.user == nil || userRepo.user.IsActive != constant.ACTIVE {
		t.Fatal("expected user to be updated to ACTIVE")
	}
}

func TestAuthenticationService_GoogleSSO_SaveRefreshTokenError(t *testing.T) {
	userRepo := &mockUserRepoSSO{user: nil, findErr: gorm.ErrRecordNotFound}
	authRepo := &mockAuthRepoSSO{user: nil, saveRefreshErr: errors.New("redis down")}
	tokenGen := &mockTokenGenSSO{accessToken: "at", refreshToken: "rt"}
	googleMock := &mockGoogleOAuth2{
		userInfo: &google.UserInfo{Email: "x@example.com", GivenName: "X", FamilyName: "Y"},
	}
	svc := &AuthenticationService{
		AuthenticationRepository: authRepo,
		UserRepository:          userRepo,
		TokenGenerator:          tokenGen,
		GoogleOAuth2:            googleMock,
	}

	res, err := svc.GoogleSSO(context.Background(), "code")
	if res != nil {
		t.Fatal("expected nil response")
	}
	if err == nil || err.StatusCode != 500 || err.Message != "failed to save session" {
		t.Fatalf("expected 500 save session, got %v", err)
	}
}

func TestAuthenticationService_GetGoogleAuthURL(t *testing.T) {
	t.Run("not_configured", func(t *testing.T) {
		svc := &AuthenticationService{GoogleOAuth2: nil}
		url := svc.GetGoogleAuthURL("state123")
		if url != "" {
			t.Fatalf("expected empty URL, got %s", url)
		}
	})
	t.Run("configured", func(t *testing.T) {
		svc := &AuthenticationService{
			GoogleOAuth2: &mockGoogleOAuth2{authURL: "https://accounts.google.com/o/oauth2/auth"},
		}
		url := svc.GetGoogleAuthURL("state123")
		if url == "" || url != "https://accounts.google.com/o/oauth2/auth?state=state123" {
			t.Fatalf("expected URL with state, got %s", url)
		}
	})
}

func TestAuthenticationService_SetSSOState(t *testing.T) {
	authRepo := &mockAuthRepoSSO{setSSOStateErr: errors.New("redis error")}
	svc := &AuthenticationService{AuthenticationRepository: authRepo}
	err := svc.SetSSOState(context.Background(), "s1")
	if err == nil || err.Error() != "redis error" {
		t.Fatalf("expected redis error, got %v", err)
	}
	authRepo.setSSOStateErr = nil
	err = svc.SetSSOState(context.Background(), "s2")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAuthenticationService_ValidateAndConsumeSSOState(t *testing.T) {
	authRepo := &mockAuthRepoSSO{validateStateOk: false, validateStateErr: errors.New("not found")}
	svc := &AuthenticationService{AuthenticationRepository: authRepo}
	ok, err := svc.ValidateAndConsumeSSOState(context.Background(), "bad")
	if ok || err == nil {
		t.Fatalf("expected false and error, got ok=%v err=%v", ok, err)
	}
	authRepo.validateStateOk = true
	authRepo.validateStateErr = nil
	ok, err = svc.ValidateAndConsumeSSOState(context.Background(), "good")
	if !ok || err != nil {
		t.Fatalf("expected true and nil, got ok=%v err=%v", ok, err)
	}
}
