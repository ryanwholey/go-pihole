package pihole

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type SessionAPI interface {
	Post(ctx context.Context) (Session, error)
	Login(ctx context.Context) (Session, error)
	Delete(ctx context.Context, sessionID string) error
}

type sessionAPI struct {
	client *Client
}

type sessionRequest struct {
	Password string `json:"password"`
}

type sessionResponse struct {
	Session sessionSessionResponse `json:"session"`
	Error   sessionErrorResponse   `json:"error"`
}

type sessionSessionResponse struct {
	Valid    bool   `json:"valid"`
	TOTP     bool   `json:"totp"`
	SID      string `json:"sid"`
	CSRF     string `json:"csrf"`
	Validity int    `json:"validity"`
	Message  string `json:"message"`
}

type sessionErrorResponse struct {
	Key     string `json:"key"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

type Session struct {
	SID        string
	TOTP       bool
	CSRF       string
	Expiration time.Time
}

func (r sessionResponse) ToSession() Session {
	s := Session{
		SID:        r.Session.SID,
		CSRF:       r.Session.CSRF,
		TOTP:       r.Session.TOTP,
		Expiration: time.Now().Add(time.Duration(r.Session.Validity) * time.Second),
	}

	return s
}

var (
	ErrorSessionNotFound        = errors.New("session not found")
	ErrorSessionUnauthorized    = errors.New("unauthorized session request")
	ErrorSessionBadRequest      = errors.New("bad session request")
	ErrorSessionTooManyRequests = errors.New("too many session requests")
)

func (s *sessionAPI) Login(ctx context.Context) (Session, error) {
	session, err := s.Post(ctx)
	if err != nil {
		return Session{}, err
	}

	s.client.auth.sid = session.SID

	return session, nil
}

// Post creates a session
func (s *sessionAPI) Post(ctx context.Context) (Session, error) {
	res, err := s.client.Post(ctx, "/api/auth", sessionRequest{
		Password: s.client.password,
	})
	if err != nil {
		return Session{}, err
	}
	defer res.Body.Close()

	var sesRes sessionResponse
	if err := json.NewDecoder(res.Body).Decode(&sesRes); err != nil {
		return Session{}, err
	}

	switch res.StatusCode {
	case http.StatusOK:
		return sesRes.ToSession(), nil
	case http.StatusBadRequest:
		return Session{}, fmt.Errorf("%w: %s", ErrorSessionBadRequest, sesRes.Error.Message)
	case http.StatusTooManyRequests:
		return Session{}, fmt.Errorf("%w: %s", ErrorSessionTooManyRequests, sesRes.Error.Message)
	default:
		return Session{}, fmt.Errorf("unexpected status code %d: %s", res.StatusCode, sesRes.Error.Message)
	}
}

// Delete cancels an active session
func (s *sessionAPI) Delete(ctx context.Context, sessionID string) error {
	res, err := s.client.Delete(ctx, fmt.Sprintf("/api/auth/%s", sessionID))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrorSessionNotFound, sessionID)
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrorSessionUnauthorized, sessionID)
	default:
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
}
