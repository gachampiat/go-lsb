package lsb

import(
	"fmt"
	"io"
	// "strconv"
	// "strings"

	"go-lsb/image"
	"go-lsb/utils"
)

type LSBEncoder struct {
	Bmp *image.BMP
	Message []byte
}

func NewLSBEncoder(bmp *image.BMP, message []byte) (*LSBEncoder, error){
	lsb := &LSBEncoder{
		Bmp : bmp,
		Message : message,
	}

	if lsb.checkCapability(){
		return nil, fmt.Errorf("Use an other stega-medium for this message (stege-medium capability=%d, message lenght=%d)", lsb.Bmp.Size/8, len(lsb.Message))
	}
	
	err := lsb.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return nil, nil
	}

	return lsb, nil
}

func (l *LSBEncoder) InsertData()(error){
	// l.writeHeader()
	for _, bits := range l.Message{
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			l.writeBit(int32(bit))
		}
	}

	return nil
}

func (l *LSBEncoder) checkCapability() bool{
	return int64(len(l.Message)) > (l.Bmp.Size / int64(8))
}

func (l *LSBEncoder) writeHeader() error{
	bits := utils.IntToBits(int64(len(l.Message)))
	for _, bit := range bits{
		l.writeBit(bit & 1)
	}
	return nil 
}

func (l *LSBEncoder) writeBit(bit int32)(error){
	buf:= make([]byte, 1)
	l.Bmp.UpdateSeekValue()
	_, err := l.Bmp.File.ReadAt(buf, l.Bmp.Seek)
	if err != nil && err != io.EOF{
		return err
	}

	var insert uint8 = byte(bit) | buf[0] >> 1 <<1
	buf_temp := make([]byte, 0)
	buf_temp = append(buf_temp, insert)
	
	_, err = l.Bmp.File.Write(buf_temp)
	if err != nil && err != io.EOF{
		return err
	}

	return nil
}