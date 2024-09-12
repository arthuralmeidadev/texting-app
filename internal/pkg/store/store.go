package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"sync"
	"texting-app/internal/pkg/models"

	_ "github.com/lib/pq"
)

type store struct {
	db *sql.DB
}

func numStr[T int | uint](n T) string {
	return strconv.Itoa(int(n))
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
        SELECT
            USERNAME
            , PASSWORD
        FROM MAIN.USER
        WHERE 1 = 1
            AND USERNAME = $1
        LIMIT 1;
    `
	usr := models.User{}
	err := s.db.QueryRow(q, usrn).Scan(&usr.Username, &usr.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("no rows")
		}
		return nil, s.queryError(q, err)
	}

	return &usr, nil
}

func (s *store) CreateUser(usrn, pw string) (*models.User, error) {
	s.checkConn()
	q := `
        INSERT INTO MAIN.USER
        VALUES (
            $1
            , $2
            , CURRENT_TIMESTAMP
        )
        RETURNING USERNAME, PASSWORD;
    `
	usr := models.User{}
	err := s.db.QueryRow(q, usrn, pw).Scan(&usr.Username, &usr.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("no rows")
		}
		return nil, s.queryError(q, err)
	}

	return &usr, nil
}

func (s *store) FindUser(forUsrn, usrn string, offset uint) ([]*models.User, error) {
	s.checkConn()
	q := `
        SELECT
            MAIN.USER.USERNAME
        FROM MAIN.USER
            LEFT JOIN MAIN.FRIENDSHIP
                ON MAIN.FRIENDSHIP.SENDER = MAIN.USER.USERNAME
                OR MAIN.FRIENDSHIP.RECIPIENT = MAIN.USER.USERNAME
        WHERE 1 = 1
            AND MAIN.FRIENDSHIP.SENDER IS NULL
            AND MAIN.FRIENDSHIP.RECIPIENT IS NULL
            AND MAIN.USER.USERNAME LIKE $1
            AND MAIN.USER.USERNAME <> $2
        ORDER BY MAIN.USER.USERNAME
        OFFSET $3 LIMIT $4;
    `

	limit := 20
	rows, err := s.db.Query(q, fmt.Sprint("%", usrn, "%"), forUsrn, numStr(offset), limit)
	if err != nil {
		fmt.Println(err)
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("no rows")
		}
		return nil, s.queryError(q, err)
	}
	defer rows.Close()

	usrs := make([]*models.User, offset)
	for rows.Next() {
		usr := models.User{Password: ""}
		err := rows.Scan(&usr.Username)
		if err != nil {
			return nil, err
		}

		usrs = append(usrs, &usr)
	}

	return usrs, nil
}

func (s *store) GetUserFriends(usrn string, offset uint) ([]*models.User, error) {
	s.checkConn()
	q := `
        SELECT
            MAIN.USER.USERNAME 
        FROM MAIN.FRIENDSHIP
            INNER JOIN MAIN.USER
                ON MAIN.USER.USERNAME = MAIN.FRIENDSHIP.SENDER
                OR MAIN.USER.USERNAME = MAIN.FRIENDSHIP.RECIPIENT
        WHERE 1 = 1
            AND MAIN.USER.USERNAME <> $1
            AND MAIN.FRIENDSHIP.STATUS = 'ACCEPTED'
        ORDER BY MAIN.USER.USERNAME
        OFFSET $2 LIMIT $3;
    `

	limit := 20
	rows, err := s.db.Query(q, usrn, numStr(offset), limit)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("no rows")
		}
		return nil, s.queryError(q, err)
	}
	defer rows.Close()

	usrs := make([]*models.User, offset)
	for rows.Next() {
		usr := models.User{Password: ""}
		err := rows.Scan(&usr.Username)
		if err != nil {
			return nil, err
		}

		usrs = append(usrs, &usr)
	}

	return usrs, nil
}

func (s *store) StoreMessage(sen, cont string, chatId uint, repTo int) error {
	s.checkConn()
	q := `
        INSERT INTO MAIN.MESSAGE
        VALUES (
            DEFAULT
            , $1
            , $2
            , $3
            , CURRENT_TIMESTAMP 
            , CASE WHEN $4 = -1 THEN NULL ELSE $4 END
        ) RETURNING ID;
    `
	var msgId uint
	err := s.db.QueryRow(
		q,
		sen,
		numStr(chatId),
		cont,
		numStr(repTo),
	).Scan(&msgId)
	if err != nil {
		return s.queryError(q, err)
	}

	q = `
        INSERT INTO MAIN.MESSAGE_STATE
        VALUES (
            DEFAULT
            , $1
            , RECIPIENT
            , NULL
            , NULL
        )
        SELECT USERNAME AS RECIPIENT
        FROM MAIN.CHAT_MEMBER
        WHERE 1 = 1
            AND CHAT = $2;
    `

	_, err = s.db.Exec(q, numStr(msgId), numStr(chatId))
	if err != nil {
		return s.queryError(q, err)
	}

	q = `
        UPDATE MAIN.CHAT
        SET LAST_ACTIVE = CURRENT_TIMESTAMP
        WHERE 1 = 1
            AND CHAT = $1
    `

	_, err = s.db.Exec(q, numStr(chatId))
	if err != nil {
		return s.queryError(q, err)
	}

	return nil
}

func (s *store) UpdateMessageStatus(msgId uint, stat string) error {
	s.checkConn()
	if !slices.Contains([]string{"DELIVERED", "SEEN"}, stat) {
		return errors.New("Invalid message status")
	}

	q := fmt.Sprintf(`
        UPDATE MAIN.MESSAGE
        SET %s_AT = CURRENT_TIMESTAMP
        WHERE 1 = 1
            AND ID = $1;
    `, stat)

	_, err := s.db.Exec(q, numStr(msgId))
	if err != nil {
		return s.queryError(q, err)
	}

	return nil
}

func (s *store) GetMessages(chatId, offset uint) ([]*models.Message, error) {
	s.checkConn()
	q := `
        SELECT
            ID
            , SENDER
            , CHAT
            , CONTENT
            , SENT_AT
            , REPLIES_TO
        FROM MAIN.MESSAGE
        WHERE 1 = 1
            AND CHAT = $2
        ORDER BY SENT_AT DESC
        OFFSET $4 LIMIT $5;
    `
	limit := 20
	rows, err := s.db.Query(q, chatId, numStr(offset), limit)
	if err != nil {
		return nil, s.queryError(q, err)
	}
	defer rows.Close()

	msgs := make([]*models.Message, limit)
	for rows.Next() {
		msg := models.Message{}
		err := rows.Scan(
			&msg.Id,
			&msg.Sender,
			&msg.Chat,
			&msg.Content,
			&msg.SentAt,
			&msg.RepliesTo,
		)

		if err != nil {
			return nil, err
		}

		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

func (s *store) GetUnseenMessagesTotal(usrn, chatId string) (uint, error) {
	return 0, nil
}

func (s *store) DeleteMessage(msgId uint) error {
	return nil
}

func (s *store) CreateChat(name string, usrs []string) (uint, error) {
	s.checkConn()
	q := `
        INSERT INTO MAIN.CHAT
        VALUES (
            DEFAULT
            , $1
            , CURRENT_TIMESTAMP
            , CURRENT_TIMESTAMP
        )
        RETURNING ID;
    `

	var chatId uint
	err := s.db.QueryRow(q, name).Scan(&chatId)
	if err != nil {
		return 0, s.queryError(q, err)
	}

	return chatId, nil
}

func (s *store) GetChats(usrn string, offset uint) ([]*models.Chat, error) {
	s.checkConn()
	q := `
        SELECT
            MAIN.CHAT.ID
            , (
                SELECT
                    CASE
                        WHEN COUNT(*) = 1 THEN $1
                        WHEN COUNT(*) <= 2 THEN (
                            SELECT USERNAME
                            FROM MAIN.CHAT_MEMBER
                            WHERE 1 = 1
                                AND USERNAME <> $1
                        )
                    ELSE NAME
                END AS NAME
                FROM MAIN.CHAT_MEMBER
                WHERE 1 = 1
                    AND MAIN.CHAT_MEMBER.CHAT = MAIN.CHAT.ID
            ) AS NAME
        MAIN.CHAT.LAST_ACTIVE
        FROM MAIN.CHAT
            INNER JOIN MAIN.CHAT_MEMBER
                ON MAIN.CHAT_MEMBER.CHAT = MAIN.CHAT.ID WHERE 1 = 1
            AND MAIN.CHAT_MEMBER.USERNAME = $1
        ORDER BY SENT_AT DESC
        OFFSET $2 LIMIT $3;
    `

	limit := 20
	rows, err := s.db.Query(q, usrn, numStr(offset), limit)
	if err != nil {
		return nil, s.queryError(q, err)
	}
	defer rows.Close()

	chats := make([]*models.Chat, limit)
	for rows.Next() {
		chat := models.Chat{}
		err := rows.Scan(&chat.Id, &chat.Name, &chat.LastActive)
		if err != nil {
			return nil, err
		}

		chats = append(chats, &chat)
	}

	return chats, nil
}

func (s *store) GetChatMembers(chatId uint) ([]string, error) {
	s.checkConn()
	q := `
        SELECT
            USERNAME
        FROM MAIN.CHAT_MEMBER
        WHERE 1 = 1
            AND CHAT = $1
    `

	rows, err := s.db.Query(q, &chatId)
	if err != nil {
		return nil, s.queryError(q, err)
	}
	defer rows.Close()

	members := make([]string, 1)
	for rows.Next() {
		var member string
		err := rows.Scan(&member)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (s *store) SendFriendRequest(usrn, recUsrn string) error {
	s.checkConn()
	q := `
        INSERT INTO MAIN.FRIENDSHIP
        VALUES (
            DEFAULT
            , $1
            , $2
            , 'PENDING'
            , CURRENT_TIMESTAMP
        )
    `

	_, err := s.db.Exec(q, usrn, recUsrn)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) GetFriendRequests(usrn string, offset uint) ([]*models.FriendRequest, error) {
	s.checkConn()
	q := `
        SELECT
            SENDER
            , RECIPIENT
            , STATUS
            , LAST_MODIFIED
        FROM MAIN.FRIENDSHIP
        WHERE 1 = 1
            AND (
                SENDER = $1
                OR RECIPIENT = $1
            )
        OFFSET $2 LIMIT $3;
    `

	limit := 20
	rows, err := s.db.Query(q, usrn, numStr(offset), limit)
	if err != nil {
		return nil, s.queryError(q, err)
	}
	defer rows.Close()

    friendReqs := make([]*models.FriendRequest, limit)
    for rows.Next() {
        friendReq := models.FriendRequest{}
        err := rows.Scan(&friendReq.SendUsrn, &friendReq.RecUsrn, &friendReq.Status, &friendReq.LastMod)
        if err != nil {
            return nil, err
        }
        friendReqs = append(friendReqs, &friendReq)
    }

	return friendReqs, nil
}

func (s *store) DeleteChat(chatId uint) error {
	return nil
}

var (
	storeInst *store
	once      sync.Once
)

func formatDSN() string {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsr := os.Getenv("DB_USER")
	dbPw := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUsr, dbPw, dbName)
}

func instantiate() (*store, error) {
	driver := os.Getenv("DB_DRIVER")
	dsn := formatDSN()
	db, err := sql.Open(driver, dsn)
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
