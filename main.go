package main

import (
	"fmt"
	"go-lsb/image"
	"go-lsb/lsb"
	"go-lsb/utils"
)

func main(){
	message := "Coucou axel tu es pas tres bon enfaite"
	key := []byte("the-key-has-to-be-32-bytes-long!")
	text := []byte(message)

	ciphertext, err := utils.Encrypt(text, key)
    if err != nil {
        // TODO: Properly handle error
        fmt.Println(err)
	}
	path := "image.bmp"
	utils.CopyFile(path, "lsb_image.bmp")
	image, err := image.New("lsb_image.bmp")
	defer image.Close()
	utils.CheckError(err)

	encoder, err := lsb.NewLSBEncoder(image, ciphertext)
	utils.CheckError(err)
	
	err = encoder.InsertData()
	utils.CheckError(err)

	decoder := lsb.NewLSBDecoder(image)
	msg_decoded, err := utils.Decrypt(decoder.Decode(len(ciphertext)), key)
	utils.CheckError(err)
	fmt.Printf("%s\n", msg_decoded)
}	