package lsb

import(
	"fmt"
	"io"
	"strconv"
	"strings"

	"go-lsb/image"
	"go-lsb/utils"
)

type LSB struct {
	Bmp *image.BMP
	Message []byte
}

func NewLSB(bmp *image.BMP, message []byte) (*LSB, error){
	lsb := &LSB{
		Bmp : bmp,
		Message : message,
	}

	if lsb.checkCapability(){
		return nil, fmt.Errorf("Use an other stega-medium for this message (stege-medium capability=%d, message lenght=%d)", lsb.Bmp.Size/8, len(lsb.Message))
	}
	
	lsb.Bmp.SetSeekAtStartAddress()

	return lsb, nil
}

func (l *LSB) InsertData()(error){
	l.writeHeader()
	for _, bits := range l.Message{
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			fmt.Print(bit)
		}
	}

	return nil
}

func (l *LSB) checkCapability() bool{
	return int64(len(l.Message)) > (l.Bmp.Size / int64(8))
}

func (l *LSB) writeHeader() error{
	bits := utils.IntToBits(int64(len(l.Message)))
	fmt.Printf("BITS : %s\n", bits)
	for _, bit := range bits{
		l.writeBit(bit & 1)
	}
	return nil 
}

func (l *LSB) writeBit(bit int32)(error){
	seek, err := l.Bmp.GetSeekValue()
	if err != nil {
		return err
	}

	buf := make([]byte, 1)
	
	fmt.Printf("SEEK VALUE : %d \n", seek)
	_, err = l.Bmp.File.ReadAt(buf, seek)
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

	fmt.Printf("SEEK VALUE : %d \n", seek)
	_, err = l.Bmp.File.WriteAt([]byte(bits), seek)
	if err != nil && err != io.EOF{
		return err
	}
	l.Bmp.SetSeekValue(seek + 1)
	fmt.Printf("SEEK VALUE : %d \n", seek)

	return nil
}