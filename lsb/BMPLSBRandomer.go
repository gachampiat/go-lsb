package lsb

import(
	"fmt"
	"io"
	"math/rand"
	
	"go-lsb/image"
	"go-lsb/utils"
)

type BMPLSBRandomer struct {
	Bmp *image.BMP
	Seed string
}

func NewBMPLSBRandomer(filename string, seed string) (*BMPLSBRandomer, error){
	image, err := image.New(filename)
	if err != nil {
		return nil, err
	}
	lsb := &BMPLSBRandomer{
		Bmp : image,
		Seed : seed,
	}

	err = lsb.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return nil, err
	}

	fmt.Println(lsb.Bmp.Seek)
	return lsb, nil
}

func (l *BMPLSBRandomer) InsertData(data []byte)(error){
	if l.checkCapability(data){
		return fmt.Errorf("Please choose another stego-medium")
	}
	index_used := make([]int, 10)
	rand.Seed(int64(utils.Hash(l.Seed)))
	for _, bits := range data{
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			number := rand.Intn(int(l.Bmp.Size))
			for !utils.GetValidRand(number, index_used){
				fmt.Println("Index already used")
				number= rand.Intn(int(l.Bmp.Size))
			}
			err := l.Bmp.SetSeekValue(int64(number))
			if err != nil{
				return err
			}
			err = l.writeBit(int32(bit))
			if err != nil{
				return err
			}
		}
	}

	err := l.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return err
	}
	l.Bmp.Close()
	return nil
}

func (l *BMPLSBRandomer) checkCapability(data []byte) bool{
	return int64(len(data)) > (l.Bmp.Size / int64(8))
}

func (b *BMPLSBRandomer) writeBit(bit int32)(error){
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

func (l *BMPLSBRandomer) RetriveData(lenght int)(msg []byte, err error){
	index_used := make([]int, 10)
	rand.Seed(int64(utils.Hash(l.Seed)))
	for i := 0; i < lenght; i++{
		buf := make([]byte,  8)
		for j := 0 ; j < 8; j ++{
			number := rand.Intn(int(l.Bmp.Size))
			for !utils.GetValidRand(number, index_used){
				fmt.Println("Index already used")
				number= rand.Intn(int(l.Bmp.Size))
			}
			err := l.Bmp.SetSeekValue(int64(number))
			if err != nil{
				return nil, err
			}
			buf[7-j] = l.ReadByte()
		}
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}

	err = l.Bmp.SetSeekAtStartAddress()
	if err != nil{
		return nil, err
	}
	l.Bmp.Close()
	return msg, nil
}

func (l *BMPLSBRandomer) ReadByte()(byte){
	buf:= make([]byte, 1)
	l.Bmp.UpdateSeekValue()
	_, err := l.Bmp.File.Read(buf)
	if err != nil{
		fmt.Errorf("%s",err)
	}
	return buf[0]
}