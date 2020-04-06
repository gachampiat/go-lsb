package lsb

import(
	"fmt"

	"go-lsb/image"
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
	for _, bits := range l.Message{
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			fmt.Print(bit)
		}
		fmt.Printf("---- %d ---- %b\n", bits, bits)
	}

	return nil
}

func (l *LSB) checkCapability() bool{
	return int64(len(l.Message)) > (l.Bmp.Size / int64(8))
}