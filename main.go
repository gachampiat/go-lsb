package main

import (
	"go-lsb/image"
	"go-lsb/lsb"
	"go-lsb/utils"
)

func main(){
	b := []byte("AZERTYUIOP")
	path := "image.bmp"
	utils.CopyFile(path, "lsb_image.bmp")
	image, err := image.New("lsb_image.bmp")
	defer image.Close()
	utils.CheckError(err)
	lsb, err := lsb.NewLSBEncoder(image, b)
	utils.CheckError(err)
	err = lsb.InsertData()
	utils.CheckError(err)

}	