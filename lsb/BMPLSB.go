package lsb

import(
	"fmt"
	"io"
  
	"go-lsb/image"
	"go-lsb/utils"
)

type BMPLSB struct {
	Bmp *image.BMP
}

func NewBMP(bmp *image.BMP) (*BMPLSB, error){
	lsb := &BMPLSB{
		Bmp : bmp,
	}

	err := lsb.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return nil, err
	}

	return lsb, nil
}

func (l *BMPLSB) InsertData(data []byte)(error){
	if l.checkCapability(data){
		return fmt.Errorf("Please choose another stego-medium")
	}
	for _, bits := range data{
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			l.writeBit(int32(bit))
		}
	}

	err := l.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return err
	}

	return nil
}

func (l *BMPLSB) checkCapability(data []byte) bool{
	return int64(len(data)) > (l.Bmp.Size / int64(8))
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

func (b *BMPLSB) RetriveData(lenght int)(msg []byte, err error){
	for i := 0; i < lenght; i++{
		buf := b.ReadNBytes(8)
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}

	err = b.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return nil, err
	}

	return msg, nil
}

func (l *BMPLSB) ReadNBytes(n int)([]byte){
	buf := make([]byte, n + 1)
	for i:= n  ; i > 0; i--{
		buf[i] = l.ReadByte()
	}

	return buf

}
func (l *BMPLSB) ReadByte()(byte){
	buf:= make([]byte, 1)
	l.Bmp.UpdateSeekValue()
	_, err := l.Bmp.File.Read(buf)
	if err != nil{
		fmt.Errorf("%s",err)
	}
	return buf[0]
}