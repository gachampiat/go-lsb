package lsb

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/cmplx"
	"strconv"

	"go-lsb/image"
	"go-lsb/utils"
	"golang.org/x/image/bmp"
)

var HEADER_SIZE = 8

type BMPLSB struct {
	Bmp *image.BMPParser
}

func NewBMP(filename string) (*BMPLSB, error) {
	image, err := image.New(filename)
	if err != nil {
		return nil, err
	}
	lsb := &BMPLSB{
		Bmp: image,
	}

	err = lsb.Bmp.SetSeekAtStartAddress()
	if err != nil {
		return nil, err
	}

	return lsb, nil
}

func (b BMPLSB) Detect() bool {
	b.Bmp.SetSeekValue(0)
	img, err := bmp.Decode(b.Bmp.File)
	if err != nil {
		fmt.Println(err)
		return false
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	var x, y, k uint64
	for i := 0; i < height; i++ {
		for j := 0; j < width-1; j++ {
			for c := 0; c < 3; c++ {
				r := pixels[i][j].RGBA[c]
				s := pixels[i][j+1].RGBA[c]
				if (s%2 == 0 && r < s) || (s%2 == 1 && r > s) {
					x += 1
				}
				if (s%2 == 0 && r > s) || (s%2 == 1 && r < s) {
					y += 1
				}
				if math.Round(float64(s)/2) == math.Round(float64(r)/2) {
					k += 1
				}
			}
		}
	}

	if k == 0 {
		fmt.Println("SPA Failed")
		return false
	}

	a := float64(2 * k)
	bBis := float64(2 * (2*int64(x) - int64(width)*(int64(height)-1)))
	c := float64(y - x)

	bp := (-complex(bBis, 0) + cmplx.Sqrt(complex(math.Pow(bBis, 2)-4*a*c, 0))) / complex(2*a, 0)
	bm := (-complex(bBis, 0) - cmplx.Sqrt(complex(math.Pow(bBis, 2)-4*a*c, 0))) / complex(2*a, 0)

	beta := math.Min(math.Abs(real(bp)), math.Abs(real(bm)))
	fmt.Printf("Estimated embedding rate: %v \n", beta)
	return beta > 0.05
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{[]int{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}}
}

// Pixel struct example

func (b BMPLSB) InsertData(data []byte) error {
	buf, err := b.ComputeHeader(data)
	if err != nil {
		return err
	}
	for _, bits := range append(buf, data...) {
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			err := b.writeBit(int32(bit))
			if err != nil {
				return err
			}
		}
	}

	b.Bmp.Close()
	return nil
}

func (b *BMPLSB) checkCapability(lenght int64) bool {
	fmt.Printf("Text lenght = %d Maximal data on the image = %d \n", lenght, uint64(b.Bmp.Size)/uint64(8))
	return lenght > (b.Bmp.Size / int64(8))
}

func (b *BMPLSB) writeBit(bit int32) error {
	buf := make([]byte, 1)
	b.Bmp.UpdateSeekValue()
	_, err := b.Bmp.File.ReadAt(buf, b.Bmp.Seek)
	if err != nil && err != io.EOF {
		return err
	}

	var insert uint8 = byte(bit) | buf[0]>>1<<1
	buf_temp := make([]byte, 0)
	buf_temp = append(buf_temp, insert)

	_, err = b.Bmp.File.Write(buf_temp)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (b BMPLSB) ComputeHeader(data []byte) ([]byte, error) {
	buf := []byte(strconv.Itoa(len(data)))
	padding := make([]byte, HEADER_SIZE-len(buf))
	if len(buf) > HEADER_SIZE {
		return nil, fmt.Errorf("The message lenght could not be bigger than 12 bytes (len=%d)\n", len(buf))
	} else if len(buf) < HEADER_SIZE {
		padding = append(padding, buf...)
	}

	if b.checkCapability(int64(len(data) + len(padding))) {
		return nil, fmt.Errorf("Please choose another stego-medium")
	}

	return padding, nil
}

func (b BMPLSB) RetriveData() (msg []byte, err error) {
	for i := 0; i < HEADER_SIZE; i++ {
		buf := b.ReadNBytes(8)
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}
	lenght, err := strconv.Atoi(string(bytes.Trim(msg, "\x00")))
	if err != nil {
		return nil, err
	}
	msg = []byte{}
	for i := 0; i < lenght; i++ {
		buf := b.ReadNBytes(8)
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}

	err = b.Bmp.SetSeekAtStartAddress()
	if err != nil {
		return nil, err
	}
	b.Bmp.Close()
	return msg, nil
}

func (b *BMPLSB) ReadNBytes(n int) []byte {
	buf := make([]byte, n+1)
	for i := n; i > 0; i-- {
		buf[i] = b.ReadByte()
	}

	return buf

}

func (b *BMPLSB) ReadByte() byte {
	buf := make([]byte, 1)
	b.Bmp.UpdateSeekValue()
	_, err := b.Bmp.File.Read(buf)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	return buf[0]
}
