package controller

import (
	"Authentication_Service/internal/dto/common"
	"Authentication_Service/internal/dto/request"
	response "Authentication_Service/internal/dto/response"
	"Authentication_Service/internal/model"
	"Authentication_Service/pkg/email"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// mockAuthServiceSSO mocks IAuthenticationService for SSO flow.
type mockAuthServiceSSO struct {
	authURL              string
	setSSOStateErr       error
	validateStateOk      bool
	validateStateErr     error
	googleSSORes         *response.LoginResDto
	googleSSOBusinessErr *model.BusinessError
}

func (m *mockAuthServiceSSO) SetSSOState(ctx context.Context, state string) error {
	return m.setSSOStateErr
}

func (m *mockAuthServiceSSO) GetGoogleAuthURL(state string) string {
	return m.authURL
}

func (m *mockAuthServiceSSO) ValidateAndConsumeSSOState(ctx context.Context, state string) (bool, error) {
	return m.validateStateOk, m.validateStateErr
}

func (m *mockAuthServiceSSO) GoogleSSO(ctx context.Context, code string) (*response.LoginResDto, *model.BusinessError) {
	return m.googleSSORes, m.googleSSOBusinessErr
}

// Stub other IAuthenticationService methods so the mock compiles (not used in SSO tests).
func (m *mockAuthServiceSSO) Login(ctx context.Context, user *request.LoginReqDto) (*response.LoginResDto, *model.BusinessError) {
	return nil, nil
}
func (m *mockAuthServiceSSO) Register(ctx context.Context, user *request.RegisterReqDto, email email.IEmailStrategy) (*response.RegisterResDto, *model.BusinessError) {
	return nil, nil
}
func (m *mockAuthServiceSSO) ValidateMail(ctx context.Context, req *request.ValidateMailReqDto) (*response.ValidateMailResDto, *model.BusinessError) {
	return nil, nil
}
func (m *mockAuthServiceSSO) Logout(ctx context.Context, req *request.LogoutReqDto) (*response.LogoutResDto, *model.BusinessError) {
	return nil, nil
}
func (m *mockAuthServiceSSO) ResendOtp(ctx context.Context, req *request.ResendOtpReqDto, email email.IEmailStrategy) (*response.ResendOtpResDto, *model.BusinessError) {
	return nil, nil
}
func (m *mockAuthServiceSSO) ForgotPassword(ctx context.Context, req *request.ForgotPasswordReqDto, email email.IEmailStrategy) (*response.ForgotPasswordResDto, *model.BusinessError) {
	return nil, nil
}
func (m *mockAuthServiceSSO) ResetPassword(ctx context.Context, req *request.ResetPasswordReqDto) (*response.ResetPasswordResDto, *model.BusinessError) {
	return nil, nil
}
func (m *mockAuthServiceSSO) ChangePassword(ctx context.Context, userID string, req *request.ChangePasswordReqDto) (*response.ChangePasswordResDto, *model.BusinessError) {
	return nil, nil
}
func (m *mockAuthServiceSSO) Introspect(ctx context.Context, req *request.IntrospectReqDto) (*response.IntrospectResDto, *model.BusinessError) {
	return nil, nil
}

func TestGoogleSSORedirectController_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/google", nil)

	mock := &mockAuthServiceSSO{
		authURL: "https://accounts.google.com/o/oauth2/auth?state=",
	}
	ctrl := &AuthenticationController{AuthenticationService: mock}
	ctrl.GoogleSSORedirectController(c)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if loc == "" || !strings.Contains(loc, "accounts.google.com") {
		t.Fatalf("expected Location to Google, got %s", loc)
	}
}

func TestGoogleSSORedirectController_NotConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/google", nil)

	mock := &mockAuthServiceSSO{authURL: ""} // empty = not configured
	ctrl := &AuthenticationController{AuthenticationService: mock}
	ctrl.GoogleSSORedirectController(c)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
	var body common.ApiResponse[any]
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Message != "Google SSO is not configured" {
		t.Fatalf("expected message about not configured, got %s", body.Message)
	}
}

func TestGoogleSSORedirectController_SetSSOStateError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/google", nil)

	mock := &mockAuthServiceSSO{
		authURL:        "https://accounts.google.com/",
		setSSOStateErr: errors.New("redis down"),
	}
	ctrl := &AuthenticationController{AuthenticationService: mock}
	ctrl.GoogleSSORedirectController(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	var body common.ApiResponse[any]
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Message != "failed to set SSO state" {
		t.Fatalf("expected message about set state, got %s", body.Message)
	}
}

func TestGoogleSSOCallbackController_MissingStateOrCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name string
		url  string
	}{
		{"missing_both", "/auth/google/callback"},
		{"missing_state", "/auth/google/callback?code=abc"},
		{"missing_code", "/auth/google/callback?state=xyz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, tt.url, nil)

			ctrl := &AuthenticationController{AuthenticationService: &mockAuthServiceSSO{}}
			ctrl.GoogleSSOCallbackController(c)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d", w.Code)
			}
			var body common.ApiResponse[any]
			if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body.Message != "missing state or code" {
				t.Fatalf("expected missing state or code, got %s", body.Message)
			}
		})
	}
}

func TestGoogleSSOCallbackController_InvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=bad&code=code123", nil)

	mock := &mockAuthServiceSSO{
		validateStateOk:  false,
		validateStateErr: errors.New("expired"),
	}
	ctrl := &AuthenticationController{AuthenticationService: mock}
	ctrl.GoogleSSOCallbackController(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var body common.ApiResponse[any]
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Message != "invalid or expired state" {
		t.Fatalf("expected invalid state message, got %s", body.Message)
	}
}

func TestGoogleSSOCallbackController_GoogleSSOError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=valid&code=code", nil)

	mock := &mockAuthServiceSSO{
		validateStateOk:      true,
		googleSSOBusinessErr: model.NewBusinessError(http.StatusBadRequest, "invalid or expired Google authorization code"),
	}
	ctrl := &AuthenticationController{AuthenticationService: mock}
	ctrl.GoogleSSOCallbackController(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var body common.ApiResponse[any]
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Message != "invalid or expired Google authorization code" {
		t.Fatalf("expected Google error message, got %s", body.Message)
	}
}

func TestGoogleSSOCallbackController_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=valid&code=goodcode", nil)

	mock := &mockAuthServiceSSO{
		validateStateOk: true,
		googleSSORes: &response.LoginResDto{
			AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		},
	}
	ctrl := &AuthenticationController{AuthenticationService: mock}
	ctrl.GoogleSSOCallbackController(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body common.ApiResponse[*response.AccessTokenDto]
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data == nil || body.Data.AccessToken != mock.googleSSORes.AccessToken {
		t.Fatalf("expected access token in body, got %+v", body.Data)
	}
	if body.Message != "Signed in with Google" {
		t.Fatalf("expected success message, got %s", body.Message)
	}

	setCookie := w.Header().Get("Set-Cookie")
	if !strings.Contains(setCookie, "refreshToken="+mock.googleSSORes.RefreshToken) {
		t.Fatalf("expected refresh token cookie to be set, got %s", setCookie)
	}
	if !strings.Contains(setCookie, "HttpOnly") {
		t.Fatalf("expected refresh token cookie to be HttpOnly, got %s", setCookie)
	}
}
