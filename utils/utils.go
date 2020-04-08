package utils

import (
	"fmt"
	"os"
	"log"
	"encoding/binary"
	"bytes"
	"strconv"
)

func CopyFile(src, dst string)error{
	buf := make([]byte, 100)

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("Impossible to create the file %s", dst)
	}

	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Impossible to open the file %s", src)
	}

	for {
		n, err := f.Read(buf)
		if err != nil {
				return err
		}
		if n == 0 {
				break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
				return err
		}
	}
	return nil 
}

func CheckError(err error){
	if err != nil{
		log.Fatal(err)
	}
}

func ByteToInt(slice []byte)(int64){
	buff := bytes.NewBuffer(slice)
	int, err := binary.ReadVarint(buff)
	if err != nil{
		log.Fatal(err) 
	}
	return int
}

func IntToBits(value int64) string{
	return fmt.Sprintf("%08s", strconv.FormatInt(value, 2))
}

func GetLSB(slice byte) uint8{
	return slice
}