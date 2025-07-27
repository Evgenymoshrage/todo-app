package main

import (
	"flag"
	"fmt"
	"os"
	csv_storage "todo-app/internal/storage/csv_storage"
	json_storage "todo-app/internal/storage/json_storage"
	manager "todo-app/internal/todo/manager"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" { //Если аргументов нет или --help или короткая версия -h,
		manager.PrintHelp() //то показываю справку
		return
	}

	manager.LoadFromFile("tasks.json") //Загружаю задачи из файла при запуске.

	// os.Args — список аргументов командной строки.
	// os.Args[0] — имя исполняемого файла (например, todo).
	// os.Args[1] — подкоманда (add, list, complete и т.д.).
	// Если нет подкоманды — показываем ошибку и выходим.

	switch os.Args[1] { //Проверяю, какая команда передана (add, list, complete, ...).
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)   //Создаю новый набор флагов для команды add
		desc := addCmd.String("desc", "", "Описание задачи") //Регистрирую флаг --desc для описания задачи
		addCmd.Parse(os.Args[2:])

		if *desc == "" { //Проверка на наличие описания задачи
			fmt.Println("Описание задачи обязательно. Используйте: todo add --desc=\"...\"")
			return
		}

		manager.Add(*desc)
		manager.SaveToFile("tasks.json")
		fmt.Println("Задача добавлена.")

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)                        //Создаю новый набор флагов для команды list
		filter := listCmd.String("filter", "all", "Фильтр задач: all/done/pending") //Регистрирую флаг --filter
		listCmd.Parse(os.Args[2:])

		tasks := manager.List(*filter)
		for _, task := range tasks {
			status := "Не выполнено"
			if task.Done {
				status = "Выполнено"
			}
			fmt.Printf("[%s] #%d: %s\n", status, task.ID, task.Description)
		}

	case "complete":
		completeCmd := flag.NewFlagSet("complete", flag.ExitOnError) //Создаю новый набор флагов для команды complete
		id := completeCmd.Int("id", 0, "ID задачи для выполнения")   //Регистрирую флаг --id
		completeCmd.Parse(os.Args[2:])

		if *id <= 0 { //Проверка на корректность ID
			fmt.Println("Укажите корректный ID: todo complete --id=1")
			return
		}

		err := manager.Complete(*id)
		if err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			manager.SaveToFile("tasks.json")
			fmt.Println("Задача выполнена.")
		}

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError) //Создаю новый набор флагов для команды delete
		id := deleteCmd.Int("id", 0, "ID задачи для удаления")   //Регистрирую флаг --id
		deleteCmd.Parse(os.Args[2:])

		if *id <= 0 {
			fmt.Println("Укажите корректный ID: todo delete --id=1")
			return
		}

		err := manager.Delete(*id)
		if err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			manager.SaveToFile("tasks.json")
			fmt.Println("Задача удалена.")
		}

	case "import": //Подкоманда import, два флага: --file, --format
		importCmd := flag.NewFlagSet("import", flag.ExitOnError)
		file := importCmd.String("file", "", "Файл импорта")
		format := importCmd.String("format", "json", "Формат: json/csv")
		importCmd.Parse(os.Args[2:])

		if *file == "" {
			fmt.Println("Укажите файл: todo import --file=tasks.json --format=json")
			return
		}

		var err error

		switch *format { //В зависимости от формата загружаю из соответствующего файла
		case "json":
			err = json_storage.LoadJSON(*file)
		case "csv":
			err = csv_storage.LoadCSV(*file)
		default:
			fmt.Println("Неверный формат импорта. Должно быть json/csv")
			return
		}

		if err != nil {
			fmt.Println("Ошибка импорта:", err)
		} else {
			fmt.Println("Импорт завершён.")
		}

	case "export": //Подкоманда export, два флага: --file, --format
		exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
		file := exportCmd.String("file", "", "Файл экспорта")
		format := exportCmd.String("format", "json", "Формат: json/csv")
		exportCmd.Parse(os.Args[2:])

		if *file == "" {
			fmt.Println("Укажите файл: todo export --file=tasks.csv --format=csv")
			return
		}

		var err error

		switch *format { //В зависимости от формата загружаю из соответствующего файла
		case "json":
			err = json_storage.SaveJSON(*file)
		case "csv":
			err = csv_storage.SaveCSV(*file)
		default:
			fmt.Println("Неверный формат экспорта. Должно быть json/csv")
			return
		}

		if err != nil {
			fmt.Println("Ошибка экспорта:", err)
		} else {
			fmt.Println("Экспорт завершён.")
		}

	default: //Если введена неизвестная команда — ошибка + список доступных.
		fmt.Println("Неизвестная команда:", os.Args[1])
		fmt.Println("Доступные команды: add, list, complete, delete, import, export")
		os.Exit(1)
	}

	fmt.Println("Работа программы завершена.")
}
