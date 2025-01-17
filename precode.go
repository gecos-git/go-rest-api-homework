package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для получения всех задач.
// Конечная точка /tasks.
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)

	// При ошибке сервер должен вернуть статус 500 Internal Server Error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// При успешном запросе сервер должен вернуть статус 200 OK.
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик для отправки задачи на сервер.
// Конечная точка /tasks.
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	// При ошибке сервер должен вернуть статус 400 Bad Request.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {

		// При ошибке сервер должен вернуть статус 400 Bad Request.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")

	// При успешном запросе сервер должен вернуть статус 201 Created.
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
// Конечная точка /tasks/{id}.
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	// Если такого ID нет, верните соответствующий статус.
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(task)

	// В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// При успешном выполнении запроса сервер должен вернуть статус 200 OK.
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик удаления задачи по ID.
// Конечная точка /tasks/{id}.
func delTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	// В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")

	// При успешном выполнении запроса сервер должен вернуть статус 200 OK.
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// Метод GET.
	r.Get("/tasks", getTasks)
	// Метод POST.
	r.Post("/tasks", postTask)
	// Метод GET.
	r.Get("/tasks/{id}", getTask)
	// Метод DELETE.
	r.Delete("/tasks/{id}", delTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
