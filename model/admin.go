// myapp/model/admin.go
package model

import "myapp/utils/postgres"

// Admin mirrors the columns in the admin table.
type Admin struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

const (
	queryInsertAdmin = `INSERT INTO admin(firstname, lastname, email, password)
                        VALUES($1, $2, $3, $4) RETURNING email;`

	queryGetAdmin = `SELECT email, password FROM admin
                     WHERE email=$1 AND password=$2;`
)

// Create inserts a new admin (sign-up).
func (adm *Admin) Create() error {
	row := postgres.Db.QueryRow(queryInsertAdmin,
		adm.FirstName, adm.LastName, adm.Email, adm.Password)
	return row.Scan(&adm.Email)
}

// Get validates admin credentials (login).
func (adm *Admin) Get() error {
	return postgres.Db.QueryRow(queryGetAdmin, adm.Email, adm.Password).
		Scan(&adm.Email, &adm.Password)
}
