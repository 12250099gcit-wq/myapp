package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"

	"myapp/routes"
	"myapp/utils/postgres"
)

var (
	serverURL  string
	testClient *http.Client
)

func TestMain(m *testing.M) {
	if err := postgres.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create cookie jar: %v\n", err)
		os.Exit(1)
	}
	testClient = &http.Client{Jar: jar}

	ts := httptest.NewServer(routes.NewRouter())
	defer ts.Close()
	serverURL = ts.URL

	if err := ensureAdminUser(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ensure admin user exists: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func ensureAdminUser() error {
	if err := loginAdmin("pc@gmail.com", "pass"); err == nil {
		return nil
	}

	if err := signupAdmin("Pema", "Choden", "pc@gmail.com", "pass"); err != nil {
		return err
	}

	return loginAdmin("pc@gmail.com", "pass")
}

func signupAdmin(fname, lname, email, password string) error {
	payload := map[string]string{
		"firstname": fname,
		"lastname":  lname,
		"email":     email,
		"password":  password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, serverURL+"/signup", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := testClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("unexpected signup status: %d", resp.StatusCode)
	}
	return nil
}

func loginAdmin(email, password string) error {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, serverURL+"/login", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := testClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d", resp.StatusCode)
	}
	return nil
}
