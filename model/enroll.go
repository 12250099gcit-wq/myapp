// myapp/model/enroll.go
package model

import "myapp/utils/postgres"

// Enroll mirrors the columns in the enroll table.
type Enroll struct {
    StdId        int64  `json:"stdid"`
    CourseID     string `json:"cid"`
    Date_Enrolled string `json:"date"`
}

const (
    queryEnrollStd = `INSERT INTO enroll(std_id, course_id, date_enrolled)
                      VALUES($1, $2, $3) RETURNING std_id;`

    queryGetEnroll = `SELECT std_id, course_id, date_enrolled
                      FROM enroll WHERE std_id=$1 AND course_id=$2;`

    queryGetAllEnrolls = `SELECT std_id, course_id, date_enrolled FROM enroll;`

    queryDeleteEnroll = `DELETE FROM enroll
                         WHERE std_id=$1 AND course_id=$2 RETURNING std_id;`
)

// EnrollStud inserts a new enrollment record.
func (e *Enroll) EnrollStud() error {
    row := postgres.Db.QueryRow(queryEnrollStd, e.StdId, e.CourseID, e.Date_Enrolled)
    return row.Scan(&e.StdId)
}

// Get retrieves one enrollment by student ID and course ID.
func (e *Enroll) Get() error {
    return postgres.Db.QueryRow(queryGetEnroll, e.StdId, e.CourseID).
        Scan(&e.StdId, &e.CourseID, &e.Date_Enrolled)
}

// GetAllEnrolls retrieves every enrollment record.
func GetAllEnrolls() ([]Enroll, error) {
    rows, err := postgres.Db.Query(queryGetAllEnrolls)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    enrolls := []Enroll{}
    for rows.Next() {
        var e Enroll
        if dbErr := rows.Scan(&e.StdId, &e.CourseID, &e.Date_Enrolled); dbErr != nil {
            return nil, dbErr
        }
        enrolls = append(enrolls, e)
    }
    rows.Close()
    return enrolls, nil
}

// Delete removes an enrollment record.
func (e *Enroll) Delete() error {
    row := postgres.Db.QueryRow(queryDeleteEnroll, e.StdId, e.CourseID)
    return row.Scan(&e.StdId)
}