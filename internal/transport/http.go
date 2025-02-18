package transport

import (
	"encoding/json"
	"github.com/irongollem/goland-api-project/internal/model"
	"github.com/irongollem/goland-api-project/internal/todo"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(s *todo.Service) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		var b []byte
		var err = error(nil)
		if query != "" {
			foundItems, err := s.FindByItem(query)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			b, err = json.Marshal(foundItems)
		} else {
			allItems, err := s.GetAll()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			b, err = json.Marshal(allItems)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		var t model.TodoItem
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = s.Add(&t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	})

	mux.HandleFunc("DELETE /todo/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/todo/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = s.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	})

	return &Server{mux: mux}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.mux)
}
