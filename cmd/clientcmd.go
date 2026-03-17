package cmd

import (
	"sfe/client"
	"strings"

	"github.com/spf13/cobra"
)

// TODO: w downloadDir, uploadFile, uploadDir
// TODO: changePass from clientside

var downloadDirCmd = &cobra.Command{
	Use:   "dirdownload",
	Short: "[{path}] Pobiera rekurencyjnie folder spod danej ścieżki",
	Run: func(cmd *cobra.Command, args []string) {
		client.ConnectServer()
		client.DownloadDir(strings.Join(args, " "), "")
	},
}

var downloadFileCmd = &cobra.Command{
	Use:   "download",
	Short: "[{path}] Pobiera plik spod danej ścieżki",
	Run: func(cmd *cobra.Command, args []string) {
		// eg ./sfe download path/to/file/foo.bar
		client.ConnectServer()
		path := strings.Split(strings.Join(args, ""), "/")
		client.DownloadFile(strings.Join(path[:len(path)-1], "/"), path[len(path)-1], "")
	},
}

var exploreDirCmd = &cobra.Command{
	Use:   "explore",
	Short: "[{path}] Pobiera informacje o zadanym folderze i zwraca informacje",
	Long:  "[{path}] Pobiera informacje o zadanym folderze i zwraca informacje",
	Run: func(cmd *cobra.Command, args []string) {
		client.ConnectServer()

		client.ExploreDir(strings.Join(args, " "))
	},
}
