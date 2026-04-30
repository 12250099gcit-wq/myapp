// myapp/routes/routes.go
package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"myapp/controller"
)

// NewRouter returns a configured mux router for the application.
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// ── Student routes ──────────────────────────────────────────
	router.HandleFunc("/students", controller.GetAllStudents).Methods("GET")
	router.HandleFunc("/student", controller.AddStudent).Methods("POST")
	router.HandleFunc("/student/{id}", controller.GetStudent).Methods("GET")
	router.HandleFunc("/student/{id}", controller.UpdateStudent).Methods("PUT")
	router.HandleFunc("/student/{id}", controller.DeleteStudent).Methods("DELETE")

	// ── Course routes ───────────────────────────────────────────
	router.HandleFunc("/courses", controller.GetAllCourses).Methods("GET")
	router.HandleFunc("/course", controller.AddCourse).Methods("POST")
	router.HandleFunc("/course/{id}", controller.GetCourse).Methods("GET")
	router.HandleFunc("/course/{id}", controller.UpdateCourse).Methods("PUT")
	router.HandleFunc("/course/{id}", controller.DeleteCourse).Methods("DELETE")

	// ── Enroll routes ───────────────────────────────────────────
	router.HandleFunc("/enrolls", controller.GetEnrolls).Methods("GET")
	router.HandleFunc("/enroll", controller.Enroll).Methods("POST")
	router.HandleFunc("/enroll/{sid}/{cid}", controller.GetEnroll).Methods("GET")
	router.HandleFunc("/enroll/{sid}/{cid}", controller.DeleteEnroll).Methods("DELETE")

	// ── Auth routes ─────────────────────────────────────────────
	router.HandleFunc("/signup", controller.Signup).Methods("POST")
	router.HandleFunc("/login", controller.Login).Methods("POST")
	router.HandleFunc("/logout", controller.Logout)

	// ── Static file server ──────────────────────────────────────
	fhandler := http.FileServer(http.Dir("./view"))
	router.PathPrefix("/").Handler(fhandler)

	return router
}

// StartServer starts the HTTP server on port 8080.
func StartServer() {
	router := NewRouter()
	port := "8080"
	log.Println("Application running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}



