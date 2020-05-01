package lsb

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/gachampiat/go-lsb/utils"
)

// BMPLSBRandomer est une structure basée sur le patron de conception
// décorateur. En effet son objectif est d'ajouté dynamiquement de
// comportement à la stuctrure de basse BmpLsb.
// La seed nous permettra de réaliser une insertion "aléatoire". Elle
// sera transformé en int via une fonction de hash.
type BMPLSBRandomer struct {
	BmpLsb *BMPLSB
	Seed   string
}

// NewBMPLSBRandomer créer et retourne un pointeur sur la structure BMPLSBRandomer.
func NewBMPLSBRandomer(BmpLsb *BMPLSB, seed string) (*BMPLSBRandomer, error) {
	lsb := &BMPLSBRandomer{
		BmpLsb: BmpLsb,
		Seed:   seed,
	}
	rand.Seed(int64(utils.Hash(seed)))
	return lsb, nil
}

// Detect utilise la méthode de détéction
// de la structure de base. En effet ici nous
// n'avons pas besoin de rajouter du comportement
// à la fonction de détéction.
func (b BMPLSBRandomer) Detect() bool {
	return b.BmpLsb.Detect()
}

// InsertData est une redéfinition de la méthode.
// Nous ajoutons un header au message de base afin
// d'insérer la taille du message. Puis, nous utilisons
// la seed afin de faire une selection "aléatoire" de pixel.
// Cette selection se base entre la borne du début de l'image
// et le fin de l'image. Ces bornes sont trouvés par
// la classe BMPParser.
func (b BMPLSBRandomer) InsertData(data []byte) error {
	// Dans tous les cas, ont ferme le fichier.
	defer b.BmpLsb.Bmp.Close()

	// On calcule automatique le header.
	header, err := b.BmpLsb.ComputeHeader(data)
	if err != nil {
		return err
	}
	// Pour stocker les index utilisés.
	// Evite d'insérer des données au même endroit
	indexUsed := make([]int64, 10)
	for _, bits := range append(header, data...) {
		// Pour stocker les valeurs de chaque bit
		var i uint8
		// une boucle de 8 car un byte = 8 bits
		for i = 0; i < 8; i++ {
			bit := (bits & byte(1<<i)) >> i
			number := getNextInt(int64(b.BmpLsb.Bmp.Start), int64(b.BmpLsb.Bmp.Size), indexUsed)

			// on update la valeur du seek dans le fichier
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

	return nil
}

// RetriveData va récupérer les données dans une image
func (b BMPLSBRandomer) RetriveData() (msg []byte, err error) {
	// Dans tous les cas, ont ferme le fichier.
	defer b.BmpLsb.Bmp.Close()

	// stocke les index déjà lu
	indexUsed := make([]int64, 10)

	// Permet de décoder le header
	for i := 0; i < HeaderSize; i++ {
		buf := make([]byte, 8)
		for j := 0; j < 8; j++ {
			number := getNextInt(b.BmpLsb.Bmp.Start, b.BmpLsb.Bmp.Size, indexUsed)
			err := b.BmpLsb.Bmp.SetSeekValue(number)
			if err != nil {
				return nil, err
			}
			buf[7-j], _ = b.BmpLsb.ReadByte()
		}
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}
	// TODO : Trouver une meilleure solution.
	// Permet de supprimer les octets nuls lors de la conversion
	// bytes to string to int
	lenght, err := strconv.Atoi(strings.Trim(string(msg), "\x00"))
	if err != nil {
		return nil, err
	}

	msg = []byte{}
	for i := 0; i < lenght; i++ {
		buf := make([]byte, 8)
		for j := 0; j < 8; j++ {
			number := getNextInt(b.BmpLsb.Bmp.Start, b.BmpLsb.Bmp.Size, indexUsed)
			err := b.BmpLsb.Bmp.SetSeekValue(number)
			if err != nil {
				return nil, err
			}
			// ici nous inverson l'ordre des bits car pendant l'insertion
			// nous insérons les bits dans le sens contraire.
			buf[7-j], _ = b.BmpLsb.ReadByte()
		}
		msg = append(msg, byte(utils.ByteSliceToInt(buf)))
	}

	return msg, nil
}

// getNextInt retourne une valeur "aléatoire" compris entre
// la valeur minimale + 1 et la valeur maximale et qui n'est pas
// présente dans le tableau indexUsed. Une fois la valeur trouvé,
// elle est automatiquement ajouté au tableau indexUsed.
// Note: Par défaut en golang ce tableau est un pointeur, alors il
// sera automatiquement mis à jour en dehors de la fonction.
func getNextInt(minVal, maxVal int64, indexUsed []int64) int64 {
	number := rand.Int63n(maxVal-minVal+1) + minVal + 1
	for !utils.InSlice(number, indexUsed) {
		number = rand.Int63n(maxVal)
	}

	indexUsed = append(indexUsed, number)
	return number
}
