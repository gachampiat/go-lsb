package utils

import (
	"fmt"
	"os"
	"log"
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

func IntToBits(value int64) string{
	return fmt.Sprintf("%08s", strconv.FormatInt(value, 2))
}

func ByteSliceToInt(slice []byte)int{
	int_val := 0
	for i := 0; i <len(slice); i++{
		int_val = int_val << 1
		int_val = int(GetLsb(slice[i])) | int_val
	}
	return int_val
}

func GetLsb (value byte)byte{
	return  value << 7 >> 7
}