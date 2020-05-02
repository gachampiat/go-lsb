// Package utils permet la définition de toutes les méthodes
// utiles pour le programme
package utils

import (
	"crypto/rc4"
	"hash/fnv"
	"io"
	"log"
	"os"
)

// Hash prend en entrée la chaine de caractère s afin de la
// transformée en uint64 via une fonction de hash
func Hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

// CopyFile permet de copier le contenu d'un fichier vers le fichier
// destination. Cette fonction remonte des erreurs si le fichier source
// n'existe pas, si le programme n'arrive pas a créer le fichier destination
// ou s'il y a une erreur pendant la copie.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// CheckError prend une erreur en entrée, si celle-ci
// n'est pas nulle, alors l'erreur est affiché et le programme
// se finit.
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ByteSliceToInt prend un slice de byte afin de le transformer
// en int. Cette fonction récupére tous les lsb de chaque byte
// présent dans le slice. Ensuite via des décalages binaires
// elle créer un int.
func ByteSliceToInt(slice []byte) int {
	intVal := 0
	for i := 0; i < len(slice); i++ {
		intVal = intVal << 1
		intVal = int(GetLsb(slice[i])) | intVal
	}
	return intVal
}

// GetLsb retroune la valeur du dérnier bit d'un Byte.
// Pour cela on fait un premier décalage à gauche de 7 bit, car
// un octet est encodé sous 8 bits.
func GetLsb(value byte) byte {
	return value << 7 >> 7
}

// InSlice retourne false si la valeur n'est pas dans
// le slice, sinon elle retourne true
func InSlice(number int64, array []int64) bool {
	for _, value := range array {
		if number == value {
			return true
		}
	}
	return false
}

// RC4Encryption permet de chiffrer le paramètre src
// avec le clef de chiffrement key.
func RC4Encryption(key, src []byte) ([]byte, error) {
	encrypt, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(src))
	encrypt.XORKeyStream(dst, src)
	return dst, nil
}
