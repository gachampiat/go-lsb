package main

import (
	// "fmt"
	"go-lsb/image"
	"go-lsb/lsb"
	"go-lsb/utils"
)

func main(){
	b := []byte("AZv")
	path := "image.bmp"
	utils.CopyFile(path, "lsb_image.bmp")
	image, err := image.New("lsb_image.bmp")
	defer image.Close()
	utils.CheckError(err)

	encoder, err := lsb.NewLSBEncoder(image, b)
	utils.CheckError(err)
	
	err = encoder.InsertData()
	utils.CheckError(err)

	decoder := lsb.NewLSBDecoder(image)
	decoder.Decode()
}	