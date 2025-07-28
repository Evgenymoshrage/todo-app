package json_storage

import (
	"encoding/json"
	"errors"
	"os"
	manager "todo-app/internal/todo/manager"
	task "todo-app/internal/todo/task"
)

const mainFile = "tasks.json"

// LoadJSON импортирует задачи из указанного JSON-файла и перезаписывает основной
func LoadJSON(srcPath string) error {
	// Проверка наличия файла
	_, err := os.Stat(srcPath)
	if os.IsNotExist(err) {
		return errors.New("Файл импорта не найден: " + srcPath)
	} else if err != nil {
		return err
	}

	// Чтение и разбор JSON
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	var imported []task.Task
	err = json.Unmarshal(data, &imported)
	if err != nil {
		return err
	}

	// Передаю задачи в менеджер и сохраняю
	manager.SetTasks(imported)
	return manager.SaveToFile(mainFile)
}

// SaveJSON экспортирует задачи из основного файла в указанный путь
func SaveJSON(destPath string) error {

	tasks := manager.GetTasks()
	// Загружаю из основного файла
	err := manager.LoadFromFile(mainFile)
	if err != nil {
		return err
	}

	// Сериализация и сохранение
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, data, 0644)
}
