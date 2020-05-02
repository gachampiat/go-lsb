// Package image permet de rassembler tous les types d'image
// que notre programme peut gérer.
package image

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// BMPParser est la structure qui va contenir
// les informations importantes pour notre image.
// C'est-à-dire l'adresse de début de l'image, la taille
// de l'image et l'adresse actuel de seek.
type BMPParser struct {
	File  *os.File
	Start int64
	Size  int64
	Seek  int64
}

// New retourne un pointeur sur une structure
// de type BMPParser.
func New(path string) (*BMPParser, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("Impossible d'ouvrir le fichier %s", path)
	}

	// D'après les spécification des images BMP
	// l'adresse de début est codée sur 4 octets.
	start := make([]byte, 4)
	_, err = f.ReadAt(start, 10)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture de l'adresse de début : %s", path)
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &BMPParser{
		File:  f,
		Start: int64(binary.LittleEndian.Uint32(start)),
		Size:  int64(stat.Size()) - int64(binary.LittleEndian.Uint32(start)),
	}, nil
}

// Close permet de fermer le fichier
// contenu dans la stucture
func (b *BMPParser) Close() {
	b.File.Close()
}

// SetSeekAtStartAddress permet de mettre la seek value
// sur l' adresse de départ des données images.
func (b *BMPParser) SetSeekAtStartAddress() error {
	if _, err := b.File.Seek(b.Start, 0); err != nil {
		return err
	}

	b.UpdateSeekValue()
	return nil
}

// UpdateSeekValue permet de mettre à jour la valeur
// Seek de la structure BMPParser.
func (b *BMPParser) UpdateSeekValue() {
	if n, err := b.File.Seek(0, 1); err != nil {
		log.Fatalf("Error update value seek : %s \n", err)
	} else {
		b.Seek = n
	}
}

// SetSeekValue permet de mettre à jour la valeur de la
// variable Seek de la structure BMPParser.
func (b *BMPParser) SetSeekValue(value int64) error {
	if _, err := b.File.Seek(value, 0); err != nil {
		return err
	}

	b.UpdateSeekValue()
	return nil
}
