package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/handlers"
)

func main() {
	//создаем новый роутер
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", handlers.GetTasks)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTask`
	r.Post("/tasks", handlers.PostTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks/{id}", handlers.GetTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTask`
	r.Delete("/tasks/{id}", handlers.DeleteTask)

	// запускаем сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
