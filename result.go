package simplemysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/rnben/simplemysql/utils"
)

func (sess *session) formatOutput(rows *sql.Rows, dest interface{}) error {
	var result []map[string]interface{}

	cols, _ := rows.Columns()
	for rows.Next() {
		row := make([]interface{}, len(cols))
		rowPointers := make([]interface{}, len(cols))

		for i := range row {
			rowPointers[i] = &row[i]
		}

		if err := rows.Scan(rowPointers...); err != nil {
			return err
		}

		result = append(result, sess.assignSliceMapConvert(row, cols))
	}
	err := rows.Err()
	if err != nil {
		return err
	}

	js, err := json.Marshal(result)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, &dest)

	return err
}

func (sess *session) assignSliceMapConvert(rowVal []interface{}, rowCols []string) map[string]interface{} {
	if rowVal == nil || rowCols == nil {
		return nil
	}

	var (
		colsNum = len(rowVal)
		data    = map[string]interface{}{}
	)

	for idx, field := range rowCols {
		if idx >= colsNum {
			continue
		}

		data[field] = convert(sess.fields[field], rowVal[idx])
	}

	return data
}

func convert(st string, data interface{}) interface{} {
	switch st {
	case "int", "int32", "int64":
		return utils.ToInt64(data)
	case "float", "float32", "float64":
		return utils.ToFloat64(data)
	case "string":
		return utils.ToString(data)
	case "bool":
		return utils.ToBool(data)
	default:
		return data
	}
}

func getSliceFields(dest interface{}) (map[string]string, error) {
	destType1 := reflect.TypeOf(dest)
	if destType1.Kind() != reflect.Ptr {
		return nil, errors.New("not ptr")
	}

	destType2 := destType1.Elem() // []myStruct
	if destType2.Kind() != reflect.Slice {
		return nil, errors.New("not slice")
	}

	destType3 := destType2.Elem()
	if destType3.Kind() != reflect.Ptr {
		return nil, errors.New("elem not ptr")
	}

	destType4 := destType3.Elem()
	if destType4.Kind() != reflect.Struct {
		return nil, errors.New("slice not struct")
	}

	structFields := destType4.NumField()
	fields := make(map[string]string, structFields)

	for i := 0; i < structFields; i++ {
		f := destType4.Field(i)
		fields[f.Tag.Get("json")] = f.Type.String()
	}

	return fields, nil
}
