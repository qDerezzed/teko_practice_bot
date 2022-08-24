package store

import (
	"context"
)

type ChatPosition int

const (
	Start ChatPosition = iota
	EnterCompanyName
	EnterCompanyID
	CompanyMenu
	EnterSumExpense
	EnterDescExpense
	EnterEmailExpense
)

type User struct {
	userID  int64
	chatPos ChatPosition
	name    string
}

func NewUser(userID int64, chatPos ChatPosition, userName string) *User {
	return &User{
		userID:  userID,
		chatPos: chatPos,
		name:    userName,
	}
}

func (s *Store) AddUser(user *User) error {
	_, err := s.dbPool.Exec(context.Background(),
		`INSERT INTO users (user_id, chat_position, user_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO 
		UPDATE SET chat_position = 0, user_name = EXCLUDED.user_name;`,
		user.userID, user.chatPos, user.name)
	return err
}

func (s *Store) SetChatPosition(chatPos ChatPosition, userID int64) error {
	_, err := s.dbPool.Exec(context.Background(),
		`UPDATE users SET chat_position = $1 WHERE user_id = $2;`,
		chatPos, userID)
	return err
}

func (s *Store) GetChatPosition(userID int64) (ChatPosition, error) {
	var chatPos ChatPosition
	err := s.dbPool.QueryRow(
		context.Background(),
		"SELECT chat_position FROM users WHERE user_id = $1;",
		userID).Scan(&chatPos)

	return chatPos, err
}

func (s *Store) GetUserName(userID int64) (string, error) {
	var name string
	err := s.dbPool.QueryRow(
		context.Background(),
		"SELECT user_name FROM users WHERE user_id = $1;",
		userID).Scan(&name)

	return name, err
}

func (s *Store) SetCurrentCompanyID(companyID int, userID int64) error {
	_, err := s.dbPool.Exec(context.Background(),
		`UPDATE users SET current_company_id = $1 WHERE user_id = $2;`,
		companyID, userID)
	return err
}

func (s *Store) SetCurrentExpenseID(expenseID int, userID int64) error {
	_, err := s.dbPool.Exec(context.Background(),
		`UPDATE users SET current_expense_id = $1 WHERE user_id = $2;`,
		expenseID, userID)
	return err
}
