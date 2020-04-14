package main

import (
	"fmt"
	"go-lsb/image"
	"go-lsb/lsb"
	"go-lsb/utils"
)

func main(){
	message := "Coucou axel tu es pas trés bon enfaite"
	text := []byte(message)

	path := "image.bmp"
	utils.CopyFile(path, "lsb_image.bmp")
	image, err := image.New("lsb_image.bmp")
	defer image.Close()
	utils.CheckError(err)

	var LSB lsb.Ilsb
	LSB, err = lsb.NewBMP(image)
	utils.CheckError(err)
 	err = LSB.InsertData(text)
	utils.CheckError(err)

	buf := make([]byte, len(message))
	buf, err = LSB.RetriveData(len(message))
	fmt.Println(string(buf))
}	