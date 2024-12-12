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
type Project struct {
	Id        int       `json:"-"`
	Name      string    `json:"name"`
	LeaderID  int       `json:"-"`
	CompanyID int       `json:"-"`
	Created   time.Time `json:"created"`
	Status    bool      `json:"status"`
}

type ProjectModel struct {
	DB *sql.DB
}

func (m *ProjectModel) FindProj(name string, idCompany int) error {
	var count int
	stmt := `SELECT COUNT(*) FROM projects WHERE name = ? AND companyID = ?`
	err := m.DB.QueryRow(stmt, name, idCompany).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrInvalidProjectName
	}
	return nil
}

func (m *ProjectModel) Insert(name string, idCompany, idLeader int) error {
	err := m.FindProj(name, idCompany)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO projects (name, leaderID ,companyID, created, status)
		VALUES(?,?,?, UTC_TIMESTAMP(), 0)`
	_, err = m.DB.Exec(stmt, name, idLeader, idCompany)
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

func (m *ProjectModel) GetIDForName(name string, idCompany int) (int, error) {
	var id int
	stmt := `SELECT id FROM projects WHERE companyID = ? AND name = ?`
	err := m.DB.QueryRow(stmt, idCompany, name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *ProjectModel) getAllCompanyProjects(stmt string, companyID int) ([]*Project, error) {
	rows, err := m.DB.Query(stmt, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Projects := []*Project{}
	for rows.Next() {
		s := &Project{}
		err = rows.Scan(&s.Name, &s.Created, &s.Status)
		if err != nil {
			return nil, err
		}
		Projects = append(Projects, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Projects, nil
}

func (m *ProjectModel) GetAllCompanyProjects(companyID int) ([]*Project, error) {
	stmt := `SELECT 
   		 projects.name AS project_name, 
    	projects.created AS project_created,
		projects.status AS project_status
	FROM 
    	projects
	JOIN 
    	company ON projects.companyID = company.id
	WHERE 
   		 company.id = ?`
	return m.getAllCompanyProjects(stmt, companyID)
}

func (m *ProjectModel) GetAllCompanyProjectsFilteComplited(companyID int) ([]*Project, error) {
	stmt := `SELECT 
   		 projects.name AS project_name, 
    	projects.created AS project_created,
		projects.status AS project_status
	FROM 
    	projects
	JOIN 
    	company ON projects.companyID = company.id
	WHERE 
    	company.id = ? AND projects.status = TRUE`
	return m.getAllCompanyProjects(stmt, companyID)
}

func (m *ProjectModel) GetAllCompanyProjectsFilterOutstanding(companyID int) ([]*Project, error) {
	stmt := `SELECT 
   		 projects.name AS project_name, 
    	projects.created AS project_created,
		projects.status AS project_status
	FROM 
    	projects
	JOIN 
    	company ON projects.companyID = company.id
	WHERE 
    	company.id = ? AND projects.status = FALSE`
	return m.getAllCompanyProjects(stmt, companyID)
}
