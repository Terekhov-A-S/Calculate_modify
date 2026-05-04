package system

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UserData хранит данные пользователя
type UserData struct {
	Username string
	FilePath string // путь к файлу с историей
}

// CheckOrCreateUser проверяет, есть ли пользователь, если нет - создаёт
func CheckOrCreateUser(username string) (UserData, error) {
	// Получаем путь к папке, где находится исполняемый файл
	execPath, err := os.Executable()
	if err != nil {
		return UserData{}, fmt.Errorf("не удалось определить путь к программе: %v", err)
	}

	// Получаем папку, где лежит .exe
	appDir := filepath.Dir(execPath)

	// Создаём полные пути к файлам и папкам
	usersFilePath := filepath.Join(appDir, "data", "users.txt")
	historyDir := filepath.Join(appDir, "data", "history")

	// Создаём ВСЕ необходимые папки (data и history)
	err = os.MkdirAll(historyDir, 0755)
	if err != nil {
		return UserData{}, fmt.Errorf("не удалось создать папку history: %v", err)
	}

	// Открываем файл users.txt (создаём если нет)
	file, err := os.OpenFile(usersFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return UserData{}, fmt.Errorf("не удалось открыть файл пользователей: %v", err)
	}
	defer file.Close()

	// Читаем всех существующих пользователей
	scanner := bufio.NewScanner(file)
	userExists := false

	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == username {
			userExists = true
			break
		}
	}

	// Если пользователя нет - добавляем
	if !userExists {
		_, err = file.WriteString(username + "\n")
		if err != nil {
			return UserData{}, fmt.Errorf("не удалось добавить пользователя: %v", err)
		}
		fmt.Printf("✅ Новый пользователь '%s' создан!\n", username)
	} else {
		fmt.Printf("👋 С возвращением, '%s'!\n", username)
	}

	// Создаём путь к файлу истории пользователя
	userHistoryFile := filepath.Join(historyDir, username+".txt")

	return UserData{
		Username: username,
		FilePath: userHistoryFile,
	}, nil
}
