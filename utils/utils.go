package utils

import (
	"os"
	"log"
	"hash/fnv"
	"crypto/rc4"
	"io"
)

func Hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

func CopyFile(src, dst string)error{
	in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }
    return out.Close()
}

func CheckError(err error){
	if err != nil{
		log.Fatal(err)
		os.Exit(-1)
	}
}

func ByteSliceToInt(slice []byte)int{
	int_val := 0
	for i := 0; i <len(slice); i++{
		int_val = int_val << 1
		int_val = int(GetLsb(slice[i])) | int_val
	}
	return int_val
}

func GetLsb(value byte)byte{
	return  value << 7 >> 7
}

func GetValidRand(number int64, array []int64)bool{
	for _, value := range(array){
		if number == value{
			return false
		}
	}
	return true
}

func RC4Encryption(key, src []byte)([]byte, error){
	encrypt, err := rc4.NewCipher(key)
	if err != nil{
		return nil, err
	}
	dst := make([]byte, len(src))
	encrypt.XORKeyStream(dst, src)
	return dst, nil
}