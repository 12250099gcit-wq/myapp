// myapp/controller/enroll.go
package controller

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "github.com/gorilla/mux"

    "myapp/model"
    "myapp/utils/date"
    "myapp/utils/httpResp"
)

// Enroll handles POST /enroll
func Enroll(w http.ResponseWriter, r *http.Request) {
    if !VerifyCookie(w, r) {
        return
    }

    var e model.Enroll
    if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
        httpResp.RespondWithError(w, http.StatusBadRequest, "invalid json body")
        return
    }
    defer r.Body.Close()

    // Automatically set the enrollment date to now
    e.Date_Enrolled = date.GetDate()

    if saveErr := e.EnrollStud(); saveErr != nil {
        if strings.Contains(saveErr.Error(), "duplicate key") {
            httpResp.RespondWithError(w, http.StatusForbidden, "Duplicate keys")
            return
        }
        httpResp.RespondWithError(w, http.StatusInternalServerError, saveErr.Error())
        return
    }
    httpResp.RespondWithJSON(w, http.StatusCreated, map[string]string{"status": "enrolled"})
}

// GetEnroll handles GET /enroll/{sid}/{cid}
func GetEnroll(w http.ResponseWriter, r *http.Request) {
    if !VerifyCookie(w, r) {
        return
    }

    params := mux.Vars(r)
    sid := params["sid"]
    cid := params["cid"]

    stdid, _ := strconv.ParseInt(sid, 10, 64)
    e := model.Enroll{StdId: stdid, CourseID: cid}

    if getErr := e.Get(); getErr != nil {
        switch getErr.Error() {
        case "sql: no rows in result set":
            httpResp.RespondWithError(w, http.StatusNotFound, "No such enrollments")
        default:
            httpResp.RespondWithError(w, http.StatusInternalServerError, getErr.Error())
        }
        return
    }
    httpResp.RespondWithJSON(w, http.StatusOK, e)
}

// GetEnrolls handles GET /enrolls
func GetEnrolls(w http.ResponseWriter, r *http.Request) {
    if !VerifyCookie(w, r) {
        return
    }

    enrolls, getErr := model.GetAllEnrolls()
    if getErr != nil {
        httpResp.RespondWithError(w, http.StatusBadRequest, getErr.Error())
        return
    }
    httpResp.RespondWithJSON(w, http.StatusOK, enrolls)
}

// DeleteEnroll handles DELETE /enroll/{sid}/{cid}
func DeleteEnroll(w http.ResponseWriter, r *http.Request) {
    if !VerifyCookie(w, r) {
        return
    }

    params := mux.Vars(r)
    sid := params["sid"]
    cid := params["cid"]

    stdid, _ := strconv.ParseInt(sid, 10, 64)
    e := model.Enroll{StdId: stdid, CourseID: cid}

    if err := e.Delete(); err != nil {
        httpResp.RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }
    httpResp.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}