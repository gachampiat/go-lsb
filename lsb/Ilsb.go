// Package lsb permt de définir les structures / méthodes utiles pour
// réaliser des lsb.
package lsb

// Ilsb permet de définir les méthodes nécessaires
// pour qu'une structure puisse dans notre programme réaliser
// des lsb.
type Ilsb interface {
	InsertData(data []byte) error
	RetriveData() (msg []byte, err error)
	Detect() bool
}

// Pixel permet de représenter en mémoire un pixel.
// Nous avons fait un slice de int afin de stocker les
// valeurs de RGBA de chaque pixel.
// Nous mettons un type de type uint8 car la valeurs de chaque
// pixel ne peut pas dépasser 255.
type Pixel struct {
	RGBA []uint8
}

// HeaderSize défini une taille par défaut de l'header.
// Cette taille peut-être changer volontairement.
var HeaderSize = 8
