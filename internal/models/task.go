package models

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // New import
)

type Task struct {
	Name        string         `json:"name"`
	Category    string         `json:"category"`
	Created     time.Time      `json:"created"`
	Expired     time.Time      `json:"expired"`
	IsDone      bool           `json:"isDone"`
	Whocomplete sql.NullString `json:"whocomplete"`
}

type TaskModel struct {
	DB *sql.DB
}

func (m *TaskModel) getAllCompanyProjectTasks(stmt string, companyID, projectID int) ([]*Task, error) {
	rows, err := m.DB.Query(stmt, companyID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Tasks := []*Task{}
	for rows.Next() {
		s := &Task{}
		err = rows.Scan(&s.Name, &s.Category, &s.Created, &s.Expired, &s.IsDone, &s.Whocomplete)
		if err != nil {
			return nil, err
		}
		Tasks = append(Tasks, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Tasks, nil
}

func (m *TaskModel) GetAllCompanyProjectTasks(companyID, projectID int) ([]*Task, error) {
	stmt := `SELECT name, category, created, expired, isDone, whocomplete 
		FROM tasks 
		WHERE companyID = ? AND projID = ?;`
	return m.getAllCompanyProjectTasks(stmt, companyID, projectID)
}

func (m *TaskModel) GetAllCompanyProjectTasksComplited(companyID, projectID int) ([]*Task, error) {
	stmt := `SELECT name, category, created, expired, isDone, whocomplete 
		FROM tasks 
		WHERE companyID = ? AND projID = ? AND tasks.isDone = TRUE;`
	return m.getAllCompanyProjectTasks(stmt, companyID, projectID)
}

func (m *TaskModel) GetAllCompanyProjectTasksUnomplited(companyID, projectID int) ([]*Task, error) {
	stmt := `SELECT name, category, created, expired, isDone, whocomplete 
		FROM tasks 
		WHERE companyID = ? AND projID = ? AND tasks.isDone = FALSE;`
	return m.getAllCompanyProjectTasks(stmt, companyID, projectID)
}

func (m *TaskModel) GetAllCompanyProjectTasksSortByCategoryDescending(companyID, projectID int) ([]*Task, error) {
	stmt := `SELECT name, category, created, expired, isDone, whocomplete 
		FROM tasks
	WHERE companyID = ? AND projID = ?
ORDER BY 
    CASE 
        WHEN tasks.category = 'Urgent' THEN 1
        WHEN tasks.category = 'High Priority' THEN 2
        WHEN tasks.category = 'Low Priority' THEN 3
        ELSE 4
    END DESC`
	return m.getAllCompanyProjectTasks(stmt, companyID, projectID)
}

func (m *TaskModel) GetAllCompanyProjectTasksSortByCategoryAscending(companyID, projectID int) ([]*Task, error) {
	stmt := `SELECT name, category, created, expired, isDone, whocomplete 
		FROM tasks
	WHERE companyID = ? AND projID = ?
ORDER BY 
    CASE 
        WHEN tasks.category = 'Urgent' THEN 1
        WHEN tasks.category = 'High Priority' THEN 2
        WHEN tasks.category = 'Low Priority' THEN 3
        ELSE 4
    END ASC`
	return m.getAllCompanyProjectTasks(stmt, companyID, projectID)
}

func (m *TaskModel) GetAllCompanyProjectTasksSortByExpiredAscending(companyID, projectID int) ([]*Task, error) {
	stmt := `SELECT name, category, created, expired, isDone, whocomplete 
		FROM tasks
	WHERE companyID = ? AND projID = ?
ORDER BY 
   tasks.expired ASC`
	return m.getAllCompanyProjectTasks(stmt, companyID, projectID)
}

func (m *TaskModel) GetAllCompanyProjectTasksSortByExpiredDescending(companyID, projectID int) ([]*Task, error) {
	stmt := `SELECT name, category, created, expired, isDone, whocomplete 
		FROM tasks
	WHERE companyID = ? AND projID = ?
ORDER BY 
   tasks.expired DESC`
	return m.getAllCompanyProjectTasks(stmt, companyID, projectID)
}

func (m *TaskModel) FindTask(name string, projID, companyID int) error {
	var count int
	stmt := `SELECT COUNT(*) FROM tasks WHERE name = ? AND companyID = ? AND projID = ?`
	err := m.DB.QueryRow(stmt, name, companyID, projID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrInvalidTaskName
	}
	return nil
}

func (m *TaskModel) Insert(TaskName, category string, expired time.Time, projID, companyID int) error {
	err := m.FindTask(TaskName, projID, companyID)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO tasks (name, category, created, expired, isDone, projID, companyID)
		VALUES(?, ?, UTC_TIMESTAMP(), ?, false, ?, ?)`

	_, err = m.DB.Exec(stmt, TaskName, category, expired, projID, companyID)
	if err != nil {
		return err
	}

	return nil
}

func (m *TaskModel) CompleateTask(projID, companyID, userID int, TaskName string) error {
	var isDone int
	checkStmt := `SELECT isDone FROM tasks WHERE name = ? AND projID = ? AND companyID = ?`
	err := m.DB.QueryRow(checkStmt, TaskName, projID, companyID).Scan(&isDone)
	if err != nil {
		if err == sql.ErrNoRows {
			return TaskNotFound
		}
		return err
	}

	// Если задача уже выполнена, возвращаем ошибку
	if isDone == 1 {
		return TaskIsAlredyDone
	}

	stmt := `UPDATE tasks 
	SET isDone = 1, 
		whocomplete = (
			SELECT name 
			FROM users 
			WHERE users.id = ?
		) 
	WHERE tasks.name = ? and tasks.projID = ? and tasks.companyID = ?;`

	_, err = m.DB.Exec(stmt, userID, TaskName, projID, companyID)
	if err != nil {
		return err
	}

	return nil
}
