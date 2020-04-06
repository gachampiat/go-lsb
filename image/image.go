package image

import (
	"os"
	"fmt"
	"io"
	"encoding/binary"
)

type BMP struct {
	File *os.File
	Start int64
}

func (b BMP) String() string {
	return fmt.Sprintf("Path = %-20s\tStart = %-6d", b.File.Name(), b.Start)
}

func New(path string) (*BMP, error){	
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("Impossible to open the file %s", path)
	}

	
	// get the start address
	var start [4]byte
	_, err = f.ReadAt(start[:], 10)
	if err != nil || err == io.EOF{
		return nil, fmt.Errorf("Read at error : %s", path)
	}

	return &BMP{
		File: f,
		Start: int64(binary.BigEndian.Uint16(start[:])),
	}, nil
}
