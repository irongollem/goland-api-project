package todo

import (
	"api-project/internal/model"
	"context"
	"errors"
	"fmt"
	"strings"
)

type DBManager interface {
	InsertItem(ctx context.Context, item *model.TodoItem) error
	GetAllItems(ctx context.Context) ([]model.TodoItem, error)
	DeleteItem(ctx context.Context, id int) error
}

type Service struct {
	db DBManager
}

func NewService(db DBManager) *Service {
	return &Service{db: db}
}

func (s *Service) Add(todo *model.TodoItem) error {
	items, err := s.GetAll()
	if err != nil {
		return fmt.Errorf("failed to read from DB: %w", err)
	}
	for _, item := range items {
		if item.Task == todo.Task {
			return errors.New("todo item already exist")
		}
	}

	err = s.db.InsertItem(context.Background(), todo)
	if err != nil {
		return fmt.Errorf("failed to insert item into DB: %w", err)
	}
	return nil
}

func (s *Service) FindByItem(query string) ([]model.TodoItem, error) {
	matchingTodos := make([]model.TodoItem, 0)
	lowerQuery := strings.ToLower(query)
	items, err := s.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read from DB: %w", err)
	}
	for _, todo := range items {
		if strings.Contains(strings.ToLower(todo.Task), lowerQuery) {
			matchingTodos = append(matchingTodos, todo)
		}
	}
	return matchingTodos, nil
}

func (s *Service) Delete(id int) error {
	err := s.db.DeleteItem(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to delete item from DB: %w", err)
	}
	return nil
}

func (s *Service) GetAll() ([]model.TodoItem, error) {
	items, err := s.db.GetAllItems(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get all items from DB: %w", err)
	}
	return items, nil
}
