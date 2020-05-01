package lsb

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"math/cmplx"
	"strconv"

	"github.com/gachampiat/go-lsb/image"
	"github.com/gachampiat/go-lsb/utils"

	"golang.org/x/image/bmp"
)

// BMPLSB Structure de base permettant de socker un BMPParser
type BMPLSB struct {
	Bmp *image.BMPParser
}

// NewBMP retourne un pointeur sur une structure
// de type BMPLSB. Le paramètre filename, sera passé
// à la méthode de construction de la structure BMPParser.
// Note : Cette méthode place la seek value au pixel de début
// de l'image du bmp.
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

// Detect la fonction de détéction de notre programme
// Nous avons implémenter la méthode de détéction SPA.
func (b BMPLSB) Detect() bool {
	// on place la seek value à l'index 0
	b.Bmp.SetSeekValue(0)
	// on decode le fichier en structure image.Image
	img, err := bmp.Decode(b.Bmp.File)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// va permettre de stocker les pixels de l'image
	var pixels [][]Pixel
	// on parcourt l'image pour remplir notre variable pixels
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
					x++
				}
				if (s%2 == 0 && r > s) || (s%2 == 1 && r < s) {
					y++
				}
				if math.Round(float64(s)/2) == math.Round(float64(r)/2) {
					k++
				}
			}
		}
	}

	if k == 0 {
		log.Fatal("SPA FAILED")
	}

	a := float64(2 * k)
	bBis := float64(2 * (2*int64(x) - int64(width)*(int64(height)-1)))
	c := float64(y - x)

	bp := (-complex(bBis, 0) + cmplx.Sqrt(complex(math.Pow(bBis, 2)-4*a*c, 0))) / complex(2*a, 0)
	bm := (-complex(bBis, 0) - cmplx.Sqrt(complex(math.Pow(bBis, 2)-4*a*c, 0))) / complex(2*a, 0)

	beta := math.Min(math.Abs(real(bp)), math.Abs(real(bm)))
	log.Println(beta)
	return beta > 0.05
}

// InsertData insert les données dans l'image.
// Cette fonction insére les données les unes
// à la suite des autres.
func (b BMPLSB) InsertData(data []byte) error {
	defer b.Bmp.Close()

	// on calcul le header correspondant à la donnée
	header, err := b.ComputeHeader(data)
	if err != nil {
		return err
	}

	// on parcourt tous les bytes a insérer.
	for _, bits := range append(header, data...) {
		var i uint8
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			err := b.writeBit(int32(bit))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ComputeHeader est une fonction qui calcule la taille
// du message est l'insert dans un slice de byte. La taille de
// ce slice est défini par la variable global HeaderSize.
// Le slice retourné est rempli par la gauche de 0 afin d'atteindre
// la taille défini par la variable.
func (b BMPLSB) ComputeHeader(data []byte) ([]byte, error) {
	buf := []byte(strconv.Itoa(len(data)))
	padding := make([]byte, HeaderSize-len(buf))

	if len(buf) > HeaderSize {
		return nil, fmt.Errorf("The message lenght could not be bigger than 12 bytes (len=%d)", len(buf))
	} else if len(buf) < HeaderSize {
		padding = append(padding, buf...)
	}

	if b.checkCapability(uint64(len(data) + len(padding))) {
		return nil, fmt.Errorf("Please choose another stego-medium")
	}

	return padding, nil
}

// RetriveData récupére le message caché dans l'image.
// Pour la récupération, nous allons lire les données
// les unes à la suite des autres.
func (b BMPLSB) RetriveData() (msg []byte, err error) {
	defer b.Bmp.Close()

	for i := 0; i < HeaderSize; i++ {
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

	return msg, nil
}

// ReadNBytes lit n bytes dans l'image et
// retourne un slice avec les bytes lu
func (b *BMPLSB) ReadNBytes(n int) []byte {
	buf := make([]byte, n+1)
	for i := n; i > 0; i-- {
		buf[i], _ = b.ReadByte()
	}
	return buf
}

// ReadByte lit et retourne un byte lu dans l'image.
// Avant de lire, le programme mets à jour la valeur
// de seek.
func (b *BMPLSB) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	b.Bmp.UpdateSeekValue()
	_, err := b.Bmp.File.Read(buf)
	if err != nil {
		log.Fatalf("Erreur lors de la fonction de lecture : %s\n Seek value = %d", err, b.Bmp.Seek)
	}
	return buf[0], nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{[]uint8{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)}}
}

func (b *BMPLSB) checkCapability(lenght uint64) bool {
	log.Printf("Taille du message = %d \n", lenght)
	log.Printf("Nombre de pixel dans l'image = %d \n", b.Bmp.Size)
	fmt.Printf("Taux Stéganographique = %f bits par pixel \n", float64(lenght)/(float64(b.Bmp.Size)/3))
	return lenght > (uint64(b.Bmp.Size) / uint64(8))
}

func (b *BMPLSB) writeBit(bit int32) error {
	buf := make([]byte, 1)
	b.Bmp.UpdateSeekValue()
	_, err := b.Bmp.File.ReadAt(buf, b.Bmp.Seek)
	if err != nil && err != io.EOF {
		return err
	}

	var insert uint8 = byte(bit) | buf[0]>>1<<1
	bufTemp := make([]byte, 0)
	bufTemp = append(bufTemp, insert)

	_, err = b.Bmp.File.Write(bufTemp)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
