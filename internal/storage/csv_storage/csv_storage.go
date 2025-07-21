package csv_storage

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	manager "todo-app/internal/todo/manager"
	task "todo-app/internal/todo/task"
)

const mainFile = "tasks.json"

// Функция импортирует задачи из CSV и сохраняет их в основной JSON
func LoadCSV(path string) error {
	// Открываю файл
	file, err := os.Open(path)
	if err != nil {
		return errors.New("Не удалось открыть CSV-файл: " + err.Error())
	}
	defer file.Close()

	// Считываю все записи
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return errors.New("Не удалось прочитать CSV: " + err.Error())
	}

	// Пропускаю заголовок
	if len(records) <= 1 {
		return errors.New("CSV не содержит данных (только заголовок или пустой)")
	}
	records = records[1:]

	// Преобразываю строки
	var imported []task.Task
	for _, record := range records {
		if len(record) != 3 {
			continue // Пропускаем некорректные строки
		}

		id, err1 := strconv.Atoi(record[0])
		description := record[1]
		done, err2 := strconv.ParseBool(record[2])

		if err1 != nil || err2 != nil {
			continue // Пропускаем строки с ошибками
		}

		imported = append(imported, task.Task{
			ID:          id,
			Description: description,
			Done:        done,
		})
	}

	manager.SetTasks(imported)
	return manager.SaveToFile(mainFile)
}

// Функция экспортирует текущие задачи из менеджера в CSV-файл
func SaveCSV(path string) error {
	err := manager.LoadFromFile(mainFile)
	if err != nil {
		return err
	}

	tasks := manager.GetTasks()

	// 1. Создать или перезаписать файл
	file, err := os.Create(path)
	if err != nil {
		return errors.New("не удалось создать CSV-файл: " + err.Error())
	}
	defer file.Close()

	// 2. Инициализировать writer и записать заголовок
	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"ID", "Description", "Done"})
	if err != nil {
		return errors.New("не удалось записать заголовок в CSV: " + err.Error())
	}

	// 3. Записать задачи
	for _, t := range tasks {
		record := []string{
			strconv.Itoa(t.ID),
			t.Description,
			strconv.FormatBool(t.Done),
		}
		if err := writer.Write(record); err != nil {
			return errors.New("ошибка записи строки в CSV: " + err.Error())
		}
	}

	// 4. Завершить запись
	writer.Flush()

	if err := writer.Error(); err != nil {
		return errors.New("ошибка при завершении записи CSV: " + err.Error())
	}

	return nil
}
