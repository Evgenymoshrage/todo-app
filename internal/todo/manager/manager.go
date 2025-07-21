package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	task "todo-app/internal/todo/task"
)

// tasks — внутренний список задач
var tasks []task.Task // Импортирую структуру из пакета "task.go"
var nextID = 1

func LoadFromFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []task.Task{}
			nextID = 1 // новый файл — начинаем с 1
			return nil
		}
		return fmt.Errorf("Ошибка чтения файла: %w", err)
	}

	if err := json.Unmarshal(content, &tasks); err != nil {
		return fmt.Errorf("Ошибка разбора JSON: %w", err)
	}

	// Вычисляем nextID как max(tasks.ID) + 1
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	nextID = maxID + 1
	return nil
}

// Add добавляет задачу
func Add(description string) {
	newTask := task.Task{
		ID:          nextID,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, newTask)
	nextID++
}

func SaveToFile(filename string) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("Ошибка сериализации JSON: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

// List возвращает задачи по фильтру
func List(filter string) []task.Task {
	var result []task.Task
	switch filter {
	case "done":
		for _, task := range tasks {
			if task.Done {
				result = append(result, task)
			}
		}
	case "pending":
		for _, task := range tasks {
			if !task.Done {
				result = append(result, task)
			}
		}
	default: // "all"
		result = tasks
	}
	// Сортирую по ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

// Complete отмечает задачу выполненной
func Complete(id int) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			return nil
		}
	}
	return errors.New("Задача с указанным ID не найдена")
}

// Delete удаляет задачу по ID
func Delete(id int) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("Задача с указанным ID не найдена")
}

// Функция пересчитывает ID
func SetTasks(newTasks []task.Task) {
	tasks = newTasks

	maxID := 0
	for _, task := range tasks { // Нахожу максимальный существующий ID в новом списке.
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	nextID = maxID + 1
}

// Функция возвращает текущий список всех задач
func GetTasks() []task.Task {
	return tasks // Возвращаю tasks для других модулей
}
