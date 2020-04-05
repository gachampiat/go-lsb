package main

import (
	"fmt"
	
	"go-lsb/images"
)

func main(){
	bmp := images.New("image2.bmp")
	fmt.Println(bmp)
	// array := []byte("1010")
	// fmt.Printf("%b\n", array)
}