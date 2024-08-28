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

func (s *store) checkConn() {
	if s.db == nil {
		panic("Missing database connection")
	}
}

func (s *store) queryError(q string, err error) error {
	return fmt.Errorf("Error executing query:\n%s\n\t%v", q, err)
}

func (s *store) Close() error {
	s.checkConn()
	return s.db.Close()
}

func (s *store) GetUser(usrn string) (*models.User, error) {
	s.checkConn()
	q := `
        SELECT TOP 1
            USERNAME
            , PASSWORD
        FROM USERS (NOLOCK)
        WHERE 1 = 1
            AND USERNAME = ?;
    `
	usr := models.User{}
	err := s.db.QueryRow(q, usrn).Scan(&usr.Username, &usr.Password)
	if err != nil {
		return nil, s.queryError(q, err)
	}

	return &usr, nil
}

func (s *store) CreateUser(usrn, pw string) (*models.User, error) {
	s.checkConn()
	q := `
        INSERT INTO
        USERS VALUES (
            ?,
            ?
        );

        SELECT TOP 1
            USERNAME
            , PASSWORD
        FROM USERS (NOLOCK)
        WHERE 1 = 1
            AND USERNAME = ?;
    `
	usr := models.User{}
	err := s.db.QueryRow(q, usrn, pw).Scan(&usr.Username, &usr.Password)
	if err != nil {
		return nil, s.queryError(q, err)
	}

	return &usr, nil
}

var (
	storeInst *store
	once      sync.Once
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
		storeInst, err = instantiate()
	})

	return storeInst, err
}
