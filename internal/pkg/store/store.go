package store 

import (
	"database/sql"
	"fmt"
	"sync"
	"texting-app/internal/pkg/models"
)

const (
	DRIVER     = "mssql"
	CONNSTRING = ""
)

type store struct {
	db *sql.DB
}

func (p *store) Close() error {
	if p.db == nil {
		panic("Missing database connection")
	}

	return p.db.Close()
}

func (p *store) GetUser(username string) (*models.User, error) {
	if p.db == nil {
		panic("Missing database connection")
	}

	query := `
        SELECT TOP 1
            USERNAME
            , PASSWORD
        FROM USERS (NOLOCK)
        WHERE 1 = 1
            AND USERNAME = ?;
    `
	user := models.User{}
	err := p.db.QueryRow(query, username).Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("Error executing query:\n%s\n\t%v", query, err)
	}

	return &user, nil
}

var (
	provInstance *store
	once         sync.Once
)

func instantiate() (*store, error) {
	db, err := sql.Open(DRIVER, CONNSTRING)
	if err != nil {
		return nil, fmt.Errorf("Error on database open: %v", err)
	}

	return &store{db: db}, nil
}

func GetStore() (*store, error) {
	var err error
	once.Do(func() {
		provInstance, err = instantiate()
	})

	return provInstance, err
}
