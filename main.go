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
	lsb, err := lsb.NewLSB(image, b)
	utils.CheckError(err)
	err = lsb.InsertData()
	utils.CheckError(err)

}	