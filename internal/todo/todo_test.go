package todo_test

import (
	"api-project/internal/model"
	"api-project/internal/todo"
	"context"
	"errors"
	"reflect"
	"testing"
)

type MockDB struct {
	items []model.TodoItem
}

func (m *MockDB) InsertItem(_ context.Context, item *model.TodoItem) error {
	m.items = append(m.items, *item)
	return nil
}

func (m *MockDB) GetAllItems(_ context.Context) ([]model.TodoItem, error) {
	return m.items, nil
}

func (m *MockDB) DeleteItem(_ context.Context, id int) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	return errors.New("item not found")
}

func TestService_FindByItem(t *testing.T) {
	tests := []struct {
		name       string
		toDosToAdd []model.TodoItem
		query      string
		want       []model.TodoItem
	}{
		{
			name:       "given a todo of shop and a search of sh, it should get shop back",
			toDosToAdd: []model.TodoItem{{Task: "shop", Status: false}},
			query:      "sh",
			want:       []model.TodoItem{{ID: 0, Task: "shop", Status: false}},
		}, {
			name:       "given a todo of shop and a search of SH, it should get shop back",
			toDosToAdd: []model.TodoItem{{Task: "shop", Status: false}},
			query:      "SH",
			want:       []model.TodoItem{{ID: 0, Task: "shop", Status: false}},
		}, {
			name:       "given a todo of shop and a search of foo, it should get nothing back",
			toDosToAdd: []model.TodoItem{{Task: "shop", Status: false}},
			query:      "foo",
			want:       []model.TodoItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := todo.NewService(&MockDB{})
			for _, toAdd := range tt.toDosToAdd {
				err := s.Add(&toAdd)
				if err != nil {
					return
				}
			}

			if got, _ := s.FindByItem(tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Add(t *testing.T) {
	tests := []struct {
		name       string
		toDosToAdd []model.TodoItem
		want       bool
	}{
		{"Adding a todo of shop should return no error", []model.TodoItem{{Task: "shop", Status: false}}, false},
		{"Adding a todo of shop twice should return an error", []model.TodoItem{{Task: "shop", Status: false}, {Task: "shop", Status: true}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := todo.NewService(&MockDB{})
			var err error
			for _, toAdd := range tt.toDosToAdd {
				err = s.Add(&toAdd)
			}

			if (err != nil) != tt.want {
				t.Errorf("Add() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestService_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		todosToAdd []model.TodoItem
		want       []model.TodoItem
	}{

		{name: "Given a todo of shop, it should return a todo of shop", todosToAdd: []model.TodoItem{{Task: "shop", Status: false}}, want: []model.TodoItem{{ID: 0, Task: "shop", Status: false}}},
		{name: "Given an empty todo list, it should return an empty todo list", todosToAdd: []model.TodoItem{}, want: []model.TodoItem{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := todo.NewService(&MockDB{items: []model.TodoItem{}})
			for _, toAdd := range tt.todosToAdd {
				_ = s.Add(&toAdd)
			}

			if got, _ := s.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		idToDel    int
		todosToAdd []model.TodoItem
		wantErr    bool
	}{
		{name: "Given a todo of shop and an id of 1, it should delete the todo", idToDel: 0, todosToAdd: []model.TodoItem{{Task: "shop", Status: false}}, wantErr: false},
		{name: "Given a todo of shop and an id of 2, it should return an error", idToDel: 1, todosToAdd: []model.TodoItem{{Task: "shop", Status: false}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := todo.NewService(&MockDB{})
			for _, toAdd := range tt.todosToAdd {
				_ = s.Add(&toAdd)
			}

			if err := s.Delete(tt.idToDel); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
