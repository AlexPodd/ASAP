package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Add a new ErrInvalidCredentials error. We'll use this later if a user
	// tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Add a new ErrDuplicateEmail error. We'll use this later if a user
	// tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")

	ErrDuplicateCompanyName = errors.New("models: duplicate company name")

	ErrInvalidUserID = errors.New("models: invalid userID")

	ErrInvalidProjectName = errors.New("models: duplicate project")

	ErrInvalidTaskName = errors.New("models: duplicate task")

	ErrWrongTimeFormat = errors.New("models: wrong time format")

	ErrDuplicateNameInCompany = errors.New("models: duplicate user in company")

	TaskNotFound = errors.New("models: task not found")

	TaskIsAlredyDone = errors.New("models: task is done")
)
