// myapp/model/course.go
package model

import "myapp/utils/postgres"

// Course mirrors the columns in the course table.
type Course struct {
    Cid         string `json:"cid"`
    CourseName  string `json:"cname"`
    Description string `json:"description"`
}

const (
    queryInsertCourse = `INSERT INTO course(cid, cname, description)
                         VALUES($1, $2, $3) RETURNING cid;`

    queryGetCourse    = `SELECT cid, cname, description FROM course WHERE cid=$1;`

    queryGetAllCourses = `SELECT cid, cname, description FROM course;`

    queryUpdateCourse = `UPDATE course
                         SET cname=$1, description=$2
                         WHERE cid=$3 RETURNING cid;`

    queryDeleteCourse = `DELETE FROM course WHERE cid=$1 RETURNING cid;`
)

func (c *Course) Create() error {
    row := postgres.Db.QueryRow(queryInsertCourse, c.Cid, c.CourseName, c.Description)
    return row.Scan(&c.Cid)
}

func (c *Course) Get() error {
    row := postgres.Db.QueryRow(queryGetCourse, c.Cid)
    return row.Scan(&c.Cid, &c.CourseName, &c.Description)
}

func GetAllCourses() ([]Course, error) {
    rows, err := postgres.Db.Query(queryGetAllCourses)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    courses := []Course{}
    for rows.Next() {
        var c Course
        if err := rows.Scan(&c.Cid, &c.CourseName, &c.Description); err != nil {
            return nil, err
        }
        courses = append(courses, c)
    }
    return courses, nil
}

func (c *Course) Update() error {
    row := postgres.Db.QueryRow(queryUpdateCourse, c.CourseName, c.Description, c.Cid)
    return row.Scan(&c.Cid)
}

func (c *Course) Delete() error {
    row := postgres.Db.QueryRow(queryDeleteCourse, c.Cid)
    return row.Scan(&c.Cid)
}