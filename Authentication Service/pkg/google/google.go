package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

// OAuth2Client is the subset of OAuth2 config used for SSO (allows mocking in tests).
type OAuth2Client interface {
	ExchangeCodeAndGetUserInfo(ctx context.Context, code string) (*UserInfo, error)
	AuthCodeURL(state string) string
}

// UserInfo represents Google OAuth2 userinfo response.
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// OAuth2Config holds Google OAuth2 configuration.
type OAuth2Config struct {
	*oauth2.Config
}

// NewOAuth2Config returns a Google OAuth2 config for the given client ID, secret and redirect URL.
func NewOAuth2Config(clientID, clientSecret, redirectURL string) *OAuth2Config {
	return &OAuth2Config{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

// AuthCodeURL returns the URL to redirect the user to for Google sign-in.
func (c *OAuth2Config) AuthCodeURL(state string) string {
	return c.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
}

// Exchange exchanges the authorization code for tokens and returns the token.
func (c *OAuth2Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.Config.Exchange(ctx, code)
}

// UserInfo fetches the user's profile from Google's userinfo endpoint.
func (c *OAuth2Config) UserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google userinfo returned %d: %s", resp.StatusCode, string(body))
	}

	var info UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ExchangeCodeAndGetUserInfo exchanges the code for tokens and fetches user info in one call.
func (c *OAuth2Config) ExchangeCodeAndGetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	token, err := c.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return c.UserInfo(ctx, token)
}
