package model

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/luiscvega/squid"
)

func Delete(table string, id string, db *sql.DB) error {
	stmt := squid.Delete(table, id)

	res, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 1 {
		return nil
	}

	return errors.New("sql: rows affected should be 1")
}

func Create(table string, s interface{}, db *sql.DB) (int, error) {
	var id int

	stmt := squid.Insert(table, s)

	err := db.QueryRow(stmt).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func All(table string, listPtr interface{}, db *sql.DB) error {
	t := reflect.TypeOf(listPtr).Elem().Elem()
	stmt := squid.SelectAll(table, t)

	rows, err := db.Query(stmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	listValue := reflect.ValueOf(listPtr).Elem()

	for rows.Next() {
		v := reflect.New(t).Elem()

		fields := make([]interface{}, 0)
		for i := 0; i < v.NumField(); i++ {
			fields = append(fields, v.Field(i).Addr().Interface())
		}

		rows.Scan(fields...)

		listValue.Set(reflect.Append(listValue, v))
	}

	return nil
}
