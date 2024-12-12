package models

import (
	"database/sql"
)

// Define a new User type. Notice how the field names and types align
// with the columns in the database "users" table?
type Invite struct {
	UserID         int
	CompanyID      int
	AdditionalInfo string
	CompanyName    string
}

type InviteModel struct {
	DB *sql.DB
}

func (m *InviteModel) GetAllUserInvite(userID int) ([]*Invite, error) {
	stmt := `SELECT 
	invites.companyID,
	company.name AS CompanyName,
	invites.AdditionalInfo
	FROM invites
	JOIN company ON invites.companyID = company.id
	WHERE invites.userID = ?`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Invites := []*Invite{}
	for rows.Next() {
		s := &Invite{}
		err = rows.Scan(&s.CompanyID, &s.CompanyName, &s.AdditionalInfo)
		if err != nil {
			return nil, err
		}
		Invites = append(Invites, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Invites, nil
}
func (m *InviteModel) DeleteInvite(userID, companyID int) error {
	stmt := `DELETE FROM invites WHERE userID = ? AND companyID = ?`
	_, err := m.DB.Exec(stmt, userID, companyID)
	return err
}

func (m *InviteModel) AddInvite(userID, companyID int, AdditionalInfo string) error {
	stmt := `insert into invites (userID, companyID, AdditionalInfo) values (?, ?, ?)`
	_, err := m.DB.Exec(stmt, userID, companyID, AdditionalInfo)
	return err
}
