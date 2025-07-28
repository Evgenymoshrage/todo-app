package main_test

import (
	"os"
	"testing"
	"todo-app/internal/storage/csv_storage"
	"todo-app/internal/storage/json_storage"
	"todo-app/internal/todo/manager"
)

func TestAddTask(t *testing.T) { // Начинается с Test, чтобы Go её увидел как тест
	manager.ClearTasks() // Функция, очищающая все задачи для теста

	manager.Add("Что-то сделать") // Тестирую добавление задачи
	tasks := manager.List("all")

	if len(tasks) != 1 { // Проверка на количество задач
		t.Fatalf("Ожидалась 1 задача, получено %d", len(tasks))
	}

	if tasks[0].Description != "Изучить тесты в Go" { // Проверка на описание задачи
		t.Errorf("Ожидалось описание 'Изучить тесты в Go', получено '%s'", tasks[0].Description)
	}
}

func TestDeleteTask(t *testing.T) { // Начинается с Test, чтобы Go её увидел как тест
	manager.ClearTasks()          // Функция, очищающая все задачи для теста
	manager.Add("Что-то сделать") // Добавляю одну задачу

	tasks := manager.List("pending") // Проверяю добавилась ли задача
	if len(tasks) == 0 {
		t.Fatal("Не удалось добавить задачу")
	}

	manager.Delete(1) // Тестирую удаление задачи

	tasks = manager.List("pending") // Обновляю список задач

	if len(tasks) != 0 { // Проверка на количество задач
		t.Fatalf("Ожидалась 0 задач, получено на %d задач", len(tasks))
	}
}

func TestCompleteTask(t *testing.T) { // Начинается с Test, чтобы Go её увидел как тест
	manager.ClearTasks() // Функция, очищающая все задачи для теста

	manager.Add("Что-то сделать") // Добавляю одну задачу

	tasks := manager.List("pending") // Проверяю добавилась ли задача
	if len(tasks) == 0 {
		t.Fatal("Не удалось добавить задачу")
	}

	manager.Complete(1)          // Тестирую выполнение задачи
	tasks = manager.List("done") // Обновляю список задач

	if tasks[0].Done != true { // Проверка на описание задачи
		t.Errorf("Ожидалась выполненная задача.")
	}
}

func TestListTaskWithFilters(t *testing.T) {
	manager.ClearTasks() // Очищаю список перед тестом

	// Добавляю три задачи
	manager.Add("Задача 1") // ID 1
	manager.Add("Задача 2") // ID 2
	manager.Add("Задача 3") // ID 3

	// Отмечаю одну задачу выполненной
	manager.Complete(1)

	// Проверка фильтра "all"
	allTasks := manager.List("all")
	if len(allTasks) != 3 {
		t.Errorf("Ожидалось 3 задачи при фильтре 'all', получено %d", len(allTasks))
	}

	// Проверка фильтра "done"
	doneTasks := manager.List("done")
	if len(doneTasks) != 1 {
		t.Errorf("Ожидалась 1 выполненная задача при фильтре 'done', получено %d", len(doneTasks))
	} else if !doneTasks[0].Done || doneTasks[0].Description != "Задача 2" {
		t.Errorf("Неверная выполненная задача при фильтре 'done'")
	}

	// Проверка фильтра "pending"
	pendingTasks := manager.List("pending")
	if len(pendingTasks) != 2 {
		t.Errorf("Ожидалось 2 невыполненные задачи при фильтре 'pending', получено %d", len(pendingTasks))
	}
}

func TestLoadJSON_FileNotFound(t *testing.T) { // Проверка открытия несузествующего файла
	err := json_storage.LoadJSON("nonexistent.json")
	if err == nil {
		t.Error("Ожидалась ошибка при загрузке несуществующего файла, но её не произошло")
	}
}

func TestSaveAndLoadJSON_Success(t *testing.T) {
	manager.ClearTasks() // Очищаю список перед тестом

	manager.Add("Задача 1") // Добавляю две задачи
	manager.Add("Задача 2")

	tmpFile := "test_tasks.json" // Экспортиркю в файл

	defer os.Remove(tmpFile) // Удаляю файл после теста

	err := json_storage.SaveJSON(tmpFile) // Сохраняю задачи в файл
	if err != nil {
		t.Fatalf("Ошибка при сохранении задач: %v", err)
	}

	manager.ClearTasks() // Очищаю память перед загрузкой

	err = json_storage.LoadJSON(tmpFile) // Загружаю задачи из файла
	if err != nil {
		t.Fatalf("Ошибка при загрузке задач: %v", err)
	}

	test_tasks := manager.List("all")
	if len(test_tasks) != 2 {
		t.Errorf("Ожидалось 2 задачи после загрузки, получено %d", len(test_tasks))
	}
}

func TestLoadCSV_FileNotFound(t *testing.T) { // Проверка открытия несузествующего файла
	err := csv_storage.LoadCSV("nonexistent.csv")
	if err == nil {
		t.Error("Ожидалась ошибка при загрузке несуществующего файла, но её не произошло")
	}
	if !os.IsNotExist(err) {
		t.Errorf("Ожидалась ошибка os.IsNotExist, но получено: %v", err)
	}
}

func TestSaveAndLoadCSV_Success(t *testing.T) {
	manager.ClearTasks() // Очищаю список перед тестом

	manager.Add("Задача 1") // Добавляю две задачи
	manager.Add("Задача 2")

	tmpFile := "test_tasks.csv" // Экспортирую в файл

	defer os.Remove(tmpFile) // Удаляю файл после теста

	err := csv_storage.SaveCSV(tmpFile) // Сохраняю задачи в файл
	if err != nil {
		t.Fatalf("Ошибка при сохранении задач в CSV: %v", err)
	}

	manager.ClearTasks() // Очищаю память перед загрузкой

	err = csv_storage.LoadCSV(tmpFile) // Загружаю задачи из файла
	if err != nil {
		t.Fatalf("Ошибка при загрузке задач из CSV: %v", err)
	}

	tasks := manager.List("all")
	if len(tasks) != 2 {
		t.Fatalf("Ожидалось 2 задачи, получено %d", len(tasks))
	}
}
