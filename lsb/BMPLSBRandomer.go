package lsb

import(
	"fmt"
	"math/rand"
	"go-lsb/utils"
	
)

type BMPLSBRandomer struct {
	BmpLsb *BMPLSB 
	Seed string
}

func NewBMPLSBRandomer(BmpLsb *BMPLSB, seed string) (*BMPLSBRandomer, error){
	lsb := &BMPLSBRandomer{
		BmpLsb : BmpLsb,
		Seed : seed,
	}

	return lsb, nil
}

func (b BMPLSBRandomer) InsertData(data []byte)(error){
	if b.BmpLsb.checkCapability(data){
		return fmt.Errorf("Please choose another stego-medium")
	}
	index_used := make([]int, 10)
	rand.Seed(int64(utils.Hash(b.Seed)))
	for _, bits := range data{
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			number := rand.Intn(int(b.BmpLsb.Bmp.Size))
			for !utils.GetValidRand(number, index_used){
				fmt.Println("Index already used")
				number= rand.Intn(int(b.BmpLsb.Bmp.Size))
			}
			err := b.BmpLsb.Bmp.SetSeekValue(int64(number))
			if err != nil{
				return err
			}
			err = b.BmpLsb.writeBit(int32(bit))
			if err != nil{
				return err
			}
		}
	}
	
	b.BmpLsb.Bmp.Close()
	return nil
}

func (b BMPLSBRandomer) RetriveData(lenght int)(msg []byte, err error){
	index_used := make([]int, 10)
	rand.Seed(int64(utils.Hash(b.Seed)))
	for i := 0; i < lenght; i++{
		buf := make([]byte,  8)
		for j := 0 ; j < 8; j ++{
			number := rand.Intn(int(b.BmpLsb.Bmp.Size))
			for !utils.GetValidRand(number, index_used){
				fmt.Println("Index already used")
				number= rand.Intn(int(b.BmpLsb.Bmp.Size))
			}
			err := b.BmpLsb.Bmp.SetSeekValue(int64(number))
			if err != nil{
				return nil, err
			}
			buf[7-j] = b.BmpLsb.ReadByte()
		}
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}

	b.BmpLsb.Bmp.Close()
	return msg, nil
}