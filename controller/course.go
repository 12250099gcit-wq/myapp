// myapp/controller/course.go
package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"myapp/model"
	"myapp/utils/httpResp"
)

func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	courses, err := model.GetAllCourses()
	if err != nil {
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, courses)
}

func AddCourse(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	var c model.Course
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()

	if saveErr := c.Create(); saveErr != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, saveErr.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusCreated, map[string]string{"status": "course added"})
}

func GetCourse(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	params := mux.Vars(r)
	c := model.Course{Cid: params["id"]}

	if err := c.Get(); err != nil {
		switch err {
		case sql.ErrNoRows:
			httpResp.RespondWithError(w, http.StatusNotFound, "Course not found")
		default:
			httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, c)
}

func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	params := mux.Vars(r)
	var c model.Course
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()
	c.Cid = params["id"]

	if err := c.Update(); err != nil {
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	if !VerifyCookie(w, r) {
		return
	}

	params := mux.Vars(r)
	c := model.Course{Cid: params["id"]}

	if err := c.Delete(); err != nil {
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
