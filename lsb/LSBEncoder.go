package lsb

import(
	"fmt"
	"io"
	"strconv"
	"strings"

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
	
	lsb.Bmp.SetSeekAtStartAddress()

	return lsb, nil
}

func (l *LSBEncoder) InsertData()(error){
	l.writeHeader()
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
	fmt.Printf("BITS : %s\n", bits)
	for _, bit := range bits{
		l.writeBit(bit & 1)
	}
	return nil 
}

func (l *LSBEncoder) writeBit(bit int32)(error){
	seek := l.Bmp.GetSeekValue()

	buf := make([]byte, 1)
	_, err := l.Bmp.File.ReadAt(buf, seek)
	if err != nil && err != io.EOF{
		return err
	}

	bits := utils.IntToBits(utils.ByteToInt(buf))
	comp, err := strconv.ParseInt(bits[7:8], 10, 32)
	if err != nil {
		return err
	}

	if  comp != int64(bit){
		slice := strings.Split(bits, "")
		slice [7] = strconv.Itoa(int(bit))
		bits = strings.Join(slice, "")
	}
	i, err := strconv.ParseInt(bits, 2, 64)	
	if err != nil {
		return err
	}
	_, err = l.Bmp.File.Write([]byte(strconv.Itoa(int(i))))
	if err != nil && err != io.EOF{
		return err
	}

	return nil
}