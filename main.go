package main

import "fmt"

func main() {
	var cipherText string
	n, err := fmt.Scanln(&cipherText)
	fmt.Println("Input read:")
	fmt.Println(n, err)
	return
}
