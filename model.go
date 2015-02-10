package model

import (
	"database/sql"
	"reflect"
	"strconv"

	"github.com/luiscvega/squid"
)

func GetFields(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)

	fields := make([]string, t.NumField())

	for i := range fields {
		fields[i] = t.Field(i).Tag.Get("column")
	}

	return fields, nil
}

func Create(table string, s interface{}, db *sql.DB) (int, error) {
	var id int

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	tuples := make([][]string, 0, v.NumField()-1)
	for i := 0; i < v.NumField(); i++ {
		tuple := make([]string, 2)

		column := t.Field(i).Tag.Get("column")
		if column == "id" {
			continue
		}

		tuple[0] = column

		if v.Field(i).Kind() == reflect.String {
			tuple[1] = v.Field(i).String()
		} else {
			tuple[1] = strconv.Itoa(int(v.Field(i).Int()))
		}

		tuples = append(tuples, tuple)
	}

	stmt := squid.Insert(table, tuples)

	err := db.QueryRow(stmt).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

//func Fetch(table, s interface{}, db *sql.DB) error {
//stmt := squid.Fetch(table, s)

//err := db.QueryRow(stmt).Scan(fields...)
//if err != nil {
//return id, err
//}

//return nil
//}

//func Update(table string, s interface{}, db *sql.DB) error {
//stmt := squid.Update(table, s)

//res, err := db.Exec(stmt)
//if err != nil {
//return err
//}

//count, err := res.RowsAffected()
//if err != nil {
//return err
//}

//if count == 1 {
//return nil
//}

//return nil
//}

//func Delete(table string, id interface{}, db *sql.DB) error {
//stmt := squid.Delete(table, id)

//res, err := db.Exec(stmt)
//if err != nil {
//return err
//}

//count, err := res.RowsAffected()
//if err != nil {
//return err
//}

//if count == 1 {
//return nil
//}

//return errors.New("sql: rows affected should be 1")
//}

//func All(table string, listPtr interface{}, db *sql.DB) error {
//t := reflect.TypeOf(listPtr).Elem().Elem()
//stmt := squid.SelectAll(table, t)

//rows, err := db.Query(stmt)
//if err != nil {
//return err
//}
//defer rows.Close()

//listValue := reflect.ValueOf(listPtr).Elem()

//for rows.Next() {
//v := reflect.New(t).Elem()

//fields := make([]interface{}, 0)
//for i := 0; i < v.NumField(); i++ {
//fields = append(fields, v.Field(i).Addr().Interface())
//}

//rows.Scan(fields...)

//listValue.Set(reflect.Append(listValue, v))
//}

//return nil
//}
