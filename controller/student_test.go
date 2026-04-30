// myapp/controller/student_test.go
package controller_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAddStudent verifies that a new student can be created.
func TestAddStudent(t *testing.T) {
	url := serverURL + "/student"
	var jsonStr = []byte(`{"stdid":1002,"fname":"Sangay","lname":"Lhamo","email":"sl@gmail.com"}`)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.JSONEq(t, `{"status":"student added"}`, string(body))
}

// TestGetStudent verifies that a student can be retrieved by ID.
func TestGetStudent(t *testing.T) {
	r, err := testClient.Get(serverURL + "/student/1002")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	body, _ := io.ReadAll(r.Body)

	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.JSONEq(t,
		`{"stdid":1002,"fname":"Sangay","lname":"Lhamo","email":"sl@gmail.com"}`,
		string(body))
}

// TestDeleteStudent verifies that a student can be deleted.
func TestDeleteStudent(t *testing.T) {
	url := serverURL + "/student/1002"
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, err := testClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"status":"deleted"}`, string(body))
}

// TestStudentNotFound verifies 404 response when student does not exist.
// Run AFTER TestDeleteStudent so student 1002 no longer exists.
func TestStudentNotFound(t *testing.T) {
	assert := assert.New(t)

	r, err := testClient.Get(serverURL + "/student/1002")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	body, _ := io.ReadAll(r.Body)

	assert.Equal(http.StatusNotFound, r.StatusCode)
	assert.JSONEq(`{"error":"Student not found"}`, string(body))
}
