package main

import (
	"fmt"
	"log"

	"go-lsb/image"
)

func main(){
	test, e := image.New("image.bmp")
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(test)
}