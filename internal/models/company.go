package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Define a new User type. Notice how the field names and types align
// with the columns in the database "users" table?
type Company struct {
	ID      int
	Name    string
	Owner   string
	Created time.Time
}

type CompanyModel struct {
	DB *sql.DB
}

func (m *CompanyModel) Insert(name string, ownerid int) (int, error) {

	id := -1

	stmt := `INSERT INTO company (name, owner ,created)
		VALUES(?,?, UTC_TIMESTAMP())`
	_, err := m.DB.Exec(stmt, name, ownerid)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1452 && strings.Contains(mySQLError.Message, "foreign key") {
				return id, ErrInvalidUserID
			}
		}
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "company.unique_name") {
				return id, ErrDuplicateCompanyName
			}
		}
		return id, err
	}

	stmt = "SELECT id FROM company WHERE name = ?"
	err = m.DB.QueryRow(stmt, name).Scan(&id)
	if err != nil {
		return id, ErrNoRecord
	}

	return id, nil
}
func (m *CompanyModel) CompanyTable() ([]*Company, error) {
	query := `	SELECT * FROM company `
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Companyies := []*Company{}
	for rows.Next() {
		s := &Company{}

		err = rows.Scan(&s.ID, &s.Name, &s.Owner, &s.Created)
		if err != nil {
			return nil, err
		}
		Companyies = append(Companyies, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Companyies, nil
}
func (m *CompanyModel) DeleteCompany(companyID int) (bool, error) {
	stmt := `DELETE FROM company WHERE id = ?`
	result, err := m.DB.Exec(stmt, companyID)
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
