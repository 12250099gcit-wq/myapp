// myapp/utils/postgres/db.go
package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"strings"

	_ "github.com/lib/pq" // registers the "postgres" driver; blank import is intentional
)

// Db is the shared database connection used by all model files.
// It is exported (capital D) so other packages can access it.
var Db *sql.DB

// Connection defaults. Override with environment variables in tests or deployment.
// const (
//
//	defaultHost     = "dpg-d7p49rog4nts73b6inl0-a.singapore-postgres.render.com"
//	defaultPort     = 5432
//	defaultUser     = "postgres_admin"
//	defaultPassword = "2QEZRwJYOIa4s1OVQYWpi0OSxilE7BCd"
//	defaultDBName   = "my_db_ssa7"
//	defaultSSLMode  = "require"
//
// )
func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func quoteIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}

func openDB(host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslmode,
	)
	return sql.Open("postgres", connStr)
}

func ensureSchema() error {
	queries := []string{

		`CREATE TABLE IF NOT EXISTS student (
			StdId INT NOT NULL,
			fname VARCHAR(45) NOT NULL,
			lname VARCHAR(45),
			email VARCHAR(45) NOT NULL UNIQUE,
			PRIMARY KEY (StdId)
		);`,

		`CREATE TABLE IF NOT EXISTS course (
			courseid VARCHAR(45) NOT NULL,
			coursename VARCHAR(45) NOT NULL,
			PRIMARY KEY (courseid)
		);`,

		`CREATE TABLE IF NOT EXISTS admin (
			firstname VARCHAR(45) NOT NULL,
			lastname VARCHAR(45),
			email VARCHAR(45) NOT NULL,
			password VARCHAR(45) NOT NULL,
			PRIMARY KEY (email)
		);`,

		`CREATE TABLE IF NOT EXISTS enroll (
			std_id INT NOT NULL,
			course_id VARCHAR(45) NOT NULL,
			date_enrolled VARCHAR(45),

			PRIMARY KEY (std_id, course_id),

			CONSTRAINT course_fk 
				FOREIGN KEY (course_id) 
				REFERENCES course(courseid)
				ON DELETE CASCADE 
				ON UPDATE CASCADE,

			CONSTRAINT std_fk 
				FOREIGN KEY (std_id) 
				REFERENCES student(StdId)
				ON DELETE CASCADE 
				ON UPDATE CASCADE
		);`,
	}

	for _, query := range queries {
		if _, err := Db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

// Init opens and validates connection to PostgreSQL database
func Init() error {
	if Db != nil {
		return nil
	}

	// Use local PostgreSQL service
	// For local connections, use empty password (peer auth) or trust
	host := envOrDefault("POSTGRES_HOST", "localhost")
	portStr := envOrDefault("POSTGRES_PORT", "5432")
	user := envOrDefault("POSTGRES_USER", "postgres")
	password := envOrDefault("POSTGRES_PASSWORD", "")
	dbname := envOrDefault("POSTGRES_DB", "myapp")
	sslmode := envOrDefault("POSTGRES_SSLMODE", "disable")

	var err error
	Db, err = openDB(host, portStr, user, password, dbname, sslmode)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %w", err)
	}

	if err = Db.Ping(); err != nil {
		_ = Db.Close()
		Db = nil
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	if err = ensureSchema(); err != nil {
		return err
	}

	// log.Println("DB Config:", host, portStr, user, password, dbname, sslmode)
	log.Println("Database connection established")
	return nil
}
