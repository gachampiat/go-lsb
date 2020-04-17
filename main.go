package main

import (
	"fmt"
	

	// "go-lsb/image"
	// "go-lsb/lsb"
	"go-lsb/utils"
)

func main(){
	message := "Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok Coucou axel je vais bien et toi ? j'espere que toi c'est ok"
	text := []byte(message)
	key := []byte("This is my key test")
	
	dst := make([]byte, len(message))
	dst, err := utils.RC4Encryption(key, text)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Printf("%s\n", dst)
	// path := "image.bmp"
	// utils.CopyFile(path, "lsb_image.bmp")
	// image, err := image.New("lsb_image.bmp")
	// defer image.Close()
	// utils.CheckError(err)

	// var LSB lsb.Ilsb
	// LSB, err = lsb.NewBMPLSBRandomer(image, "coucou")
	// utils.CheckError(err)
 	// err = LSB.InsertData(text)
	// utils.CheckError(err)
	// buf := make([]byte, len(message))
	// buf, err = LSB.RetriveData(len(message))
	// fmt.Println(string(buf))
}	