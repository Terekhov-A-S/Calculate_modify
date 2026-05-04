package main

import (
	history "Calculate_modify/system"
	user "Calculate_modify/system"
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

	// Защита от пустого имени
	for username == "" {
		fmt.Print("Имя не может быть пустым! Введите имя: ")
		usernameInput, _ := reader.ReadString('\n')
		username = strings.TrimSpace(usernameInput)
	}

	// Проверяем/создаём пользователя
	userData, err := user.CheckOrCreateUser(username)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		fmt.Println("Нажмите Enter для выхода...")
		fmt.Scanln()
		return
	}

	fmt.Println("Доступные операции: +, -, *, /, % (остаток от деления)")
	fmt.Println("Поддерживаются сложные выражения: 2+3*4, (2+3)*4, 10/2+3")
	fmt.Println("Поддерживаются унарные операторы: +5, -3, (+5), -(+3)")
	fmt.Println("Можно вводить с пробелами или без")
	fmt.Println("Для выхода введите 'exit'")
	fmt.Println("Для просмотра истории введите 'history'")
	fmt.Println()

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Защита от пустого ввода
		if input == "" {
			fmt.Println("⚠️ Ввод не может быть пустым! Введите выражение.")
			continue
		}

		// Защита от слишком длинного ввода
		if len(input) > 100 {
			fmt.Println("⚠️ Слишком длинное выражение (максимум 100 символов)")
			continue
		}

		// Обработка команд
		if input == "exit" {
			fmt.Println("До свидания!")
			break
		}

		if input == "history" {
			err := history.ShowHistory(userData.FilePath, 0)
			if err != nil {
				fmt.Printf("Ошибка при чтении истории: %v\n", err)
			}
			continue
		}

		// Удаляем ВСЕ пробелы из строки
		inputNoSpaces := strings.ReplaceAll(input, " ", "")

		// Дополнительная проверка: не пустая ли строка после удаления пробелов
		if inputNoSpaces == "" {
			fmt.Println("⚠️ Введены только пробелы! Введите выражение.")
			continue
		}

		// Вычисляем
		result, err := evaluateExpression(inputNoSpaces)
		if err != nil {
			fmt.Printf("❌ Ошибка: %v\n", err)
			continue
		}

		fmt.Printf("✅ Результат: %.2f\n", result)

		// Сохраняем в историю исходный ввод
		err = history.AddRecord(userData.FilePath, input, result)
		if err != nil {
			fmt.Printf("⚠️ Не удалось сохранить историю: %v\n", err)
		}
	}
}

// evaluateExpression вычисляет сложное выражение с приоритетом операций
func evaluateExpression(expr string) (float64, error) {
	// Удаляем все пробелы
	expr = strings.ReplaceAll(expr, " ", "")

	if expr == "" {
		return 0, fmt.Errorf("пустое выражение")
	}

	// Проверка на допустимые символы
	for _, ch := range expr {
		isValid := (ch >= '0' && ch <= '9') ||
			ch == '.' ||
			ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' ||
			ch == '(' || ch == ')'

		if !isValid {
			return 0, fmt.Errorf("недопустимый символ '%c'", ch)
		}
	}

	// Нормализуем выражение: обрабатываем унарные операторы
	expr = normalizeUnaryOperators(expr)

	// Проверка на сбалансированность скобок
	if err := checkParentheses(expr); err != nil {
		return 0, err
	}

	// Вычисляем выражение
	result, err := calculateExpression(expr)
	if err != nil {
		return 0, err
	}

	// Проверка на слишком большой результат
	if result > 1e15 || result < -1e15 {
		return 0, fmt.Errorf("результат слишком большой")
	}

	return result, nil
}

// normalizeUnaryOperators заменяет унарные операторы на специальные маркеры
func normalizeUnaryOperators(expr string) string {
	result := make([]rune, 0, len(expr))
	runes := []rune(expr)

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		// Проверяем, является ли + унарным (в начале выражения или после открывающей скобки или после другого оператора)
		if ch == '+' {
			// Унарный плюс, если:
			// 1. Это первый символ
			// 2. Или перед ним открывающая скобка
			// 3. Или перед ним другой оператор
			isUnary := i == 0 ||
				runes[i-1] == '(' ||
				runes[i-1] == '+' ||
				runes[i-1] == '-' ||
				runes[i-1] == '*' ||
				runes[i-1] == '/' ||
				runes[i-1] == '%'

			if isUnary {
				// Унарный плюс просто игнорируем (не добавляем в результат)
				continue
			}
		}

		// Проверяем, является ли - унарным
		if ch == '-' {
			// Унарный минус, если:
			// 1. Это первый символ
			// 2. Или перед ним открывающая скобка
			// 3. Или перед ним другой оператор
			isUnary := i == 0 ||
				runes[i-1] == '(' ||
				runes[i-1] == '+' ||
				runes[i-1] == '-' ||
				runes[i-1] == '*' ||
				runes[i-1] == '/' ||
				runes[i-1] == '%'

			if isUnary {
				// Унарный минус заменяем на специальный маркер "~"
				result = append(result, '~')
				continue
			}
		}

		result = append(result, ch)
	}

	return string(result)
}

// checkParentheses проверяет правильность расстановки скобок
func checkParentheses(expr string) error {
	balance := 0
	for i, ch := range expr {
		switch ch {
		case '(':
			balance++
		case ')':
			balance--
			if balance < 0 {
				return fmt.Errorf("закрывающая скобка без открывающей на позиции %d", i)
			}
		}
	}
	if balance != 0 {
		return fmt.Errorf("незакрытых скобок: %d", balance)
	}
	return nil
}

// calculateExpression вычисляет выражение с учётом приоритетов
func calculateExpression(expr string) (float64, error) {
	// Сначала обрабатываем скобки
	expr, err := processParentheses(expr)
	if err != nil {
		return 0, err
	}

	// Разбиваем на токены (числа и операторы)
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}

	// Вычисляем с учётом приоритетов
	return evaluateTokens(tokens)
}

// processParentheses рекурсивно обрабатывает скобки
func processParentheses(expr string) (string, error) {
	for {
		// Ищем самую внутреннюю пару скобок
		openIdx := -1
		closeIdx := -1

		for i, ch := range expr {
			if ch == '(' {
				openIdx = i
			} else if ch == ')' && openIdx != -1 {
				closeIdx = i
				break
			}
		}

		// Если скобок нет - выходим
		if openIdx == -1 {
			break
		}

		// Вычисляем выражение внутри скобок
		innerExpr := expr[openIdx+1 : closeIdx]
		result, err := evaluateSimpleExpression(innerExpr)
		if err != nil {
			return "", err
		}

		// Заменяем скобки на результат
		expr = expr[:openIdx] + fmt.Sprintf("%g", result) + expr[closeIdx+1:]
	}

	return expr, nil
}

// evaluateSimpleExpression вычисляет выражение без скобок
func evaluateSimpleExpression(expr string) (float64, error) {
	// Сначала нормализуем унарные операторы внутри выражения
	expr = normalizeUnaryOperators(expr)

	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}
	return evaluateTokens(tokens)
}

// tokenize разбивает строку на числа и операторы с учётом унарного минуса
func tokenize(expr string) ([]string, error) {
	var tokens []string
	var currentNumber strings.Builder

	for i := 0; i < len(expr); i++ {
		ch := rune(expr[i])

		// Обработка унарного минуса (маркер ~)
		if ch == '~' {
			// Начинаем отрицательное число
			currentNumber.WriteRune('-')
			continue
		}

		// Если это цифра или точка
		if (ch >= '0' && ch <= '9') || ch == '.' {
			currentNumber.WriteRune(ch)
		} else if ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' {
			// Оператор
			if currentNumber.Len() > 0 {
				tokens = append(tokens, currentNumber.String())
				currentNumber.Reset()
			}
			tokens = append(tokens, string(ch))
		} else {
			return nil, fmt.Errorf("неизвестный символ: %c", ch)
		}
	}

	// Добавляем последнее число
	if currentNumber.Len() > 0 {
		tokens = append(tokens, currentNumber.String())
	}

	// Проверяем, что выражение не заканчивается оператором
	if len(tokens) > 0 && isOperator(tokens[len(tokens)-1]) {
		return nil, fmt.Errorf("выражение не может заканчиваться оператором")
	}

	// Валидация: не идут ли два оператора подряд
	for i := 0; i < len(tokens)-1; i++ {
		if isOperator(tokens[i]) && isOperator(tokens[i+1]) {
			return nil, fmt.Errorf("два оператора подряд: %s и %s", tokens[i], tokens[i+1])
		}
	}

	return tokens, nil
}

// isOperator проверяет, является ли строка оператором
func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/" || token == "%"
}

// evaluateTokens вычисляет выражение по токенам с приоритетом операций
func evaluateTokens(tokens []string) (float64, error) {
	if len(tokens) == 0 {
		return 0, fmt.Errorf("пустое выражение")
	}

	// Сначала обрабатываем умножение, деление и остаток
	result, err := evaluateHighPriority(tokens)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// evaluateHighPriority обрабатывает *, /, %
func evaluateHighPriority(tokens []string) (float64, error) {
	// Копируем токены в новый слайс
	newTokens := make([]string, len(tokens))
	copy(newTokens, tokens)

	i := 1
	for i < len(newTokens) {
		if newTokens[i] == "*" || newTokens[i] == "/" || newTokens[i] == "%" {
			if i-1 < 0 || i+1 >= len(newTokens) {
				return 0, fmt.Errorf("неверный формат выражения")
			}

			// Получаем левый и правый операнды
			leftStr := newTokens[i-1]
			rightStr := newTokens[i+1]

			// Преобразуем в числа
			left, err := strconv.ParseFloat(leftStr, 64)
			if err != nil {
				return 0, fmt.Errorf("ошибка в числе: %v", leftStr)
			}

			right, err := strconv.ParseFloat(rightStr, 64)
			if err != nil {
				return 0, fmt.Errorf("ошибка в числе: %v", rightStr)
			}

			// Вычисляем результат
			var result float64
			switch newTokens[i] {
			case "*":
				result = left * right
			case "/":
				if right == 0 {
					return 0, fmt.Errorf("деление на ноль")
				}
				result = left / right
			case "%":
				leftInt := int(left)
				rightInt := int(right)
				if float64(leftInt) != left || float64(rightInt) != right {
					return 0, fmt.Errorf("оператор %% работает только с целыми числами")
				}
				if rightInt == 0 {
					return 0, fmt.Errorf("деление на ноль")
				}
				result = float64(leftInt % rightInt)
			}

			// Заменяем три элемента на один результат
			newTokens = append(newTokens[:i-1], fmt.Sprintf("%g", result))
			newTokens = append(newTokens, newTokens[i+2:]...)
			// Остаёмся на той же позиции
		} else {
			i++
		}
	}

	// Теперь обрабатываем сложение и вычитание
	result, err := evaluateLowPriority(newTokens)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// evaluateLowPriority обрабатывает + и -
func evaluateLowPriority(tokens []string) (float64, error) {
	if len(tokens) == 0 {
		return 0, fmt.Errorf("пустое выражение")
	}

	// Начинаем с первого числа
	result, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка в числе: %v", tokens[0])
	}

	i := 1
	for i < len(tokens) {
		if i+1 >= len(tokens) {
			return 0, fmt.Errorf("неверный формат выражения")
		}

		operator := tokens[i]
		rightStr := tokens[i+1]

		right, err := strconv.ParseFloat(rightStr, 64)
		if err != nil {
			return 0, fmt.Errorf("ошибка в числе: %v", rightStr)
		}

		switch operator {
		case "+":
			result += right
		case "-":
			result -= right
		default:
			return 0, fmt.Errorf("неизвестный оператор: %v", operator)
		}

		i += 2
	}

	return result, nil
}
