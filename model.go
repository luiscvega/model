package model

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/luiscvega/squid"
)

type field struct {
	index int
	name  string
	kind  reflect.Kind
}

func GetIdAndValues(s interface{}) (int,[]interface{}, error) {
	var id int

	v := reflect.ValueOf(s)

	values := make([]interface{}, v.NumField()-1)

	count := 0
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Tag.Get("column") == "id" {
			id = int(v.Field(i).Int())
			continue
		}

		values[count] = v.Field(i).Interface()
		count++
	}

	return id, values, nil
}

func GetIdAndPointers(s interface{}) (int, []interface{}, error) {
	val := reflect.ValueOf(s).Elem()

	pointers := make([]interface{}, val.NumField()-1)

	count := 0
	var id int
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Tag.Get("column") == "id" {
			id = int(val.Field(i).Int())
			continue
		}

		pointers[count] = val.Field(i).Addr().Interface()
		count++
	}

	return id, pointers, nil
}

func GetFieldNames(t reflect.Type) ([]string, error) {
	fieldNames := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldNames[i] = t.Field(i).Tag.Get("column")
	}

	return fieldNames, nil
}

func GetField(sf reflect.StructField) field {
	return field{sf.Index[0], sf.Tag.Get("column"), sf.Type.Kind()}
}

func GetFields(s interface{}) ([]field, error) {
	t := reflect.TypeOf(s)

	fields := make([]field, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fields[i] = GetField(t.Field(i))
	}

	return fields, nil
}

func GetFieldsWithoutID(s interface{}) ([]field, error) {
	fields, err := GetFields(s)
	if err != nil {
		return nil, err
	}

	for i, field := range fields {
		if field.name == "id" {
			fields = append(fields[:i], fields[i+1:]...)
		}
	}

	return fields, nil
}

func Create(table string, s interface{}, db *sql.DB) (int, error) {
	var id int

	fields, err := GetFieldsWithoutID(s)
	if err != nil {
		return id, err
	}

	v := reflect.ValueOf(s)

	fs := make([]string, len(fields))
	vs := make([]interface{}, len(fields))
	for i, field := range fields {
		fs[i] = field.name
		vs[i] = v.Field(field.index).Interface()
	}

	stmt := squid.Insert(table, fs)

	err = db.QueryRow(stmt, vs...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func Fetch(table string, s interface{}, db *sql.DB) error {
	fieldNames, err := GetFieldNames(reflect.TypeOf(s).Elem())
	fieldNames = fieldNames[1:]
	if err != nil {
		return err
	}

	id, pointers, err := GetIdAndPointers(s)
	if err != nil {
		return err
	}

	stmt := squid.Fetch(table, fieldNames, id)

	err = db.QueryRow(stmt).Scan(pointers...)
	if err != nil {
		return err
	}

	return nil
}

func All(table string, listPtr interface{}, db *sql.DB) error {
	t := reflect.TypeOf(listPtr).Elem().Elem()

	fieldNames, err := GetFieldNames(t)
	if err != nil {
		return err
	}

	stmt := squid.SelectAll(table, fieldNames)

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

func Update(table string, s interface{}, db *sql.DB) error {
	fieldNames, err := GetFieldNames(reflect.TypeOf(s))
	fieldNames = fieldNames[1:]
	if err != nil {
		return err
	}

	id, values, err := GetIdAndValues(s)
	if err != nil {
		return err
	}

	stmt := squid.Update(table, fieldNames, id)

	res, err := db.Exec(stmt, values...)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return errors.New("sql: rows affected should be 1")
	}

	return nil
}

func Delete(table string, id int, db *sql.DB) error {
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
