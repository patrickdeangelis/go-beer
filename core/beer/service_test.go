package beer_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/patrickdeangelis/go-beer/core/beer"
)

// the file name xxxxx_test.go it's important to compiler ignore it
// It's a good practice to create a package xxx_test for the package i'm testing
// All test functions should start with "Test"
// to run test: go test xxxx_test.go

func TestStore(t *testing.T) {
	b := &beer.Beer{
		ID: 1,
		Name: "Heineken",
		Type: beer.TypeLager,
		Style: beer.StylePale,
	}
	db, err := sql.Open("sqlite3", "../../data/beer_test.db")
	if err != nil {
		t.Fatalf("Erro ao conectar ao bando de dados: %s", err.Error())
	}
	defer db.Close()
	err = clearDB(db)
	if err != nil {
		t.Fatalf("Erro limpando o banco de dados: %s", err.Error())
	}
	service := beer.NewService(db)
	err = service.Store(b)
	if err != nil {
		t.Fatalf("Erro salvando no bando de dados: %s", err.Error())
	}
	saved, err := service.Get(1)
	if err != nil {
		t.Fatalf("Erro buscando no bando de dados: %s", err.Error())
	}
	if saved.ID != 1 {
		t.Fatalf("Dados inválidos. Esperado %d, recebido %d", 1, saved.ID)
	}
}

func clearDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from beer")
	tx.Commit()
	return err
}
