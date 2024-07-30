// В данном пакете хранятся хендлеры для корректной обработки запросов клиента
package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/db"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/models"
)

// статически создаем ошибку для обработки ситуации, когда в JSON от клиента указана
// задача с ID, который уже присутствует в мапе
var ErrExistingID = errors.New("задача с таким ID уже присутствует в списке дел")

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
	// записываем сериализованные в формат JSON данные в тело ответа
	_, err = w.Write(resp)
	if err != nil {
		// если возникает ошибка при записи сериализованных данных в ответ,
		// логируем ее в консоль и возвращаем клиенту статус "500 Internal Server Error"
		fmt.Printf("Ошибка при записи сериализованных данных в тело ответа: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	// проверяем, существует ли уже в мапе задача с тем же ID, который был получен от клиента
	_, found := db.Tasks[task.ID]
	// если задача с таким ID уже есть в мапе с делами, возвращаем ошибку ErrExistingID и статус "400 Bad Request"
	if found {
		http.Error(w, ErrExistingID.Error(), http.StatusBadRequest)
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
	// при ошибке возвращаем ответ "500 Internal Server Error"
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// в заголовок записываем тип контента: данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// записываем в заголовок статус OK
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		// если возникает ошибка при записи сериализованных данных в ответ,
		// логируем ее в консоль и возвращаем клиенту статус "500 Internal Server Error"
		fmt.Printf("Ошибка при записи сериализованных данных в тело ответа: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
