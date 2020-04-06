package lsb

import(
	"fmt"

	"go-lsb/image"
)

type LSB struct {
	Bmp *image.BMP
	Message []byte
}

func (l *LSB) InsertData()(error){
	if l.checkCapability(){
		return fmt.Errorf("Use an other stega-medium for this message (stege-medium capability=%d, message lenght=%d)", l.Bmp.Size/8, len(l.Message))
	}
	l.Bmp.SetSeekAtStartAddress()

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