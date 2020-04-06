package main

import (
	"go-lsb/image"
	"go-lsb/lsb"
	"go-lsb/utils"
)

func main(){
	b := []byte("AZERTYUIOP")
	image, err := image.New("image.bmp")
	defer image.Close()
	utils.CheckError(err)
	lsb := lsb.LSB{
		Bmp : image,
		Message : b,
	}
	err = lsb.InsertData()
	utils.CheckError(err)

}	