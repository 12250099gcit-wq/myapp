// myapp/controller/student.go
package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"

	"myapp/model"
	"myapp/utils/httpResp"
)

// GetAllStudents handles GET /students
// Returns all students as a JSON array.
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	students, err := model.GetAllStudents()
	if err != nil {
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, students)
}

// AddStudent handles POST /student
// Reads one student from the request body and inserts it.
func AddStudent(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	var stud model.Student
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&stud); err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()

	if saveErr := stud.Create(); saveErr != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, saveErr.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusCreated, map[string]string{"status": "student added"})
}

// GetStudent handles GET /student/{id}
// Returns a single student by ID.
func GetStudent(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	// mux.Vars(r) extracts URL path parameters, e.g. /student/1001 → id = "1001"
	params := mux.Vars(r)
	var stud model.Student
	// strconv.ParseInt converts the string "1001" to int64 1001
	id, err := parseInt64(params["id"])
	if err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid student id")
		return
	}
	stud.StdId = id

	if getErr := stud.Get(); getErr != nil {
		switch getErr {
		case sql.ErrNoRows:
			httpResp.RespondWithError(w, http.StatusNotFound, "Student not found")
		default:
			httpResp.RespondWithError(w, http.StatusInternalServerError, getErr.Error())
		}
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, stud)
}

// UpdateStudent handles PUT /student/{id}
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	params := mux.Vars(r)
	id, err := parseInt64(params["id"])
	if err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	var stud model.Student
	if err := json.NewDecoder(r.Body).Decode(&stud); err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()
	stud.StdId = id

	if updateErr := stud.Update(); updateErr != nil {
		httpResp.RespondWithError(w, http.StatusInternalServerError, updateErr.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteStudent handles DELETE /student/{id}
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	params := mux.Vars(r)
	id, err := parseInt64(params["id"])
	if err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	stud := model.Student{StdId: id}
	if delErr := stud.Delete(); delErr != nil {
		httpResp.RespondWithError(w, http.StatusInternalServerError, delErr.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// parseInt64 is a small helper so we don't repeat strconv calls everywhere.
func parseInt64(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
