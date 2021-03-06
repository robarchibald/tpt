package tpt

import (
	"fmt"

	"github.com/daemonl/tpt.go/tptobjects"
)

// User wraps an oAuth-ish token which can be used for user authenticated API
// calls. It extends RequestBuilder for custom http calls, and wraps some of
// the API calls to the TPT API
type User struct {
	Token  string
	Client *Client
}

// NewRequest returns a request builder which is pre-populated with
// authentication headers for the user and application
func (u *User) NewRequest(reqPath string) *Request {
	if len(u.Token) < 1 {
		return &Request{
			firstError: fmt.Errorf("User has no token"),
		}
	}
	return u.Client.NewRequest(reqPath).AddHeader("User-Token", u.Token)
}

/////////////////////////
// Wrapped API Methods //
/////////////////////////

// RevokeToken revokes an existing access token. This token will no longer be
// able to be used for authentication.
func (u *User) RevokeToken() error {
	resp := &struct {
		Revoked bool `json:"revoked"`
	}{}
	err := u.NewRequest("/v1/user/oauth/revoke").PostJSON(map[string]string{
		"token": u.Token,
	}).DecodeInto(resp)
	if err != nil {
		return err
	}
	if resp.Revoked {
		u.Token = ""
		return nil
	}
	return fmt.Errorf("Not Revoked")
}

// GetAccountDetails returns the user’s account details.
func (u *User) GetAccountDetails() (*tptobjects.UserAccountDetails, error) {
	resp := &tptobjects.UserAccountDetails{}
	err := u.NewRequest("/v1/user/account").DecodeInto(resp)
	return resp, err
}
