package lsb

import(
	"fmt"
	"io"
	"strconv"
	"bytes"
  
	"go-lsb/image"
	"go-lsb/utils"
)

var HEADER_SIZE = 8

type BMPLSB struct {
	Bmp *image.BMP
}

func NewBMP(filename string) (*BMPLSB, error){
	image, err := image.New(filename)
	if err != nil {
		return nil, err
	}
	lsb := &BMPLSB{
		Bmp : image,
	}

	err = lsb.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return nil, err
	}

	return lsb, nil 
}

func (b BMPLSB) Detect() bool{
	return false
}

func (b BMPLSB) InsertData(data []byte)(error){
	buf, err := b.ComputeHeader(data)
	if err != nil {
		return err
	}
	for _, bits := range append(buf, data...){
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			err := b.writeBit(int32(bit))
			if err != nil{
				return err
			}
		}
	}

	err = b.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return err
	}
	b.Bmp.Close()
	return nil
}

func (b *BMPLSB) checkCapability(lenght int64) bool{
	return lenght > (b.Bmp.Size / int64(8))
}

func (b *BMPLSB) writeBit(bit int32)(error){
	buf:= make([]byte, 1)
	b.Bmp.UpdateSeekValue()
	_, err := b.Bmp.File.ReadAt(buf, b.Bmp.Seek)
	if err != nil && err != io.EOF{
		return err
	}
	
	var insert uint8 = byte(bit) | buf[0] >> 1 <<1
	buf_temp := make([]byte, 0)
	buf_temp = append(buf_temp, insert)
	
	_, err = b.Bmp.File.Write(buf_temp)
	if err != nil && err != io.EOF{
		return err
	}

	return nil
}

func (b BMPLSB) ComputeHeader(data []byte)([]byte, error){
	buf := []byte(strconv.Itoa(len(data)))
	padding := make([]byte, HEADER_SIZE - len(buf))
	if len(buf) > HEADER_SIZE {
		return nil, fmt.Errorf("The message lenght could not be bigger than 12 bytes (len=%d)\n", len(buf))
	} else if len(buf) < HEADER_SIZE {
		padding = append(padding, buf...)
	}

	if b.checkCapability(int64(len(data) + len(padding))){
		return nil, fmt.Errorf("Please choose another stego-medium")
	}
	
	return padding, nil
}

func (b BMPLSB) RetriveData()(msg []byte, err error){
	for i := 0; i < HEADER_SIZE; i++{
		buf := b.ReadNBytes(8)
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}
	lenght, err := strconv.Atoi(string(bytes.Trim(msg, "\x00")))
	if err != nil{
		return nil, err
	}
	msg = []byte{}
	for i := 0; i < lenght; i++{
		buf := b.ReadNBytes(8)
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}

	err = b.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return nil, err
	}
	b.Bmp.Close()
	return msg, nil
}

func (b *BMPLSB) ReadNBytes(n int)([]byte){
	buf := make([]byte, n + 1)
	for i:= n  ; i > 0; i--{
		buf[i] = b.ReadByte()
	}

	return buf

}

func (b *BMPLSB) ReadByte()(byte){
	buf:= make([]byte, 1)
	b.Bmp.UpdateSeekValue()
	_, err := b.Bmp.File.Read(buf)
	if err != nil{
		fmt.Errorf("%s",err)
	}
	return buf[0]
}