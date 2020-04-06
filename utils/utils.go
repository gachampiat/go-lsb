package utils

import (
	"fmt"
	"os"
)

func CopyFile(src, dst string)error{
	buf := make([]byte, 100)

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("Impossible to create the file %s", dst)
	}

	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Impossible to open the file %s", src)
	}

	for {
		n, err := f.Read(buf)
		if err != nil {
				return err
		}
		if n == 0 {
				break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
				return err
		}
	}
	return nil 
}