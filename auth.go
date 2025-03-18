package pihole

import (
	"context"
	"encoding/json"
)

type AuthAPI interface {
	Authenticate(context.Context) error
}

type authAPI struct {
	client *Client
}

type authRequest struct {
	Password string `json:"password"`
}

type authResponse struct {
	Session authSessionResponse `json:"session"`
	Took    float64             `json:"took"`
}

type authSessionResponse struct {
	Valid    bool   `json:"valid"`
	TOTP     bool   `json:"totp"`
	SID      string `json:"sid"`
	CSRF     string `json:"csrf"`
	Validity int    `json:"validity"`
	Message  string `json:"message"`
}

func (a *authAPI) Authenticate(ctx context.Context) error {
	res, err := a.client.Post(ctx, "/api/auth", authRequest{
		Password: a.client.password,
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var authRes authResponse
	if err := json.NewDecoder(res.Body).Decode(&authRes); err != nil {
		return err
	}

	a.client.auth = auth{
		csrf: authRes.Session.CSRF,
		sid:  authRes.Session.SID,
	}

	return nil
}
