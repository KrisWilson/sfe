package cmd

import (
	"fmt"
	"log"
	"os"
	"sfe/client"
	"sfe/listener"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "SmallFileExchanger",
	Short: "Mała aplikacja do wymiany plików, ponieważ FTP/SMB setupowanie na chwile jest irytujące tylko aby użyć dla paru plików >:(",
	Long:  "Mały server/client app do wymiany plików,\npozwala on bez męczarni z pełnymi protokołami FTP/SMB itp przesyłać dane na inne urządzenia",

	Run: func(cmd *cobra.Command, args []string) {
		client.Run()
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Uruchamia się w trybie listen - server",
	Long:  "Uruchamia się jako server, zezwala na odbieranie nadchodzących żądań",

	Run: func(cmd *cobra.Command, args []string) {
		// Pobierz wartość opcji --port
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Starting server on port %d\n", port)
		listener.Host(port)
	},
}

func init() {

	// Server command
	serverCmd.Flags().IntP("port", "p", 7068, "Server port")
	serverCmd.Aliases = []string{"host", "listen"}
	rootCmd.AddCommand(serverCmd)

	// Database commands
	rootCmd.AddCommand(addUserCmd)
	rootCmd.AddCommand(rmUserCmd)
	rootCmd.AddCommand(viewUsersCmd)

	// Client commands
	rootCmd.AddCommand(exploreDirCmd)
	rootCmd.AddCommand(downloadFileCmd)
	rootCmd.AddCommand(downloadDirCmd)
	rootCmd.AddCommand(uploadFileCmd)
	rootCmd.AddCommand(uploadDirCmd)

	// Dodaj komendę zatrzymującą serwer HTTP
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the HTTP server",
		Long:  "Stops the HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Stopping the HTTP server...")
			// Implementacja zatrzymywania serwera HTTP
		},
	}

	rootCmd.AddCommand(stopCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
