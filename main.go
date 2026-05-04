package main

import (
	history "Calculate_modify/History"
	user "Calculate_modify/User"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Калькулятор ===")
	fmt.Print("Введите имя пользователя: ")
	usernameInput, _ := reader.ReadString('\n')
	username := strings.TrimSpace(usernameInput)

	// Проверяем/создаём пользователя
	userData, err := user.CheckOrCreateUser(username)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Println("Доступные операции: +, -, *, /")
	fmt.Println("Формат ввода: 2 + 2")
	fmt.Println("Для выхода введите 'exit'")
	fmt.Println("Для просмотра истории введите 'history'")
	fmt.Println()

	// Основной цикл калькулятора
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("До свидания!")
			break
		}

		if input == "history" {
			// Показываем историю
			err := history.ShowHistory(userData.FilePath, 0)
			if err != nil {
				fmt.Println("История пуста или ошибка чтения")
			}
			continue
		}

		// Вычисляем выражение
		result, err := calculate(input)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			continue
		}

		// Выводим результат
		fmt.Printf("Результат: %.2f\n", result)

		// Сохраняем в историю
		err = history.AddRecord(userData.FilePath, input, result)
		if err != nil {
			fmt.Printf("Не удалось сохранить историю: %v\n", err)
		}
	}
}

// calculate парсит и вычисляет простое арифметическое выражение
func calculate(input string) (float64, error) {
	// Разделяем пробелами: "2 + 2" -> ["2", "+", "2"]
	parts := strings.Fields(input)
	if len(parts) != 3 {
		return 0, fmt.Errorf("неверный формат. Используйте: \"число оператор число\"")
	}

	// Парсим числа
	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("первый аргумент не число: %v", parts[0])
	}

	b, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, fmt.Errorf("второй аргумент не число: %v", parts[2])
	}

	// Выполняем операцию
	operator := parts[1]
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("неизвестный оператор: %v", operator)
	}
}
