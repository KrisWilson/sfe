package cmd

import (
	"fmt"
	"runtime"
	"sfe/client"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var compileInfoCmd = &cobra.Command{
	Use:   "compile-info",
	Short: "Information about compiler parameters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Compiler Parameters")
		fmt.Println("GO OS: " + runtime.GOOS)
		fmt.Println("GO ARCH: " + runtime.GOARCH)
		fmt.Println("GO COMPILER: " + runtime.Compiler)
		fmt.Println("GO VERSION: " + runtime.Version())
		fmt.Println("GO BUILD: " + runtime.Version())
	},
}

var downloadDirCmd = &cobra.Command{
	Use:   "dirdownload",
	Short: "[{path}] Pobiera rekurencyjnie folder spod danej ścieżki",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		client.ConnectServer()
		timing := time.Now()
		wg.Add(1)
		var files uint
		var bytes uint64
		client.DownloadDir(strings.Join(args, " "), "", &wg, &files, &bytes)
		wg.Wait()
		fmt.Println("Job ended, time passed: " + strconv.FormatFloat(time.Now().Sub(timing).Seconds(), 'f', 2, 64) + "s, [ " + client.BytesShortener(bytes) + " ] in " + strconv.Itoa(int(files)) + " files")
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
		timing := time.Now()
		var bytes uint64
		client.DownloadFile(strings.Join(path[:len(path)-1], "/"), path[len(path)-1], "", &wg, &bytes)
		wg.Wait()
		fmt.Println("Job ended, time passed: " + strconv.FormatFloat(time.Now().Sub(timing).Seconds(), 'f', 2, 64))
	},
}

//goland:noinspection DuplicatedCode
var uploadFileCmd = &cobra.Command{
	Use:   "upload",
	Short: "[{path}] [{uploadPath}] Wysyła dany plik do folderu uzytkownika",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("[SFE] Expected at least 1 argument (got 0)")
			return
		}
		if len(args) == 1 {
			args = append(args, "")
		}
		client.ConnectServer()
		var wg sync.WaitGroup
		wg.Add(1)
		timing := time.Now()
		client.UploadFile(args[0], args[1], &wg)
		wg.Wait()
		fmt.Println("Job ended, time passed: " + strconv.FormatFloat(time.Now().Sub(timing).Seconds(), 'f', 2, 64))
	},
}

//goland:noinspection DuplicatedCode
var uploadDirCmd = &cobra.Command{
	Use:   "dirupload",
	Short: "[{dirPath}] [{uploadPath}] Przesyła folder rekurencyjnie na serwer do folderu użytkownika",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("[SFE] Expected at least 1 argument (got 0)")
			return
		}

		if len(args) == 1 {
			args = append(args, "")
		}

		client.ConnectServer()
		var wg sync.WaitGroup
		wg.Add(1)
		timing := time.Now()
		client.UploadDir(args[0], args[1], &wg)
		wg.Wait()
		fmt.Println("Job ended, time passed: " + strconv.FormatFloat(time.Now().Sub(timing).Seconds(), 'f', 2, 64))
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
