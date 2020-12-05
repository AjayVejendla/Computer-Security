package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var commandArgs []string = os.Args[1:]

func main() {
	keyFile, err1 := ioutil.ReadFile(commandArgs[0])
	plaintext, err2 := os.Open(commandArgs[1])
	ciphtertext, err3 := os.Create(commandArgs[2])

	keyLength := int64(len(keyFile))

	if err1 != nil || err2 != nil || err3 != nil {

		fmt.Println("An error has occurred")
		return

	}
	fileBuffer := make([]byte, 128)
	outputBuffer := make([]byte, 128)
	var position int64 = 0
	var numBytes int
	for err1 == nil {
		numBytes, err1 = plaintext.Read(fileBuffer)

		for i := 0; i < numBytes; i++ {

			outputBuffer[i] = (fileBuffer[i] + keyFile[(position%keyLength)]) % 255
			position = position + 1

		}

		ciphtertext.Write(outputBuffer[:numBytes])
	}

}
