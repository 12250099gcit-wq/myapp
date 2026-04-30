// myapp/model/student.go
package model

import (
	"database/sql"

	"myapp/utils/postgres"
)

// Student mirrors the columns in the student table.
// The json tags control how Go field names appear in JSON output.
type Student struct {
	StdId int64  `json:"stdid"`
	Fname string `json:"fname"`
	Lname string `json:"lname"`
	Email string `json:"email"`
}

// SQL query constants — using $1, $2 placeholders prevents SQL injection.
const (
	queryInsertStudent = `INSERT INTO student(StdId, fname, lname, email)
                          VALUES($1, $2, $3, $4) RETURNING StdId;`

	queryGetStudent = `SELECT StdId, fname, lname, email
                          FROM student WHERE StdId=$1;`

	queryGetAllStudents = `SELECT StdId, fname, lname, email FROM student;`

	queryUpdateStudent = `UPDATE student
                          SET fname=$1, lname=$2, email=$3
                          WHERE StdId=$4 RETURNING StdId;`

	queryDeleteStudent = `DELETE FROM student WHERE StdId=$1 RETURNING StdId;`
)

// Create inserts a new student record.
// (s *Student) is the "receiver" — the function belongs to the Student type.
func (s *Student) Create() error {
	row := postgres.Db.QueryRow(queryInsertStudent, s.StdId, s.Fname, s.Lname, s.Email)
	return row.Scan(&s.StdId) // Scan reads the RETURNING value back into the struct
}

// Get retrieves one student by ID.
func (s *Student) Get() error {
	row := postgres.Db.QueryRow(queryGetStudent, s.StdId)
	return row.Scan(&s.StdId, &s.Fname, &s.Lname, &s.Email)
}

// GetAll retrieves every student and returns them as a slice.
func GetAllStudents() ([]Student, error) {
	rows, err := postgres.Db.Query(queryGetAllStudents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.StdId, &s.Fname, &s.Lname, &s.Email); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

// Update modifies an existing student record.
func (s *Student) Update() error {
	row := postgres.Db.QueryRow(queryUpdateStudent, s.Fname, s.Lname, s.Email, s.StdId)
	return row.Scan(&s.StdId)
}

// Delete removes a student record.
func (s *Student) Delete() error {
	row := postgres.Db.QueryRow(queryDeleteStudent, s.StdId)
	return row.Scan(&s.StdId)
}

// ErrNoRows is a convenience alias for the sentinel error returned when
// a query finds no matching row.
var ErrNoRows = sql.ErrNoRows
