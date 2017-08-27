package main

import (
	"crypto/aes"
	"flag"
	"fmt"
	"io/ioutil"
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

	if !shouldDecrypt && !shouldDecrypt {
		fmt.Println(help)
		return
	}

	pw := flag.Arg(0)
	file := flag.Arg(1)

	cipher, e := aes.NewCipher([]byte(pw))
	if e != nil {
		fmt.Println("Error creating cipher!\nLength of password must be 16,24 or 32 characters! Current password length:", len(pw))
		errHandler(e)
	}

	data, e := ioutil.ReadFile(file)
	errHandler(e)

	if shouldDecrypt {

	} else {

	}

	fmt.Println("FILE:", file)
}
