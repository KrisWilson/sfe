package client

import (
	"bufio"
	"fmt"
	"os"

	"sfe/listener"
	"sfe/settings"
)

func readKey() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, _ := reader.ReadRune()
	return char
}

func Run() {

	fmt.Println("<<< SFE - Small File Exchanger >>>")
	fmt.Println("[1] Connect to Server")
	fmt.Println("[2] Host a server")
	fmt.Println("[3] Show config")
	fmt.Println("[4] Exit")
	//fmt.Println("Your choice: \"" + string(input) + "\"")
	input := readKey()

	switch string(input) {
	case "1":

	case "2":
		listener.Host(8670)

	case "3":
		settings.Load()
		Run()

	case "4":
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice.")
	}

}
