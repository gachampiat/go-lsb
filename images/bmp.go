package images

import (
	"os"
	"fmt"
	"bufio"
)

type Header struct {
	MagicNumber [2]byte
	Size [4]byte
	Reserved [4]byte
	Start [4]byte
}

type Dib struct {
	
}

type Bmp struct {
	Header Header
}

func (h Header) String() string{
	return fmt.Sprintf("Magic Number = %-4x\t Size = %-4x\t Starting address = %-4x", h.MagicNumber, h.Size, h.Start)
}

func (b Bmp) GetHeader() Header{
	return b.Header
}

func New(fileSrc string)(error){
	f, err := os.Open(fileSrc)
	defer f.Close()
	if err != nil{
		return fmt.Errorf("Cannot open the file %v ", fileSrc)
	}

	reader := bufio.NewReader(f)

	data := make([]byte, 50)
	_, err = reader.Read(data)
	if err != nil{
		return err
	}
	var b Bmp
	b.CreateFromBytes(data)
	return nil

}

func (b *Bmp)CreateFromBytes(data []byte){
	fmt.Printf("%x\n", data)
	var h Header
	copy(h.MagicNumber[:], data[:2])
	copy(h.Size[:], data[2:6])
	copy(h.Reserved[:], data[6:10])
	copy(h.Start[:], data[10:12])
	fmt.Println(h)
}