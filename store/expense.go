package store

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type Expense struct {
	ExpenseID    int       `db:"expense_id"`
	Sum          int       `db:"sum"`
	Description  string    `db:"description"`
	Email        string    `db:"email"`
	CreationDate time.Time `db:"creation_date"`
	CompanyId    int       `db:"fk_company_id"`
	UserID       int64     `db:"fk_user_id"`
}

func NewExpense(sum int, description string, email string,
	creationDate time.Time, companyId int, userID int64) *Expense {
	return &Expense{
		Sum:          sum,
		Description:  description,
		Email:        email,
		CreationDate: creationDate,
		CompanyId:    companyId,
		UserID:       userID,
	}
}

func (s *Store) AddExpense(expense *Expense) (int, error) {
	var expenseID int
	err := s.dbPool.QueryRow(context.Background(),
		`INSERT INTO expenses (sum, description, email, creation_date, fk_company_id, fk_user_id)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING expense_id;`,
		expense.Sum, expense.Description, expense.Email, expense.CreationDate,
		expense.CompanyId, expense.UserID).Scan(&expenseID)
	return expenseID, err
}

func (s *Store) GetLastExpenses(companyID int, count int) ([]*Expense, error) {
	var expenses []*Expense
	err := pgxscan.Select(
		context.Background(), s.dbPool, &expenses,
		`SELECT * FROM expenses WHERE fk_company_id = $1 
		ORDER BY creation_date DESC 
		LIMIT $2;`,
		companyID, count,
	)

	return expenses, err
}

func (s *Store) GetSumPerDay(companyID int) (int, error) {
	var sum int
	err := s.dbPool.QueryRow(
		context.Background(),
		"SELECT COALESCE(SUM(sum), 0) FROM expenses WHERE fk_company_id = $1 AND creation_date::date = $2;",
		companyID, time.Now().Format("2006-01-02")).Scan(&sum)

	return sum, err
}

func (s *Store) GetCurrentExpenseID(userID int64) (int, error) {
	var expenseID int
	err := s.dbPool.QueryRow(
		context.Background(),
		"SELECT current_expense_id FROM users WHERE user_id = $1;",
		userID).Scan(&expenseID)

	return expenseID, err
}

func (s *Store) SetEmailExpense(email string, expenseID int) error {
	_, err := s.dbPool.Exec(context.Background(),
		`UPDATE expenses SET email = $1 WHERE expense_id = $2;`,
		email, expenseID)
	return err
}

func (s *Store) SetDescriptionExpense(desc string, expenseID int) error {
	_, err := s.dbPool.Exec(context.Background(),
		`UPDATE expenses SET description = $1 WHERE expense_id = $2;`,
		desc, expenseID)
	return err
}

func (s *Store) DeleteExpense(expenseID int) error {
	_, err := s.dbPool.Exec(context.Background(),
		`DELETE FROM expenses WHERE expense_id = $1;`, expenseID)
	return err
}
