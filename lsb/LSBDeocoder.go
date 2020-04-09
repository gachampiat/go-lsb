package lsb

import (
	"fmt"
	"go-lsb/image"
	"go-lsb/utils"
)

type LSBDecoder struct {
	Bmp *image.BMP
}

func NewLSBDecoder(bmp *image.BMP)(*LSBDecoder){
	lsb := &LSBDecoder{
		Bmp : bmp,
	}

	lsb.Bmp.SetSeekAtStartAddress()

	return lsb
}

func (l *LSBDecoder) Decode()(string, error){
	buf := l.ReadNBytes(8)
	fmt.Printf("Value : %c\n", utils.ByteSliceToInt(buf))
	buf = l.ReadNBytes(8)
	fmt.Printf("Value : %c\n", utils.ByteSliceToInt(buf))
	buf = l.ReadNBytes(8)
	fmt.Printf("Value : %c\n", utils.ByteSliceToInt(buf))
	return "", nil 
}

func (l *LSBDecoder) ReadNBytes(n int)([]byte){
	buf := make([]byte, n + 1)
	for i:= n  ; i > 0; i--{
		buf[i] = l.ReadByte()
	}

	return buf

}
func (l *LSBDecoder) ReadByte()(byte){
	buf:= make([]byte, 1)
	l.Bmp.UpdateSeekValue()
	_, err := l.Bmp.File.Read(buf)
	if err != nil{
		fmt.Errorf("%s",err)
	}
	return buf[0]
}