package models

import (
	"database/sql"
	"errors"  // New import
	"strings" // New import
	"time"

	"github.com/go-sql-driver/mysql" // New import
	"golang.org/x/crypto/bcrypt"     // New import
)

// Define a new User type. Notice how the field names and types align
// with the columns in the database "users" table?
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Role           string
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password associated with the given email. If
	// no matching email exists we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
		VALUES(?, ?, ?, UTC_TIMESTAMP())`
	// Use the Exec() method to insert the user details and hashed password
	// into the users table.
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type *mysql.MySQLError. If it does, the
		// error will be assigned to the mySQLError variable. We can then check
		// whether or not the error relates to our users_uc_email key by
		// checking if the error code equals 1062 and the contents of the error
		// message string. If it does, we return an ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users.email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// We'll use the Exists method to check if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) FindForIdOrUsername(usernameORid string) ([]*User, error) {
	stmt := "SELECT id, name FROM users WHERE id = ? OR name LIKE CONCAT('%', ?, '%')"

	rows, err := m.DB.Query(stmt, usernameORid, usernameORid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Users := []*User{}
	for rows.Next() {
		s := &User{}
		err = rows.Scan(&s.ID, &s.Name)
		if err != nil {
			return nil, err
		}
		Users = append(Users, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Users, nil

}
func (m *UserModel) IsUserASiteAdmin(userID int) (bool, error) {
	query := `	SELECT 1 
	FROM users 
	WHERE id = ? AND role = 'admin'	`
	var exists int
	err := m.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (m *UserModel) UsersTable() ([]*User, error) {
	query := `	SELECT * FROM users `
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Users := []*User{}
	for rows.Next() {
		s := &User{}

		err = rows.Scan(&s.ID, &s.Name, &s.Email, &s.HashedPassword, &s.Created, &s.Role)
		if err != nil {
			return nil, err
		}
		Users = append(Users, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Users, nil
}
func (m *UserModel) DeleteUser(userID int) (bool, error) {
	stmt := `DELETE FROM users WHERE id = ?`
	result, err := m.DB.Exec(stmt, userID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}
