package repository

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

type provider struct {
	db *sql.DB
}

func (p *provider) Close() error {
	return p.db.Close()
}

func (p *provider) GetUser(username string) (*models.User, error) {
    query := ""
    user := models.User{}
	err := p.db.QueryRow(query, username).Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("Error executing query:\n%s\n\t%v", query, err)
	}

	return &user, nil
}

var (
	provInstance *provider
	once         sync.Once
)

func instantiate() (*provider, error) {
	db, err := sql.Open(DRIVER, CONNSTRING)
	if err != nil {
		return nil, fmt.Errorf("Error on database open: %v", err)
	}

	return &provider{db: db}, nil
}

func GetProvider() (*provider, error) {
	var err error
	once.Do(func() {
		provInstance, err = instantiate()
	})

	return provInstance, err
}
