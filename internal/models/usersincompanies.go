package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// Define a new User type. Notice how the field names and types align
// with the columns in the database "users" table?
type Usersincompanies struct {
	UserID    int
	CompanyID int
	Role      string
}
type UserWithRole struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
type CompanyWithUsers struct {
	CompanyID   int    `json:"companyID"`
	CompanyName string `json:"companyName"`
	Role        string `json:"role"`
}

type UsersincompaniesModel struct {
	DB *sql.DB
}

func (m *UsersincompaniesModel) Insert(userID, companyID int, role string) error {
	flag, err := m.IsUserInCompany(userID, companyID)
	if flag {
		return ErrDuplicateNameInCompany
	} else {
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	stmt := `INSERT INTO usersincompanies (user , company ,role)
		VALUES(?,?,?)`

	_, err = m.DB.Exec(stmt, userID, companyID, role)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1452 && strings.Contains(mySQLError.Message, "foreign key") {
				return ErrInvalidUserID
			}
		}
		return err
	}
	return nil
}

// Добавить обработку случая пустого списка!
// Добавить имя компании
func (m *UsersincompaniesModel) getAllUsers(stmt string, companyID int) ([]*UserWithRole, error) {
	rows, err := m.DB.Query(stmt, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	UsersWithRole := []*UserWithRole{}
	for rows.Next() {
		s := &UserWithRole{}
		err = rows.Scan(&s.UserID, &s.Username, &s.Role)
		if err != nil {
			return nil, err
		}
		UsersWithRole = append(UsersWithRole, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return UsersWithRole, nil
}

func (m *UsersincompaniesModel) GetAllUsers(companyID int) ([]*UserWithRole, error) {
	stmt := `SELECT  users.id, users.name AS user_name, usersInCompanies.role AS user_role
FROM usersInCompanies JOIN users ON usersInCompanies.user = users.id
JOIN company ON usersInCompanies.company = company.id
WHERE company.id = ?`
	return m.getAllUsers(stmt, companyID)
}

func (m *UsersincompaniesModel) GetAllUsersSortByRoleAscending(companyID int) ([]*UserWithRole, error) {
	stmt := `SELECT users.id, users.name AS user_name, usersInCompanies.role AS user_role
	FROM usersInCompanies JOIN users ON usersInCompanies.user = users.id
	JOIN company ON usersInCompanies.company = company.id
	WHERE company.id = ?
ORDER BY 
    CASE 
        WHEN usersInCompanies.role = 'owner' THEN 1
        WHEN usersInCompanies.role = 'admin' THEN 2
        WHEN usersInCompanies.role = 'worker' THEN 3
        ELSE 4
    END ASC`
	return m.getAllUsers(stmt, companyID)
}

func (m *UsersincompaniesModel) GetAllUsersSortByRoleDescending(companyID int) ([]*UserWithRole, error) {
	stmt := `SELECT users.id, users.name AS user_name, usersInCompanies.role AS user_role
	FROM usersInCompanies JOIN users ON usersInCompanies.user = users.id
	JOIN company ON usersInCompanies.company = company.id
	WHERE company.id = ?
ORDER BY 
    CASE 
        WHEN usersInCompanies.role = 'owner' THEN 1
        WHEN usersInCompanies.role = 'admin' THEN 2
        WHEN usersInCompanies.role = 'worker' THEN 3
        ELSE 4
    END DESC`
	return m.getAllUsers(stmt, companyID)
}

func (m *UsersincompaniesModel) GetAllCompanyWhereUserAdminOrOwner(userID int) ([]*CompanyWithUsers, error) {
	stmt := `SELECT 
    company.id AS company_id, 
    company.name AS company_name
FROM 
    usersInCompanies 
JOIN 
    company ON usersInCompanies.company = company.id
WHERE 
    usersInCompanies.user = ? 
    AND (usersInCompanies.role = 'admin' OR usersInCompanies.role = 'owner');`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	CompaniesWithUsers := []*CompanyWithUsers{}
	for rows.Next() {
		s := &CompanyWithUsers{}
		err = rows.Scan(&s.CompanyID, &s.CompanyName)
		if err != nil {
			return nil, err
		}
		CompaniesWithUsers = append(CompaniesWithUsers, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return CompaniesWithUsers, nil
}

func (m *UsersincompaniesModel) IsUserInCompany(userID int, companyID int) (bool, error) {
	query := "SELECT 1 FROM usersincompanies WHERE user = ? AND company = ?"
	var exists int
	err := m.DB.QueryRow(query, userID, companyID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *UsersincompaniesModel) IsUserAAdminOrOwner(userID int, companyID int) (bool, error) {
	query := `
        SELECT 1 
        FROM usersincompanies 
        WHERE user = ? AND company = ? AND role IN ('admin', 'owner')
    `
	var exists int
	err := m.DB.QueryRow(query, userID, companyID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (m *UsersincompaniesModel) IsUserAOwner(userID int, companyID int) (bool, error) {
	query := `	SELECT 1 
	FROM usersincompanies 
	WHERE user = ? AND company = ? AND role = 'owner'	`
	var exists int
	err := m.DB.QueryRow(query, userID, companyID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *UsersincompaniesModel) getAllCompanyWithUser(stmt string, userID int) ([]*CompanyWithUsers, error) {
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	CompaniesWithUsers := []*CompanyWithUsers{}
	for rows.Next() {
		s := &CompanyWithUsers{}
		err = rows.Scan(&s.CompanyID, &s.CompanyName, &s.Role)
		if err != nil {
			return nil, err
		}
		CompaniesWithUsers = append(CompaniesWithUsers, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return CompaniesWithUsers, nil
}

func (m *UsersincompaniesModel) GetAllCompanyWithUser(userID int) ([]*CompanyWithUsers, error) {
	stmt := `SELECT company.id AS company_id, company.name AS company_name,
    usersInCompanies.role AS user_role FROM usersInCompanies 
	JOIN 
    company ON usersInCompanies.company = company.id
	WHERE 
    usersInCompanies.user = ?`

	return m.getAllCompanyWithUser(stmt, userID)
}

func (m *UsersincompaniesModel) GetAllCompanyWithUserSortByRoleAscending(userID int) ([]*CompanyWithUsers, error) {
	stmt := `SELECT company.id AS company_id, company.name AS company_name,
    usersInCompanies.role AS user_role 
FROM usersInCompanies
JOIN company ON usersInCompanies.company = company.id
WHERE usersInCompanies.user = ?
ORDER BY 
    CASE 
        WHEN usersInCompanies.role = 'owner' THEN 1
        WHEN usersInCompanies.role = 'admin' THEN 2
        WHEN usersInCompanies.role = 'worker' THEN 3
        ELSE 4
    END ASC`

	return m.getAllCompanyWithUser(stmt, userID)
}

func (m *UsersincompaniesModel) GetAllCompanyWithUserSortByRoleDescending(userID int) ([]*CompanyWithUsers, error) {
	stmt := `SELECT company.id AS company_id, company.name AS company_name,
    usersInCompanies.role AS user_role 
FROM usersInCompanies
JOIN company ON usersInCompanies.company = company.id
WHERE usersInCompanies.user = ?
ORDER BY 
    CASE 
        WHEN usersInCompanies.role = 'owner' THEN 1
        WHEN usersInCompanies.role = 'admin' THEN 2
        WHEN usersInCompanies.role = 'worker' THEN 3
        ELSE 4
    END DESC`

	return m.getAllCompanyWithUser(stmt, userID)
}

func (m *UsersincompaniesModel) DeleteUser(userID, companyID int) (bool, error) {
	stmt := `DELETE FROM usersincompanies WHERE user = ? AND company = ?`
	result, err := m.DB.Exec(stmt, userID, companyID)
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
func (m *UsersincompaniesModel) SetAdminRole(userID, companyID int) (bool, error) {
	stmt := `UPDATE usersincompanies SET role = 'admin' WHERE user = ? AND company = ?`
	result, err := m.DB.Exec(stmt, userID, companyID)
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
func (m *UsersincompaniesModel) SetWorkerRole(userID, companyID int) (bool, error) {
	stmt := `UPDATE usersincompanies SET role = 'worker' WHERE user = ? AND company = ?`
	result, err := m.DB.Exec(stmt, userID, companyID)
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
