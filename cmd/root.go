package cmd

import (
	"fmt"
	"log"
	"os"
	"sfe/client"
	"sfe/listener"

	"github.com/spf13/cobra"
)

// TODO: Dodaj więcej CLI komend wykorzystujących już istniejące metody w innych plikach
// TODO: exploreDir, downloadFile, downloadDir, uploadFile, uploadDir
// TODO: changePass from clientside

var rootCmd = &cobra.Command{
	Use:   "SmallFileExchanger",
	Short: "Mała aplikacja do wymiany plików, ponieważ FTP/SMB setupowanie na chwile jest irytujące tylko aby użyć dla paru plików >:(",
	Long:  "Mały server/client app do wymiany plików, pozwala on bez męczarni z pełnymi protokołami FTP/SMB itp przesyłać dane na inne urządzenia",

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

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Dodaj uzytkownika do bazy danych",
	Long:  " ",
	//Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		listener.ConfigDB(2)
	},
}

var rmUserCmd = &cobra.Command{
	Use:   "rm",
	Short: "Usuń użytkownika z bazy danych",
	Long:  " ",
	//Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		listener.ConfigDB(3)
	},
}

var viewUsersCmd = &cobra.Command{
	Use:   "view",
	Short: "Wyswietl baze danych",
	Long:  " ",
	Run: func(cmd *cobra.Command, args []string) {
		listener.ConfigDB(1)
	},
}

func init() {
	serverCmd.Flags().IntP("port", "p", 7068, "Server port")
	serverCmd.Aliases = []string{"host", "listen"}
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(addUserCmd)
	rootCmd.AddCommand(rmUserCmd)
	rootCmd.AddCommand(viewUsersCmd)

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
