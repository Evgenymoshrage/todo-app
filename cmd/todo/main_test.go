package main_test

import (
	"testing"
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
