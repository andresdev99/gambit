package tools

import (
	"fmt"
	"github.com/andresdev99/gambit/models"
	"reflect"
	"strings"
	"time"
)

type Operation string

const (
	Insert Operation = "insert"
	Update Operation = "update"
)

var modelProperties = map[string]string{
	"products": "Prod_id",
	"category": "Categ_id",
}

type modelTypes interface {
	models.Product | models.Category
}

func BuildSQL[mType modelTypes](model mType, table string, op Operation) (string, []any, string, error) {
	v := reflect.ValueOf(model)

	pk, ok := modelProperties[table]
	if !ok {
		return "", nil, "", fmt.Errorf("no primary key defined for table %q", table)
	}

	if v.Kind() != reflect.Struct {
		return "", nil, "", fmt.Errorf("input must be a struct")
	}

	t := v.Type()
	var fields, holders, updates []string
	var args []any
	var pkValue any

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		val := v.Field(i).Interface()
		if dbTag == pk {
			pkValue = val
			continue
		}

		if reflect.ValueOf(val).IsZero() {
			continue
		}

		fields = append(fields, dbTag)
		holders = append(holders, "?")
		updates = append(updates, fmt.Sprintf("%s = ?", dbTag))
		args = append(args, val)
	}

	var query string
	switch op {
	case Insert:
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			table, strings.Join(fields, ", "), strings.Join(holders, ", "))
	case Update:
		if pkValue == nil {
			return "", nil, "", fmt.Errorf("missing primary key value")
		}
		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
			table, strings.Join(updates, ", "), pk)
		args = append(args, pkValue)
	default:
		return "", nil, "", fmt.Errorf("unsupported operation")
	}

	// Build debug version of query
	debug := interpolateQuery(query, args)

	return query, args, debug, nil
}

func interpolateQuery(query string, args []any) string {
	var b strings.Builder
	argIndex := 0

	for i := 0; i < len(query); i++ {
		if query[i] == '?' && argIndex < len(args) {
			val := args[argIndex]
			argIndex++

			switch v := val.(type) {
			case string:
				b.WriteString("'" + strings.ReplaceAll(v, "'", "''") + "'")
			default:
				b.WriteString(fmt.Sprintf("%v", v))
			}
		} else {
			b.WriteByte(query[i])
		}
	}
	return b.String()
}

func MYSQLDate() string {
	t := time.Now()
	return fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
}

func ScapeString(t string) string {
	desc := strings.ReplaceAll(t, "'", "")
	desc = strings.ReplaceAll(desc, "\"", "")
	return desc
}
