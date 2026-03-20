package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sfe/listener"
	"sfe/settings"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

func readKey() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, _ := reader.ReadRune()
	return char
}

var token string
var config settings.Config
var oldState *term.State

func ErrorLog(err error, prefix string) {
	file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		_, err = file.WriteString(prefix + " " + err.Error() + "\n")
		if err != nil {
		}
		err = file.Close()
		if err != nil {
		}
	}
}

func BytesShortener(in uint64) string {
	if in < 1000 {
		return strconv.Itoa(int(in))
	} else if in > 1000 && in < 1000000 {
		return strconv.FormatFloat(float64(in)/1000, 'f', 2, 32) + " KB"

	} else if in > 1000000 && in < 1000000000 {
		return strconv.FormatFloat(float64(in)/1000000, 'f', 2, 32) + " MB"
	} else if in > 1000000000 {
		return strconv.FormatFloat(float64(in)/1000000000, 'f', 2, 32) + " GB"
	}
	return strconv.Itoa(int(in/1000000000000)) + " TB"

}

func ExploreDir(dir string) []byte {
	data := []byte("")
	//dir = url.PathEscape(dir)
	dir = url.QueryEscape(dir)
	req, err := http.NewRequest(http.MethodGet, "http://"+config.ConnectIP+":"+strconv.Itoa(config.ClientPort)+"/explore?path=/"+dir+"/", bytes.NewBuffer(data))
	if req != nil {
		req.Header.Set("Token", token)
	}
	if err != nil {
		ErrorLog(err, "[Http Request Error]")
		return []byte("Unauthorized")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		ErrorLog(err, "[Http Client]")
		return []byte("Unauthorized")
		//panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println("\n\033[36m" + "Folder path: /" + dir + "/\u001B[0m\r")

	// print raw json
	//fmt.Println("\033[31m" + string(bodyBytes) + "\u001B[0m\r")

	var filesJson []listener.FileJSON
	err = json.Unmarshal(bodyBytes, &filesJson)
	if err != nil {
		fmt.Println("[Response] Folder can't be read", "\r")
		for _, b := range bodyBytes {
			fmt.Printf("%c", b)
		}
		fmt.Print("\n")
		ErrorLog(err, "[JSON Parser]")
		return []byte("err")
	}
	var largestFile int64 //do sformatowania później listingu plików...
	for _, file := range filesJson {
		if file.Size > largestFile {
			largestFile = file.Size
		}
	}
	//min number of space, to dont crash strings.Repeat(" ", pos - 4)
	posNeeded := len(strconv.Itoa(int(largestFile)))
	if posNeeded < 4 {
		posNeeded = 4
	}
	fmt.Printf("\u001B[36mType \tDate Modified\t\t%s \tname\r\u001B[0m\n\r", "Size"+strings.Repeat(" ", posNeeded-4))
	for _, file := range filesJson {
		//	fmt.Printf("Name: %s\nType: %s\nSize: %d bytes\nDate Modified: %s\n", file.Name, file.Type, file.Size, file.DateModified)

		fmt.Printf("%s \t%s \t%s \t%s\r\n",
			file.Type,
			file.DateModified,
			strconv.Itoa(int(file.Size))+strings.Repeat(" ", posNeeded-len(strconv.Itoa(int(file.Size)))),
			file.Name)
	}
	fmt.Println("\r")
	return bodyBytes
}

func DownloadFile(dir string, filename string, downloadDir_ string, wg *sync.WaitGroup, bytesDownload *uint64) {
	defer wg.Done()
	//dir = url.PathEscape(dir)
	dir = url.QueryEscape(dir)
	filename = url.QueryEscape(filename)
	var downloadDir string
	if len(downloadDir_) == 0 {
		downloadDir = config.DownloadDir + "/"
	} else {
		downloadDir = config.DownloadDir + "/" + downloadDir_ + "/"
	}
	err := os.Mkdir(downloadDir, os.ModePerm)
	if err != nil {
		// its ok, po prostu folder istnieje, Albo ma jakiś mentalny problem tj. nie moze pisac w folderze
	}
	data := []byte("")
	if len(dir) > 0 {
		dir = dir + "/"
	}

	fmt.Println("[Client] Pobieranie \u001B[33m" + filename + "\u001B[0m do folderu " + downloadDir + filename + "....\r")
	req, err := http.NewRequest(http.MethodGet, "http://"+config.ConnectIP+":"+strconv.Itoa(config.ClientPort)+"/explore?path=/"+dir+"&file="+filename, bytes.NewBuffer(data))
	if req != nil {
		req.Header.Set("Token", token)
	}

	if err != nil {
		ErrorLog(err, "[Http Client]")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		ErrorLog(err, "[Http Client]")
		return
		//panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	bodyBytes, err := io.ReadAll(resp.Body)

	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		ErrorLog(err, "[Mkdir]")
		// do nothing - najprawdopodobniej istnieje już taki folder
	}

	err = os.WriteFile(downloadDir+"/"+filename, bodyBytes, os.ModePerm)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("[Client] \u001B[31mPobieranie niepowiodło się "+filename+" => ", err, "\r\u001B[0m")
		ErrorLog(err, "[Client]")
	} else {
		fmt.Println("[Client] \u001B[36mPobieranie powiodło się "+filename+" [ "+BytesShortener(uint64(len(bodyBytes)))+" ]", "\r\u001B[0m")
		*bytesDownload = uint64(len(bodyBytes)) + *bytesDownload
	}

	//fmt.Println("\033[31m" + string(bodyBytes) + "\u001B[0m\r")

	//fmt.Println("[Client] Zakonczone połączenie\r")
}

func DownloadDir(dir string, downloadDir string, wg *sync.WaitGroup, filesDownload *uint, bytesDownload *uint64) {
	defer wg.Done()
	var wgInside sync.WaitGroup
	list := ExploreDir(dir)
	if len(dir) == 0 {
		dir = "./"
	}
	if len(downloadDir) == 0 {
		downloadDir = ""
	}

	//fmt.Println("Dir: "+dir, "\nDownloadDir: "+downloadDir)
	var filesJson []listener.FileJSON
	err := json.Unmarshal(list, &filesJson)
	if err != nil {
		fmt.Println("Folder can't be read", "\r")
		ErrorLog(err, "[JSON Parser [Dir]]")
	}
	for _, file := range filesJson {
		if file.Type == "File" {
			wgInside.Add(1)
			*filesDownload = *filesDownload + 1
			go DownloadFile(dir, file.Name, downloadDir, &wgInside, bytesDownload)
		} else {
			err := os.MkdirAll(config.DownloadDir+"/"+downloadDir+"/"+file.Name, os.ModePerm)
			if err != nil {
				fmt.Println("Folder can't be created (", err, ")\r")
			}
			wgInside.Add(1)
			go DownloadDir(dir+"/"+file.Name, downloadDir+"/"+file.Name, &wgInside, filesDownload, bytesDownload)
		}
	}
	wgInside.Wait()
	return
}

func UploadFile(filename string, uploadPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(uploadPath) == 0 {
		uploadPath = ""
	}
	data := []byte("")
	var err error
	if len(filename) != 0 {
		data, err = os.ReadFile(filename)
		if err != nil {
			fmt.Println("[Client] File can't be read", "\r")
			return
		}
		buff := strings.Split(filename, "/")
		filename = buff[len(buff)-1]
	} else {
		filename = "."
	}
	req, err := http.NewRequest(http.MethodPut, "http://"+config.ConnectIP+":"+strconv.Itoa(config.ClientPort)+"/upload?filename="+filename+"&uploadpath="+uploadPath, bytes.NewBuffer(data))
	if req != nil {
		req.Header.Set("Token", token)
	}
	client := &http.Client{}
	fmt.Println("[Client] Wysyłanie \u001B[33m"+filename+"\u001B[0m do folderu "+uploadPath, "\r")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		//panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("[Client] Something went wrong " + resp.Status)
	} else if filename == "." {
		fmt.Println("[Client] Empty folder "+uploadPath+" has been created successfully", "\r")
	} else {
		fmt.Println("[Client] \u001B[36mWysyłanie powiodło się "+filename+"\u001B[0m", "\r")
	}
}

func UploadDir(uploadDir string, uploadPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	folder, err := os.ReadDir(uploadDir)
	if err != nil {
		fmt.Println("Folder can't be read", "\r")
		return
	}

	wgInside := sync.WaitGroup{}
	if len(folder) == 0 {
		wgInside.Add(1)
		go UploadFile("", uploadPath, &wgInside)
	} else {
		for _, file := range folder {
			if file.IsDir() {
				wgInside.Add(1)
				//	fmt.Println("[Dir] " + uploadDir + "/" + file.Name() + " ==> " + uploadPath)
				go UploadDir(uploadDir+"/"+file.Name(), uploadPath+"/"+file.Name(), &wgInside)
			} else {
				wgInside.Add(1)
				//	fmt.Println("[File] " + uploadDir + "/" + file.Name() + " ==> " + uploadPath)
				go UploadFile(uploadDir+"/"+file.Name(), uploadPath, &wgInside)
			}
		}
	}
	wgInside.Wait()
}

func ConnectServer() {
	// load settings
	config = settings.Load()

	// create json payload to authorize
	data := []byte(`{"pass":"` + config.UserPass + `",` + `"user":"` + config.UserName + `"` + `}`)
	req, err := http.NewRequest(http.MethodPost, "http://"+config.ConnectIP+":"+strconv.Itoa(config.ClientPort)+"/authorize", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		//panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// TODO: Dodaj pętle, możliwość exploracji oraz pobierania plików
	// TODO: Dodaj weryfikacje pobranych danych

	token = string(bodyBytes)
	if len(token) != 64 {
		//	fmt.Println(config.UserPass + " => " + config.UserName)
		//	fmt.Println(token)
		fmt.Println("[Client] Błąd autoryzacji")
		os.Exit(1)
	}
	fmt.Println("[Client] Autoryzacja ukończona pomyślne\r") //\n[>>" + token + "<<]")
	//	fmt.Println(token)

}

func Run() {

	// switch stdin into 'raw' modes
	var err error
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(fd int, oldState *term.State) {
		err := term.Restore(fd, oldState)
		if err != nil {

		}
	}(int(os.Stdin.Fd()), oldState)

	//b := make([]byte, 1)
	//_, err = os.Stdin.Read(b)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Printf("the char %q was hit", string(b[0]))

	fmt.Println("\033[31m<<< \u001B[0mSFE - Small File Exchanger \u001B[31m>>>\u001B[0m\r")
	fmt.Println("[\u001B[31m1\u001B[0m] Check connection to server\r")
	fmt.Println("[\u001B[31m2\u001B[0m] Host a server\r")
	fmt.Println("[\u001B[31m3\u001B[0m] Show config\r")
	fmt.Println("[\u001B[31m4\u001B[0m] Config DB\r")
	fmt.Println("[\u001B[31mX\u001B[0m] Exit\r")
	//	fmt.Print("Your choice: \r")
	input := readKey()

	switch string(input) {
	case "1": // Connect to server
		ConnectServer()
		ExploreDir("Pics")
		var wg sync.WaitGroup
		wg.Add(2)
		var bytesDownload uint64
		go DownloadFile("Pics", "cute.jpg", "", &wg, &bytesDownload)
		go UploadFile("client.png", "upstairs", &wg)
		wg.Wait()
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		if err != nil {
			return
		}

	case "2": // Start server
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		if err != nil {
			return
		}
		fmt.Println("\u001B[31mPress Ctrl+C to quit\u001B[0m\r")
		go listener.Host(-1)
		for {
			input := readKey()
			if input == 3 {
				fmt.Println("\u001B[31mCtrl+C detected, now exit...\u001B[0m\r")
				os.Exit(0)
			}
			fmt.Println("Key pressed " + strconv.Itoa(int(input)) + "\r")
		}

	case "3": // Print config
		config := settings.Load()

		fmt.Println("\r")
		fmt.Printf("File loaded: \u001B[31m%s\n\r\u001B[0m", viper.ConfigFileUsed())
		fmt.Println("\tServer Config:\r")
		fmt.Printf("Server Port: \u001B[31m%d\n\r\u001B[0m", config.ServerPort)
		fmt.Printf("Server DB: \u001B[31m%s\n\r\u001B[0m", config.ServerDB)
		fmt.Printf("SharedDir: \u001B[31m%s\n\n\r\u001B[0m", config.SharedDir)

		fmt.Println("\tClient Config:\r")
		fmt.Printf("Connect IP: \u001B[31m%s\n\r\u001B[0m", config.ConnectIP)
		fmt.Printf("Connect Port: \u001B[31m%d\n\r\u001B[0m", config.ClientPort)
		fmt.Printf("Username: \u001B[31m%s\n\r\u001B[0m", config.UserName)
		fmt.Printf("Userpass: \u001B[31m%s\n\r\u001B[0m", config.UserPass)
		fmt.Printf("DownloadDir: \u001B[31m%s\n\r\u001B[0m", config.DownloadDir)

		fmt.Print("<< Press enter to continue\n\r")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune()

		Run() //

	case "4": // Config Database - Add/Remove/View users
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		if err != nil {
			return
		}
		listener.ConfigDB(0)

	case "x":
		break
	case "X": // Close App
		fmt.Println("Exiting...\r")
		os.Exit(0)

	default: // Looping menu
		fmt.Println("Invalid choice\r")
		fmt.Print("<< Press enter to continue\n\r")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune()
		Run()
	}

}
