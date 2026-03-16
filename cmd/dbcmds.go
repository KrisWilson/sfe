package cmd

import (
	"sfe/listener"

	"github.com/spf13/cobra"
)

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Dodaj uzytkownika do bazy danych",
	Long:  "Uruchamia menu dodawania użytkownika do bazy danych",
	//Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		listener.ConfigDB(2)
	},
}

var rmUserCmd = &cobra.Command{
	Use:   "rm",
	Short: "Usuń użytkownika z bazy danych",
	Long:  "Uruchamia menu usuwania uzytkownika z bazy danych",
	//Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		listener.ConfigDB(3)
	},
}

var viewUsersCmd = &cobra.Command{
	Use:   "view",
	Short: "Wyswietl baze danych",
	Long:  "Wyświetla bazę danych z użytkownikami",
	Run: func(cmd *cobra.Command, args []string) {
		listener.ConfigDB(1)
	},
}
