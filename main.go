package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	csv_storage "todo-app/internal/storage/csv_storage"
	json_storage "todo-app/internal/storage/json_storage"
	manager "todo-app/internal/todo/manager"
)

var answer string // Ввожу переменные для функций работы с пользователем

var ( // Ввожу флаги
	addFlag    = flag.String("add", "", "Добавить новую задачу")
	listFlag   = flag.Bool("list", false, "Показать список задач")
	filterFlag = flag.String("filter", "all", "Фильтр задач: all/done/pending")
	completeID = flag.Int("complete", 0, "Отметить задачу выполненной по ID")
	deleteID   = flag.Int("delete", 0, "Удалить задачу по ID")
	exportFlag = flag.String("export", "", "Экспорт задач в файл (указать имя)")
	exportFmt  = flag.String("exportfmt", "json", "Формат экспорта: json/csv")
	importFlag = flag.String("import", "", "Импорт задач из файла (указать имя)")
	importFmt  = flag.String("importfmt", "json", "Формат импорта: json/csv")
)

func main() {
	// Сначала работа с флагами, а потом работа с пользователем
	flag.Parse()

	// Загружаю задачи один раз
	if err := manager.LoadFromFile("tasks.json"); err != nil {
		fmt.Println("Ошибка загрузки задач:", err)
		os.Exit(1)
	}

	if *addFlag != "" { // Флаг добавления задачи
		manager.Add(*addFlag)
		if err := manager.SaveToFile("tasks.json"); err != nil {
			fmt.Println("Ошибка сохранения задач:", err)
		}
		fmt.Println("Задача добавлена.")
		return
	}

	if *listFlag { // Флаг выведения задач по фильтру
		tasks := manager.List(*filterFlag)
		for _, task := range tasks {
			status := "Не выполнено"
			if task.Done {
				status = "Выполнено"
			}
			fmt.Printf("[%s] #%d: %s\n", status, task.ID, task.Description)
		}
		return
	}

	if *completeID > 0 { // Флаг выполнения задачи
		if err := manager.Complete(*completeID); err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			manager.SaveToFile("tasks.json")
			fmt.Println("Задача выполнена.")
		}
		return
	}

	if *deleteID > 0 { // Флаг удаления задачи
		if err := manager.Delete(*deleteID); err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			manager.SaveToFile("tasks.json")
			fmt.Println("Задача удалена.")
		}
		return
	}

	if *exportFlag != "" { // Флаг экспорта задач
		switch strings.ToLower(*exportFmt) {
		case "json":
			err := json_storage.SaveJSON(*exportFlag)
			if err != nil {
				fmt.Println("Ошибка экспорта:", err)
			} else {
				fmt.Println("Экспорт завершён (JSON).")
			}
		case "csv":
			err := csv_storage.SaveCSV(*exportFlag)
			if err != nil {
				fmt.Println("Ошибка экспорта:", err)
			} else {
				fmt.Println("Экспорт завершён (CSV).")
			}
		default:
			fmt.Println("Неверный формат экспорта.")
		}
		return
	}

	if *importFlag != "" { // Флаг импорта задач
		switch strings.ToLower(*importFmt) {
		case "json":
			err := json_storage.LoadJSON(*importFlag)
			if err != nil {
				fmt.Println("Ошибка импорта:", err)
			} else {
				fmt.Println("Импорт завершён (JSON).")
			}
		case "csv":
			err := csv_storage.LoadCSV(*importFlag)
			if err != nil {
				fmt.Println("Ошибка импорта:", err)
			} else {
				fmt.Println("Импорт завершён (CSV).")
			}
		default:
			fmt.Println("Неверный формат импорта.")
		}
		return
	}

	// Работа с пользователем (сделал на базе моей программы работы с файлами,
	// интересно было доработать чуть управление)

	// Запрос функции добавления задачи
	fmt.Println("Вы хотите добавить новую задачу? (да/нет)")
	fmt.Scan(&answer)

	reader := bufio.NewReader(os.Stdin) // читаю данные из консоли с буферизацией
	reader.ReadString('\n')             // очищает лишний '\n' после Scan

	if strings.ToLower(answer) == "да" {
		err := manager.LoadFromFile("tasks.json")
		if err != nil {
			fmt.Println("Ошибка загрузки задач:", err)
			return
		}

		fmt.Print("Введите описание новой задачи: ")
		description, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			return
		}
		description = strings.TrimSpace(description)
		manager.Add(description)

		err = manager.SaveToFile("tasks.json")
		if err != nil {
			fmt.Println("Ошибка сохранения задач:", err)
			return
		}
		fmt.Println("Задача добавлена.")
	}

	// Запрос функции выведения задач
	fmt.Println("Вы хотите вывести список задач? (да/нет)")
	fmt.Scan(&answer)
	reader.ReadString('\n') // очищает лишний '\n' после Scan

	if strings.ToLower(answer) == "да" {

		err := manager.LoadFromFile("tasks.json")
		if err != nil {
			fmt.Println("Ошибка загрузки задач:", err)
			return
		}
		for { // Цикл чтобы при ошиибочном фильтре снова ввести фильтр
			fmt.Print("Укажите фильтр (all/done/pending): ")
			filter, _ := reader.ReadString('\n')
			filter = strings.TrimSpace(strings.ToLower(filter))

			tasks := manager.List(filter)

			if len(tasks) == 0 {
				fmt.Println("Нет задач по заданному фильтру. Попробуйте другой фильтр.")
				continue // возвращаемся к вводу фильтра
			}

			fmt.Println("Список задач:")
			for _, task := range tasks {
				status := "Не выполнено"
				if task.Done {
					status = "Выполнено"
				}
				fmt.Printf("[%s] #%d: %s\n", status, task.ID, task.Description)
			}
			fmt.Println("Список выведен.")
			break // выходим из цикла, если задачи найдены
		}
	}

	// Запрос функции выполнения задачи
	fmt.Println("Вы хотите отметить задачу выполненой? (да/нет)")
	fmt.Scan(&answer)
	reader.ReadString('\n') // очищает лишний '\n' после Scan

	if strings.ToLower(answer) == "да" {
		err := manager.LoadFromFile("tasks.json") // загружаю файл
		if err != nil {
			fmt.Println("Ошибка загрузки задач:", err)
			return
		}

		for { // Цикл чтобы при ошиибочном выборе вернуться
			fmt.Print("Укажите номер задачи, которую необходимо отметить выполненной: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			id, err := strconv.Atoi(idInput) // перевожу строку в число

			if err != nil {
				fmt.Println("Ошибка: введите корректное число.")
				continue // возвращаемся к запросу
			}

			err = manager.Complete(id) // выполняю функцию Complete
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue // пользователь может ввести другой ID
			}

			err = manager.SaveToFile("tasks.json") // сохраняю файл
			if err != nil {
				fmt.Println("Ошибка сохранения задач:", err)
				return
			}
			fmt.Printf("Задача %d отмечена выполненной.\n", id)
			break // выходим из цикла после успешного выполнения
		}
	}

	// Запрос функции удаления задачи
	fmt.Println("Вы хотите удалить задачу? (да/нет)")
	fmt.Scan(&answer)
	reader.ReadString('\n') // очищает лишний '\n' после Scan

	if strings.ToLower(answer) == "да" {

		err := manager.LoadFromFile("tasks.json") // загружаю файл
		if err != nil {
			fmt.Println("Ошибка загрузки задач:", err)
			return
		}

		for { // Цикл чтобы при ошиибочном выборе вернуться
			fmt.Print("Укажите номер задачи, которую необходимо удалить: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			id, err := strconv.Atoi(idInput) // перевожу строку в число
			if err != nil {
				fmt.Println("Ошибка: введите корректное число.")
				continue // возвращаемся к запросу
			}

			err = manager.Delete(id) // выполняю функцию Delete
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue // пользователь может ввести другой ID
			}

			err = manager.SaveToFile("tasks.json") // сохраняю файл
			if err != nil {
				fmt.Println("Ошибка сохранения задач:", err)
				return
			}
			fmt.Printf("Задача %d удалена.\n", id)
			break // выходим из цикла после успешного выполнения
		}
	}

	// Запрос функции экспорта задачи
	fmt.Println("Вы хотите экспортировать задачу? (да/нет)")
	fmt.Scan(&answer)
	reader.ReadString('\n') // очищает лишний '\n' после Scan

	if answer == "да" {
		fmt.Print("В каком формате экспортировать? (json/csv): ")
		format, _ := reader.ReadString('\n')
		format = strings.TrimSpace(strings.ToLower(format))

		fmt.Print("Укажите имя файла (например, export.csv): ")
		filename, _ := reader.ReadString('\n')
		filename = strings.TrimSpace(filename)

		switch format {
		case "json":
			err := json_storage.SaveJSON(filename)
			if err != nil {
				fmt.Println("Ошибка экспорта в JSON:", err)
			} else {
				fmt.Println("Экспорт в JSON выполнен успешно.")
			}
		case "csv":
			err := csv_storage.SaveCSV(filename)
			if err != nil {
				fmt.Println("Ошибка экспорта в CSV:", err)
			} else {
				fmt.Println("Экспорт в CSV выполнен успешно.")
			}
		default:
			fmt.Println("Неверный формат.")
		}
	}

	// Запрос функции импорта задачи
	fmt.Println("Вы хотите импортировать задачу? (да/нет)")
	fmt.Scan(&answer)
	reader.ReadString('\n') // очищает лишний '\n' после Scan

	if answer == "да" {
		fmt.Print("В каком формате импортировать? (json/csv): ")
		format, _ := reader.ReadString('\n')
		format = strings.TrimSpace(strings.ToLower(format))

		fmt.Print("Укажите имя файла (например, import.csv): ")
		filename, _ := reader.ReadString('\n')
		filename = strings.TrimSpace(filename)

		switch format {
		case "json":
			err := json_storage.LoadJSON(filename)
			if err != nil {
				fmt.Println("Ошибка импорта из JSON:", err)
			} else {
				fmt.Println("Импорт из JSON выполнен успешно.")
			}
		case "csv":
			err := csv_storage.LoadCSV(filename)
			if err != nil {
				fmt.Println("Ошибка импорта из CSV:", err)
			} else {
				fmt.Println("Импорт из CSV выполнен успешно.")
			}
		default:
			fmt.Println("Неверный формат.")
		}
	}

	fmt.Println("Работа программы завершена.")
}
