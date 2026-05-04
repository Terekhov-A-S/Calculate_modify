package History

import (
	"fmt"
	"os"
	"time"
)

// AddRecord добавляет запись в историю пользователя
func AddRecord(filePath, expression string, result float64) error {
	// Открываем файл для добавления (флаг O_APPEND)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Формируем запись с временем
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	record := fmt.Sprintf("[%s] %s = %.2f\n", timestamp, expression, result)

	// Записываем в файл
	_, err = file.WriteString(record)
	return err
}

// ShowHistory показывает последние N записей (опционально)
func ShowHistory(filePath string, lastLines int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Читаем весь файл (для простоты, потом можно сделать эффективнее)
	stat, _ := file.Stat()
	data := make([]byte, stat.Size())
	file.Read(data)

	fmt.Println("\n=== История вычислений ===")
	fmt.Print(string(data))
	fmt.Println("=========================")

	return nil
}
