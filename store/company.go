package store

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type Company struct {
	CompanyID    int       `db:"company_id"`
	Name         string    `db:"company_name"`
	CreationDate time.Time `db:"creation_date"`
	UserID       int64     `db:"fk_user_id"`
}

func NewCompany(UserID int64, name string, creationDate time.Time) *Company {
	return &Company{
		UserID:       UserID,
		Name:         name,
		CreationDate: creationDate,
	}
}

func (s *Store) AddCompany(company *Company) error {
	_, err := s.dbPool.Exec(context.Background(),
		`INSERT INTO companys (company_name, creation_date, fk_user_id)
		VALUES ($1, $2, $3);`,
		company.Name, company.CreationDate, company.UserID)
	return err
}

func (s *Store) GetAllCompanys() ([]*Company, error) {
	var companys []*Company
	err := pgxscan.Select(context.Background(), s.dbPool, &companys, "SELECT * FROM companys")

	return companys, err
}

func (s *Store) GetCurrentCompanyID(userID int64) (int, error) {
	var companyID int
	err := s.dbPool.QueryRow(
		context.Background(),
		"SELECT current_company_id FROM users WHERE user_id = $1;",
		userID).Scan(&companyID)

	return companyID, err
}

func (s *Store) IsValidCompanyID(companyID int) (bool, error) {
	var isValid bool
	err := s.dbPool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM companys WHERE company_id = $1);",
		companyID).Scan(&isValid)

	return isValid, err
}
