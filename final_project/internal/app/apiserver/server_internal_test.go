package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"final_project/internal/app/model"
	"final_project/internal/app/store/teststore"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

// FIXME
func TestServer_HandleUsersCreate(t *testing.T) {
	s := newServer(teststore.New())
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: model.User{
				Login:    "user",
				Email:    "user@example.com",
				Password: "Password",
				Name:     "user",
				Surname:  "usersur",
				Is_admin: false,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", nil)

			s.ServeHTTP(rec, req)
			assert.Equal(t, rec.Code, http.StatusOK)
		})
	}

}

// FIXME
func TestServer_HandleSessionsCreate(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	store.User().Create(context.Background(), u)
	s := newServer(store)
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
