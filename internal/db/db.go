package db

import (
	"context"
	"fmt"
	"github.com/irongollem/goland-api-project/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(user, password, dbname, host string, port int) (*DB, error) {
	conStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	pool, err := pgxpool.Connect(context.Background(), conStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the db: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping the db: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) InsertItem(ctx context.Context, item *model.TodoItem) error {
	query := "INSERT INTO todo_items (task, status) VALUES ($1, $2)"
	_, err := db.pool.Exec(ctx, query, item.Task, item.Status)
	return err
}

func (db *DB) GetAllItems(ctx context.Context) ([]model.TodoItem, error) {
	query := "SELECT id, task, status FROM todo_items"
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.TodoItem
	for rows.Next() {
		var item model.TodoItem
		err := rows.Scan(&item.ID, &item.Task, &item.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (db *DB) GetItemByID(ctx context.Context, id int) (*model.TodoItem, error) {
	query := "SELECT id, task, status FROM todo_items WHERE id = $1"
	row := db.pool.QueryRow(ctx, query, id)

	var item model.TodoItem
	err := row.Scan(&item.ID, &item.Task, &item.Status)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (db *DB) UpdateItem(ctx context.Context, item *model.TodoItem) (*model.TodoItem, error) {
	query := "UPDATE todo_items SET task = $1, status = $2 WHERE id = $3"
	_, err := db.pool.Exec(ctx, query, item.Task, item.Status, item.ID)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (db *DB) DeleteItem(ctx context.Context, id int) error {
	query := "DELETE FROM todo_items WHERE id = $1"
	_, err := db.pool.Exec(ctx, query, id)
	return err
}

func (db *DB) Close() {
	db.pool.Close()
}
