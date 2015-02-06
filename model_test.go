package model

import (
	"database/sql"
	"os/exec"
	"testing"

	_ "github.com/lib/pq"
)

func setup() *sql.DB {
	exec.Command("createdb", "model-test").Run()
	exec.Command("psql", "-f", "model.sql", "model-test").Run()

	db, err := sql.Open("postgres", "dbname=model-test")
	if err != nil {
		panic(err)
	}

	return db
}

func teardown(db *sql.DB) {
	db.Close()
	exec.Command("dropdb", "model-test").Run()
}

func TestAll(t *testing.T) {
	db := setup()
	defer teardown(db)

	type product struct {
		Id    int    `column:"id"`
		Name  string `column:"name"`
		Price int    `column:"price"`
	}

	// Create 2 records
	Create("products", product{Name: "luis's stuff", Price: 300}, db)
	Create("products", product{Name: "miguel's stuff", Price: 400}, db)

	var products []product
	err := All("products", &products, db)
	if err != nil {
		t.Error(err)
	}

	if len(products) != 2 {
		t.Error("slice length must be 2")
	} else if products[0].Name != "luis's stuff" && products[1].Name != "miguel's stuff" {
		t.Error("incorrect data")
	}
}

func TestCreate(t *testing.T) {
	db := setup()
	defer teardown(db)

	type product struct {
		Id    int    `column:"id"`
		Name  string `column:"name"`
		Price int    `column:"price"`
	}

	id, err := Create("products", product{Name: "luis's stuff", Price: 300}, db)
	if err != nil {
		t.Error(err)
	}
	if id != 1 {
		t.Error("id != 1")
	}
}

func TestUpdate(t *testing.T) {
	db := setup()
	defer teardown(db)

	type product struct {
		Id    int    `column:"id"`
		Name  string `column:"name"`
		Price int    `column:"price"`
	}

	id, err := Create("products", product{Name: "luis's stuff", Price: 300}, db)
	if err != nil {
		t.Error(err)
	}
	if id != 1 {
		t.Error("id != 1")
	}
}
