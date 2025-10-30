package utils

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/i-sub135/i-sub-orm/internal/constant"
)

func ScanRows(rows *sql.Rows, dest any) error {
	destVal := reflect.ValueOf(dest)

	if destVal.Kind() != reflect.Pointer {
		return constant.ErrDestination
	}

	destVal = destVal.Elem()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	switch destVal.Kind() {

	//destination == slice of struct
	case reflect.Slice:
		elemType := destVal.Type().Elem()
		for rows.Next() {
			elemPtr := reflect.New(elemType)
			if err := intoStruct(rows, elemPtr.Elem(), cols); err != nil {
				return err
			}
			destVal.Set(reflect.Append(destVal, elemPtr.Elem()))
		}
		return rows.Err()

	//destination == single struct
	case reflect.Struct:
		if rows.Next() {
			return intoStruct(rows, destVal, cols)
		}
		return sql.ErrNoRows
	default:
		return constant.ErrDestinationType
	}
}

func intoStruct(rows *sql.Rows, dest reflect.Value, cols []string) error {

	fieldMap := make(map[string]reflect.Value)
	tipe := dest.Type()

	for i := 0; i < dest.NumField(); i++ {
		f := tipe.Field(i)
		colName := strings.ToLower(f.Name)
		if tag := f.Tag.Get("db"); tag != "" {
			colName = strings.ToLower(tag)
		}
		fieldMap[colName] = dest.Field(i)
	}

	values := make([]any, len(cols))
	valuesPtr := make([]any, len(cols))

	for i, col := range cols {
		if f, ok := fieldMap[col]; ok && f.CanSet() {
			valuesPtr[i] = f.Addr().Interface()
		} else {
			var skip any
			valuesPtr[i] = &skip
		}
		values[i] = valuesPtr[i]
	}

	if err := rows.Scan(values...); err != nil {
		return err
	}

	return nil
}
