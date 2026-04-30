// myapp/controller/admin_test.go
package controller_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAdmLogin verifies that a registered admin can log in successfully.
func TestAdmLogin(t *testing.T) {
	url := serverURL + "/login"

	var jsonStr = []byte(`{"email":"pc@gmail.com", "password":"pass"}`)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"message":"login success"}`, string(body))
}

// TestAdmUserNotExist verifies that invalid credentials return 401 Unauthorized.
func TestAdmUserNotExist(t *testing.T) {
	url := serverURL + "/login"
	data := []byte(`{"email":"notexist@example.com", "password":"wrong"}`)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.JSONEq(t, `{"error":"sql: no rows in result set"}`, string(body))
}
