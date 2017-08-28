package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func errHandler(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var shouldEncrypt bool
	var shouldDecrypt bool
	help := "\nUsage: -(e|d) password filename \nChoose e or d to encrypt or decrypt the input file."

	flag.BoolVar(&shouldEncrypt, "e", false, "Set for encryption")
	flag.BoolVar(&shouldDecrypt, "d", false, "Set for decryption")

	flag.Parse()

	if !shouldDecrypt && !shouldEncrypt {
		fmt.Println(help)
		return
	}

	pw := flag.Arg(0)
	file := flag.Arg(1)

	if pw == "" && file == "" {
		fmt.Println(help, "PW:", pw, "FILE:", file)
		return
	}

	// exe, e := os.Executable()
	// errHandler(e)
	// cwd := filepath.Dir(exe)

	aesCipher, e := aes.NewCipher([]byte(pw))
	if e != nil {
		fmt.Println("Error creating cipher!\nLength of password must be 16,24 or 32 characters! Current password length:", len(pw))
		errHandler(e)
	}

	data, e := ioutil.ReadFile(file)
	errHandler(e)

	//data length must be a multiple of the aes blocksize
	if dataRemain := len(data) % aes.BlockSize; dataRemain != 0 {
		padding := bytes.Repeat([]byte("\n"), aes.BlockSize-dataRemain)
		data = append(data, padding...)
	}

	if shouldEncrypt {
		encrypted := make([]byte, aes.BlockSize+len(data))
		iv := encrypted[:aes.BlockSize]

		_, e := io.ReadFull(rand.Reader, iv)
		errHandler(e)

		cfb := cipher.NewCFBEncrypter(aesCipher, iv)
		cfb.XORKeyStream(encrypted[aes.BlockSize:], data)

		output := base64.StdEncoding.EncodeToString(encrypted)

		outFile := "filencrypted_" + filepath.Base(file)
		outFile = filepath.Join(filepath.Dir(file), outFile)

		ioutil.WriteFile(outFile, []byte(output), 0740)
		fmt.Println("Encrypted file written to ", outFile)
	} else {
		decoded, e := base64.StdEncoding.DecodeString(string(data))
		errHandler(e)

		if len(decoded)%aes.BlockSize != 0 {
			fmt.Println("Encoded data is not correct length!")
			return
		}

		iv := decoded[:aes.BlockSize]
		output := decoded[aes.BlockSize:]

		cfb := cipher.NewCFBDecrypter(aesCipher, iv)
		cfb.XORKeyStream(output, output)

		outFile := strings.Replace(filepath.Base(file), "filencrypted_", "", -1)
		outFile = filepath.Join(filepath.Dir(file), outFile)

		ioutil.WriteFile(outFile, output, 0740)
		fmt.Println("Decrypted data has been written to ", outFile)

	}

}
