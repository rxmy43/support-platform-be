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
	// Hapus setUUID karena kita menggunakan auto increment
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

	// Jika perlu mendapatkan ID yang di-generate, bisa diambil seperti ini
	id, err := result.LastInsertId()
	if err == nil {
		// Set ID ke struct jika diperlukan
		setID(entity, uint(id))
	}

	return nil
}

func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error {
	// Hapus setUUID
	setClauses := make([]string, 0)
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

func (r *BaseRepo[T]) Delete(ctx context.Context, id uint) error { // Ubah parameter menjadi uint
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.TableName)
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

// Fungsi untuk set ID ke struct (opsional, jika ingin mengisi ID setelah create)
func setID[T any](entity *T, id uint) {
	v := reflect.ValueOf(entity).Elem()
	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		// Cek tipe field ID
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
		if dbTag == "" || dbTag == "id" { // Skip field ID karena auto increment
			continue
		}
		columns = append(columns, dbTag)
		values = append(values, ":"+dbTag)
	}
	return columns, values
}

// Method tambahan yang berguna
func (r *BaseRepo[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.TableName)
	err := r.DB.GetContext(ctx, &entity, query, id)
	if err != nil {
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
