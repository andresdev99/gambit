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
	"products": "Prod_Id",
	"category": "Categ_Id",
}

type modelTypes interface {
	models.Product | models.Category
}

func BuildSelectQuery[T modelTypes](
	filter T,
	table string,
	orderBy string,
	orderDir string,
	limit int,
	offset int,
) (string, []any, string, error) {
	t := reflect.TypeOf(filter)
	v := reflect.ValueOf(filter)

	if t.Kind() != reflect.Struct {
		return "", nil, "", fmt.Errorf("filter must be a struct")
	}

	var selectFields []string
	var whereClauses []string
	var args []any

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		selectFields = append(selectFields, dbTag)

		val := v.Field(i).Interface()
		if reflect.ValueOf(val).IsZero() {
			continue
		}

		if field.Type.Kind() == reflect.String {
			whereClauses = append(whereClauses, fmt.Sprintf("%s LIKE ?", dbTag))
			args = append(args, "%"+val.(string)+"%")
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", dbTag))
			args = append(args, val)
		}
	}

	// Start building query
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(selectFields, ", "), table)

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Validate orderBy using struct's db tags
	if orderBy != "" {
		validCols := map[string]bool{}
		for i := 0; i < t.NumField(); i++ {
			if col := t.Field(i).Tag.Get("db"); col != "" {
				validCols[col] = true
			}
		}

		if validCols[orderBy] {
			orderDir = strings.ToUpper(orderDir)
			if orderDir != "ASC" && orderDir != "DESC" {
				orderDir = "ASC"
			}
			query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)
		} else {
			fmt.Printf("⚠️  Ignoring invalid orderField: %q\n", orderBy)
		}
	}

	// Add LIMIT / OFFSET
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	debug := interpolateQuery(query, args)
	return query, args, debug, nil
}

func BuildDeleteQuery[T modelTypes](model T, table string) (string, []any, string, error) {
	pk, ok := modelProperties[table]
	if !ok {
		return "", nil, "", fmt.Errorf("no primary key defined for table %q", table)
	}

	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	if v.Kind() != reflect.Struct {
		return "", nil, "", fmt.Errorf("model must be a struct")
	}

	var pkValue any
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == pk {
			pkValue = v.Field(i).Interface()
			break
		}
	}

	if pkValue == nil || reflect.ValueOf(pkValue).IsZero() {
		return "", nil, "", fmt.Errorf("primary key value is missing or zero")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, pk)
	debug := interpolateQuery(query, []any{pkValue})

	return query, []any{pkValue}, debug, nil
}

func BuildInsertUpdateQuery[mType modelTypes](model mType, table string, op Operation) (string, []any, string, error) {
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
