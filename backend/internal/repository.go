package internal

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"traceroute-optimised/internal/domain"

	_ "github.com/lib/pq"
)

type DataRepository struct {
	db *sql.DB
}

func NewDataRepository() (*DataRepository, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("PSQL_HOST")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)

	psqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("NewDataRepository error: %v", err)
	}

	if err := psqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DataRepository{
		db: psqlDB,
	}, nil
}

func (d DataRepository) GetData(ctx context.Context, dataKey string, keyName string) (*domain.DataModel, error) {
	var query string
	if keyName == "key" {
		query = `
        SELECT key, value
        FROM data
        WHERE key = $1
    `
	} else {
		query = `
        SELECT key, value
        FROM data
        WHERE other_key = $1
    `
	}

	var dataModel domain.DataModel

	err := d.db.QueryRowContext(ctx, query, dataKey).
		Scan(&dataModel.Key, &dataModel.Value)

	if err != nil {
		fmt.Println(dataModel)

		return nil, err
	}
	return &dataModel, nil
}

func (d DataRepository) GetBulkData(ctx context.Context, dataKeys []string, keyName string) (*[]domain.DataModel, error) {
	if len(dataKeys) == 0 {
		return nil, nil
	}

	// Определяем имя колонки
	col := "key"
	if keyName == "other_key" {
		col = "other_key"
	}

	placeholders := make([]string, len(dataKeys))
	args := make([]interface{}, len(dataKeys))
	for i, k := range dataKeys {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = k
	}

	query := fmt.Sprintf(`
        SELECT key, value
        FROM data
        WHERE %s IN (%s)
    `, col, strings.Join(placeholders, ","))

	// Выполняем запрос
	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.DataModel
	for rows.Next() {
		var m domain.DataModel
		if err := rows.Scan(&m.Key, &m.Value); err != nil {
			return nil, err
		}
		results = append(results, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &results, nil
}

func (d DataRepository) UpdateData(ctx context.Context, newData domain.DataModel, keyName string) (*domain.DataModel, error) {
	var query string
	if keyName == "key" {
		query = `
        UPDATE data
        SET value = $1
        WHERE key = $2
        RETURNING key, value
    `
	} else {
		query = `
        UPDATE data
        SET value = $1
        WHERE other_key = $2
        RETURNING key, value
    `
	}

	var dm domain.DataModel

	err := d.db.QueryRowContext(ctx, query, newData.Value, newData.Key).
		Scan(&dm.Key, &dm.Value)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}
