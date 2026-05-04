package main

import (
	"Calculate_modify/system"
	"bufio"
	"fmt"
	"math"
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
	userData, err := system.CheckOrCreateUser(username)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		fmt.Println("Нажмите Enter для выхода...")
		fmt.Scanln()
		return
	}

	fmt.Println("Доступные операции: +, -, *, /, % (остаток от деления)")
	fmt.Println("Примеры ввода: 2+2, 10-3, 4*5, 15/3, 10%3, 12^2")
	fmt.Println("Можно вводить с пробелами или без")
	fmt.Println("Для выхода введите 'exit'")
	fmt.Println("Для просмотра истории введите 'history'")
	fmt.Println()

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Удаляем ВСЕ пробелы из строки (чтобы "2 + 2" тоже работало)
		inputNoSpaces := strings.ReplaceAll(input, " ", "")

		if input == "exit" {
			fmt.Println("До свидания!")
			break
		}

		if input == "history" {
			err := system.ShowHistory(userData.FilePath, 0)
			if err != nil {
				fmt.Printf("Ошибка при чтении истории: %v\n", err)
			}
			continue
		}

		// Вычисляем (используем строку без пробелов)
		result, err := calculate(inputNoSpaces)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			continue
		}

		fmt.Printf("Результат: %.2f\n", result)

		// Сохраняем в историю исходный ввод (с пробелами, если они были)
		err = system.AddRecord(userData.FilePath, input, result)
		if err != nil {
			fmt.Printf("Не удалось сохранить историю: %v\n", err)
		}
	}
}

// calculate принимает строку без пробелов типа "2+2" или "10/3" или "10%3"
func calculate(input string) (float64, error) {
	// Ищем позицию оператора в строке
	var operator string
	var pos int

	// Проверяем наличие каждого оператора
	for i, ch := range input {
		switch ch {
		case '+', '-', '*', '/', '%', '^':
			operator = string(ch)
			pos = i
			break
		}
		if operator != "" {
			break
		}
	}

	if operator == "" {
		return 0, fmt.Errorf("не найден оператор (+, -, *, /, %%)")
	}

	leftStr := input[:pos]
	rightStr := input[pos+1:]

	a, err := strconv.ParseFloat(leftStr, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка в левом числе: %v", leftStr)
	}

	b, err := strconv.ParseFloat(rightStr, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка в правом числе: %v", rightStr)
	}

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
	case "^":
		return math.Pow(a, b), nil

	case "%":
		aInt := int(a)
		bInt := int(b)

		if bInt == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}

		if float64(aInt) != a || float64(bInt) != b {
			return 0, fmt.Errorf("оператор %% работает только с целыми числами")
		}

		return float64(aInt % bInt), nil
	default:
		return 0, fmt.Errorf("неизвестный оператор: %v", operator)
	}
}
