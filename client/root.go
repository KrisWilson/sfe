package client

import (
	"bufio"
	"fmt"
	"os"
)

func Run() {
	// Dodaj obsługę menu interaktywnego
	fmt.Println("Welcome to My Application!")
	fmt.Println("1. View version")
	fmt.Println("2. Exit")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	switch input {
	case "1":
		// Wywołaj polecenie view-version
	case "2":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice.")
	}

}
