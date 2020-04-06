package image

import (
	"os"
	"fmt"
	"io"
	"encoding/binary"
)

type BMP struct {
	File *os.File
	Start int64
	Size int64
}

func (b BMP) String() string {
	return fmt.Sprintf("Path = %-20s\tStart = %-6d\t Buffer Size = %d Bytes", b.File.Name(), b.Start, b.Size)
}

func New(path string) (*BMP, error){	
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("Impossible to open the file %s", path)
	}

	// get the start address
	var start [4]byte
	_, err = f.ReadAt(start[:], 10)
	if err != nil && err != io.EOF{
		return nil, fmt.Errorf("Read at error : %s", path)
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &BMP{
		File: f,
		Start: int64(binary.BigEndian.Uint16(start[:])),
		Size : stat.Size() - int64(binary.BigEndian.Uint16(start[:])),
	}, nil
}

func (b *BMP) Close(){
	b.File.Close()
}

func (b *BMP) GetFile()(*os.File){
	return b.File
}

func (b *BMP) GetStartAddress() (int64){
	return b.Start
}

func (b *BMP) SetSeekAtStartAddress()error{
	if _, err := b.File.Seek(b.Start, 0); err != nil{
		return err
	} else {
		return nil
	}
}