package main

import (
	"os"
)

const a, c int = 1103515245, 12345

func main() {

	//Number of bytes to read at once
	const chunkSize = 128
	//Create a buffer to partially read the file into
	fileBuffer := make([]byte, chunkSize)
	//Create a buffer for file writing
	outputBuffer := make([]byte, chunkSize)

	//Array for command line arguments
	flags := os.Args[1:]

	//Open the file specified in the first command line arg
	//File to read from and XOR contents of
	plaintext, err1 := os.Open(flags[1])
	ciphertext, _ := os.Create(flags[2])

	//Close files when main exits
	defer ciphertext.Close()
	defer plaintext.Close()

	//Get initial X
	var x int = int(hash(flags[0]) % 256)
	x = keystreamGenerate(x)
	//Var to store how many bytes were read with each chunk
	//Will be 128 unless last chunk is read
	var numBytes int
	for err1 == nil {

		numBytes, err1 = plaintext.Read(fileBuffer)
		//fmt.Println(string(fileBuffer[:numBytes]))
		for i := 0; i < numBytes; i++ {

			//XOR byte with keystream and put into outputBuffer
			outputBuffer[i] = (byte(x) ^ fileBuffer[i])

			//Set x to the next keystream value
			x = keystreamGenerate(x)

		}

		//Write out the contents of the ouput buffer
		//Only write up to the number of bytes read in the last chunk to avoid writing too many bytes
		ciphertext.Write(outputBuffer[:numBytes])

	}

}

func hash(input string) uint64 {

	var hash uint64

	for _, c := range input {
		hash = uint64(c) + (hash << 6) + (hash << 16) - hash
	}

	return hash
}

func keystreamGenerate(x_n int) int {

	return (((x_n * a) + c) % 256)

}
