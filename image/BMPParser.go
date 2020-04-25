package image

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type BMPParser struct {
	File  *os.File
	Start uint64
	Size  uint64
	Seek  int64
}

func New(path string) (*BMPParser, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("Impossible to open the file %s", path)
	}

	start := make([]byte, 4)
	_, err = f.ReadAt(start, 10)
	if err != nil {
		return nil, fmt.Errorf("Read at error : %s", path)
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &BMPParser{
		File: f,
		Start: uint64(binary.LittleEndian.Uint32(start)),
		Size: uint64(stat.Size()) - uint64(binary.LittleEndian.Uint32(start)),
	}, nil
}

func (b *BMPParser) Close() {
	b.File.Close()
}

func (b *BMPParser) SetSeekAtStartAddress() error {
	if _, err := b.File.Seek(int64(b.Start), 0); err != nil {
		return err
	} else {
		b.UpdateSeekValue()
		return nil
	}
}

func (b *BMPParser) UpdateSeekValue() {
	if n, err := b.File.Seek(0, 1); err != nil {
		log.Fatalf("Error update value seek : %s \n", err)
	} else {
		b.Seek = n
	}
}

func (b *BMPParser) SetSeekValue(value int64) error {
	if _, err := b.File.Seek(value, 0); err != nil {
		return err
	} else {
		b.UpdateSeekValue()
		return nil
	}
}
