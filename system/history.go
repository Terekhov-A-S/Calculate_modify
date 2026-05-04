package system

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AddRecord добавляет запись в историю пользователя
func AddRecord(filePath, expression string, result float64) error {
	// Убеждаемся, что папка для файла существует
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать папку для истории: %v", err)
	}

	// Открываем файл для добавления (флаг O_APPEND)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл истории: %v", err)
	}
	defer file.Close()

	// Формируем запись с временем
	timestamp := time.Now().Format("02.01.2006 15:04:05")
	record := fmt.Sprintf("[%s] %s = %.2f\n", timestamp, expression, result)

	// Записываем в файл
	_, err = file.WriteString(record)
	return err
}

// ShowHistory показывает историю пользователя
func ShowHistory(filePath string, lastLines int) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("История пока пуста")
			return nil
		}
		return err
	}
	defer file.Close()

	// Читаем весь файл
	stat, _ := file.Stat()
	data := make([]byte, stat.Size())
	file.Read(data)

	if len(data) == 0 {
		fmt.Println("История пока пуста")
		return nil
	}

	fmt.Println("\n=== История вычислений ===")
	fmt.Print(string(data))
	fmt.Println("=========================")

	return nil
}
