package cmd

import (
	"sfe/client"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// TODO: uploadFile, uploadDir
// TODO: changePass from clientside

var downloadDirCmd = &cobra.Command{
	Use:   "dirdownload",
	Short: "[{path}] Pobiera rekurencyjnie folder spod danej ścieżki",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		client.ConnectServer()
		wg.Add(1)
		client.DownloadDir(strings.Join(args, " "), "", &wg)
		wg.Wait()
	},
}

var downloadFileCmd = &cobra.Command{
	Use:   "download",
	Short: "[{path}] Pobiera plik spod danej ścieżki",
	Run: func(cmd *cobra.Command, args []string) {
		// eg ./sfe download path/to/file/foo.bar
		client.ConnectServer()
		path := strings.Split(strings.Join(args, ""), "/")
		var wg sync.WaitGroup
		wg.Add(1)
		client.DownloadFile(strings.Join(path[:len(path)-1], "/"), path[len(path)-1], "", &wg)
		wg.Wait()
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
