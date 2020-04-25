package lsb

import (
	"go-lsb/utils"
	"math/rand"
	"strconv"
	"strings"
)

type BMPLSBRandomer struct {
	BmpLsb *BMPLSB
	Seed   string
}

func NewBMPLSBRandomer(BmpLsb *BMPLSB, seed string) (*BMPLSBRandomer, error) {
	lsb := &BMPLSBRandomer{
		BmpLsb: BmpLsb,
		Seed:   seed,
	}
	rand.Seed(int64(utils.Hash(seed)))
	return lsb, nil
}

func (b BMPLSBRandomer) Detect() bool {
	return b.BmpLsb.Detect()
}

func (b BMPLSBRandomer) InsertData(data []byte) error {
	header, err := b.BmpLsb.ComputeHeader(data)
	if err != nil {
		return err
	}
	index_used := make([]int64, 10)
	for _, bits := range append(header, data...) {
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			number := getNextInt(b.BmpLsb.Bmp.Size, index_used)
			err := b.BmpLsb.Bmp.SetSeekValue(number)
			if err != nil {
				return err
			}
			err = b.BmpLsb.writeBit(int32(bit))
			if err != nil {
				return err
			}
		}
	}
	b.BmpLsb.Bmp.Close()
	return nil
}

func (b BMPLSBRandomer) RetriveData() (msg []byte, err error) {
	index_used := make([]int64, 10)
	for i := 0; i < HEADER_SIZE; i++ {
		buf := make([]byte, 8)
		for j := 0; j < 8; j++ {
			number := getNextInt(b.BmpLsb.Bmp.Size, index_used)
			err := b.BmpLsb.Bmp.SetSeekValue(number)
			if err != nil {
				return nil, err
			}
			buf[7-j] = b.BmpLsb.ReadByte()
		}
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}
	lenght, err := strconv.Atoi(strings.Trim(string(msg), "\x00"))
	if err != nil {
		return nil, err
	}
	msg = []byte{}
	for i := 0; i < lenght; i++ {
		buf := make([]byte, 8)
		for j := 0; j < 8; j++ {
			number := getNextInt(b.BmpLsb.Bmp.Size, index_used)
			err := b.BmpLsb.Bmp.SetSeekValue(number)
			if err != nil {
				return nil, err
			}
			buf[7-j] = b.BmpLsb.ReadByte()
		}
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}

	b.BmpLsb.Bmp.Close()
	return msg, nil
}

func getNextInt(maxVal int64, index_used []int64) int64 {
	number := rand.Int63n(maxVal)
	for !utils.GetValidRand(number, index_used) {
		number = rand.Int63n(maxVal)
	}
	index_used = append(index_used, number)
	return number
}
