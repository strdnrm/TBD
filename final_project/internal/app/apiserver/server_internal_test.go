package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"final_project/internal/app/model"
	"final_project/internal/app/store/teststore"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gopkg.in/go-playground/assert.v1"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	s := newServer(teststore.New(), sessions.NewCookieStore([]byte("aaaaa")))
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid pass",
			payload: map[string]interface{}{
				"email":    "invalid",
				"password": "a",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleSessionsCreate(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	store.User().Create(context.Background(), u)
	s := newServer(store, sessions.NewCookieStore([]byte("aaaaa")))
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_AuthenticateUser(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	store.User().Create(context.Background(), u)

	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.Id,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("aaaaa")
	s := newServer(store, sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)
	mw := s.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			mw.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
