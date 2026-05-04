package user

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// UserData хранит данные пользователя
type UserData struct {
	Username string
	FilePath string // путь к файлу с историей
}

// CheckOrCreateUser проверяет, есть ли пользователь, если нет - создаёт
func CheckOrCreateUser(username string) (UserData, error) {
	// Проверяем, существует ли файл users.txt
	usersFile := "User/users.txt"

	// Открываем файл (создаём, если нет)
	file, err := os.OpenFile(usersFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return UserData{}, err
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
			return UserData{}, err
		}
		fmt.Printf("✅ Новый пользователь '%s' создан!\n", username)
	} else {
		fmt.Printf("👋 С возвращением, '%s'!\n", username)
	}

	// Создаём папку history, если её нет
	err = os.MkdirAll("history", 0755)
	if err != nil {
		return UserData{}, err
	}

	return UserData{
		Username: username,
		FilePath: "history/" + username + ".txt",
	}, nil
}
