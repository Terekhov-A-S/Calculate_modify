package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Калькулятор ===")
	fmt.Print("Введите имя пользователя: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Здесь будем проверять пользователя
	fmt.Printf("Привет, %s!\n", username)
}
