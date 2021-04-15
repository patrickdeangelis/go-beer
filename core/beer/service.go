package beer

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type UseCase interface {
	GetAll() ([]*Beer, error)
	Get(ID int64) (*Beer, error)
	Store(b *Beer) error
	Update(b *Beer) error
	Remove(ID int64) error
}


type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}

func (s *Service) GetAll() ([]*Beer, error) {
	var result []*Beer

	rows, err := s.DB.Query("select id, name, type, style from beer")
	if err != nil {
		return nil, err
	}
	// Defer ensures that the connection will be closed when the function ends
	defer rows.Close()

	for rows.Next() {
		var b Beer
		// Scan go on current line of database and for each column set the value
		// in this case for the beer object
		err := rows.Scan(&b.ID, &b.Name, &b.Type, &b.Style)
		if err != nil {
			return nil, err
		}
		result = append(result, &b)
	}
	return result, nil
}

func (s *Service) Get(ID int64) (*Beer, error) {
	var b Beer

	stmt, err := s.DB.Prepare("select id, name, type, style from beer where id =?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(ID).Scan(&b.ID, &b.Name, &b.Type, &b.Style)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *Service) Store(b *Beer) error {
	// To start a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into beer(id, name, type, style) values (?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(b.ID, b.Name, b.Type, b.Style)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *Service) Update(b *Beer) error {
	if b.ID == 0 {
		return fmt.Errorf("invalid ID")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("update beer set name=?, type=?, style=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(b.Name, b.Type, b.Style, b.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *Service) Remove(ID int64) error {
	if ID == 0 {
		return fmt.Errorf("invalid ID")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil
	}
	stmt, err := tx.Prepare("delete from beer where id=?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		tx.Rollback()
		return nil
	}
	tx.Commit()
	return nil
}

