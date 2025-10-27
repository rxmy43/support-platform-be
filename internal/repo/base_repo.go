package repo

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type BaseRepo[T any] struct {
	DB        *sqlx.DB
	TableName string
}

func (r *BaseRepo[T]) Create(ctx context.Context, entity *T) error {
	columns, values := extractColumnsAndValues(entity)

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.TableName,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	result, err := r.DB.NamedExecContext(ctx, query, entity)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err == nil {
		setID(entity, uint(id))
	}

	return nil
}

func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error {
	setClauses := []string{}
	v := reflect.ValueOf(entity).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "id" {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s=:%s", dbTag, dbTag))
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=:id", r.TableName, strings.Join(setClauses, ", "))
	_, err := r.DB.NamedExecContext(ctx, query, entity)
	return err
}

func (r *BaseRepo[T]) Delete(ctx context.Context, id uint) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.TableName)
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

// Set ID ke struct
func setID[T any](entity *T, id uint) {
	v := reflect.ValueOf(entity).Elem()
	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		switch idField.Kind() {
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			idField.SetUint(uint64(id))
		}
	}
}

func extractColumnsAndValues[T any](entity *T) ([]string, []string) {
	v := reflect.ValueOf(entity).Elem()
	t := v.Type()
	columns := []string{}
	values := []string{}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "id" {
			continue
		}
		columns = append(columns, dbTag)
		values = append(values, ":"+dbTag)
	}
	return columns, values
}

func (r *BaseRepo[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", r.TableName)
	row := r.DB.QueryRowxContext(ctx, query, id)
	if err := row.StructScan(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepo[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	query := fmt.Sprintf("SELECT * FROM %s", r.TableName)
	err := r.DB.SelectContext(ctx, &entities, query)
	return entities, err
}
