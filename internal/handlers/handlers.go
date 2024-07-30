// В данном пакете хранятся хендлеры для корректной обработки запросов клиента
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/db"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/models"
)

// Параметры (для всех хендлеров)
// w - параметр типа http.ResponseWriter, т.е. ответ сервера клиенту
// r - параметр типа *http.Request, т.е. указатель на запрос от клиента серверу

// Хендлер GetTasks для получения всех задач из мапы tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем в JSON данные из мапы tasks
	resp, err := json.Marshal(db.Tasks)
	// если возникает ошибка, возвращаем статус "500 Internal Server Error"
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента: данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// записываем в заголовок статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// Хендлер PostTask для добавления новой задачи в мапу tasks на основании данных запроса от клиента
func PostTask(w http.ResponseWriter, r *http.Request) {
	// создаем экземпляр структуры Task, куда мы будем десериализовать данные, полученные из запроса клиента
	var task models.Task
	// создаем буфер для сохранения в нем "сырых" (сериализованных) данных из запроса
	var buf bytes.Buffer

	// читаем данные из тела запроса и записываем сериализованные данные в буфер
	_, err := buf.ReadFrom(r.Body)
	// при ошибке возвращаем клиенту ответ "400 Bad Request"
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем данные из буфера в переменную task
	// при ошибке возвращаем клиенту ответ "400 Bad Request"
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// сохраняем переменную task (полученную от клиента задачу) в мапу tasks (общий список задач)
	db.Tasks[task.ID] = task

	// в заголовок записываем тип контента: данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// при успешном запросе возвращаем статус 201 Created
	w.WriteHeader(http.StatusCreated)
}

// Хендлер GetTask для получения отдельной задачи по ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	// Функция chi.URLParam(r, "id") возвращает значение параметра из URL.
	id := chi.URLParam(r, "id")
	// проверяем наличие задачи в списке
	task, ok := db.Tasks[id]
	// если задача не найдена, возвращаем ответ "400 Bad Request" (согласно ТЗ)
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// если задача найдена, сериализуем ее данные и сохраняем в переменную resp
	resp, err := json.Marshal(task)
	// при ошибке возвращаем ответ "400 Bad Request"
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// в заголовок записываем тип контента: данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// записываем в заголовок статус OK
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Хендлер DeleteTask для удаления задачи по ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Функция chi.URLParam(r, "id") возвращает значение параметра из URL.
	id := chi.URLParam(r, "id")
	// проверяем наличие задачи в списке
	_, ok := db.Tasks[id]
	// если задача не найдена, возвращаем ответ "400 Bad Request" (согласно ТЗ)
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// если задача найдена, удаляем ее из мапы tasks
	delete(db.Tasks, id)

	// в заголовок записываем тип контента: данные в формате JSON (согласно ТЗ)
	w.Header().Set("Content-Type", "application/json")
	// записываем в заголовок статус OK
	w.WriteHeader(http.StatusOK)
}
