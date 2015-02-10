package model

import (
	"database/sql"
	"os/exec"
	"reflect"
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

type product struct {
	Id    int    `column:"id"`
	Name  string `column:"name"`
	Price int    `column:"price"`
}

func TestGetIdAndValues(t *testing.T) {
	p := product{Id: 1, Name: "Production", Price: 53421}

	id, actualValues, err := GetIdAndValues(p)
	if err != nil {
		t.Error(err)
	}
	if id != p.Id {
		t.Error(id, " != ", p.Id)
	}

	expectedValues := []interface{}{p.Name, p.Price}
	if !reflect.DeepEqual(actualValues, expectedValues) {
		t.Error(actualValues, " != ", expectedValues)
	}
}

func TestGetFields(t *testing.T) {
	actualFields, err := GetFields(product{})
	if err != nil {
		t.Error(err)
	}

	expectedFields := []field{{0, "id", reflect.Int}, {1, "name", reflect.String}, {2, "price", reflect.Int}}
	if !reflect.DeepEqual(actualFields, expectedFields) {
		t.Error(actualFields, " != ", expectedFields)
	}
}

func TestGetIdAndPointers(t *testing.T) {
	p := product{Id: 321, Name: "Luis's Cool Stuff", Price: 300}

	id, pointers, err := GetIdAndPointers(&p)
	if err != nil {
		t.Error(err)
	}
	if id != p.Id {
		t.Error(id, " != ", p.Id)
	}
	if pointers[0].(*string) != &p.Name {
		t.Error(pointers[0].(*string), " != ", &p.Name)
	}
	if pointers[1].(*int) != &p.Price {
		t.Error(pointers[1].(*int), " != ", &p.Price)
	}
}

func TestCreate(t *testing.T) {
	db := setup()
	defer teardown(db)

	id, err := Create("products", product{Name: "luis's stuff", Price: 300}, db)
	if err != nil {
		t.Error(err)
	}
	if id != 1 {
		t.Error("id != 1")
	}
}

func TestFetch(t *testing.T) {
	db := setup()
	defer teardown(db)

	id, _ := Create("products", product{Name: "luis's stuff", Price: 300}, db)

	p := product{Id: id}
	err := Fetch("products", &p, db)
	if err != nil {
		t.Error(err)
	}
	if p.Name != "luis's stuff" {
		t.Error(p.Name, "!= luis's stuff")
	}
	if p.Price != 300 {
		t.Error(p.Price, "!= 300")
	}
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

func TestUpdate(t *testing.T) {
	db := setup()
	defer teardown(db)

	p := product{Name: "luis's stuff", Price: 300}
	id, _ := Create("products", p, db)

	p.Id = id
	p.Name = "Miguel's Loco Stuff"
	p.Price = 4321

	err := Update("products", p, db)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	db := setup()
	defer teardown(db)

	type product struct {
		Id    int    `column:"id"`
		Name  string `column:"name"`
		Price int    `column:"price"`
	}

	id, _ := Create("products", product{Name: "luis's stuff", Price: 300}, db)

	err := Delete("products", id, db)
	if err != nil {
		t.Error(err)
	}
	if id != 1 {
		t.Error("id != 1")
	}
}
