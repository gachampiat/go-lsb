// Package cmd permet de définir toutes les intéractions
// possible depuis la ligne de commande.
package cmd

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/gachampiat/go-lsb/utils"

	"github.com/gachampiat/go-lsb/lsb"
)

// Execute est la fonction principale de ce package
// elle permet d'appeller les fonctionnalités en fonction
// des flags passés à la cmd.
func Execute() {
	inst := flag.Bool("insert", false, "Insertion de données\nUsage : -insert $src_file $dst_file $msg_file")
	rtve := flag.Bool("retrive", false, "Recupération de données\nUsage : -retrive $src_file")
	key := flag.String("key", "", "Chiffrer l'entrée utilisateur\nUsage : -key $key")
	seed := flag.String("seed", "", "Insertion aléatoire des données\nUsage : -seed $seed")
	dct := flag.Bool("detect", false, "Détection d'un message dans une image\nUsage : -detect $img_src")

	flag.Parse()

	if *inst {
		insert(*key, *seed, flag.Args())
	} else if *dct {
		detect(flag.Args())
	} else if *rtve {
		log.Printf("%s\n", retrive(*key, *seed, flag.Args()))
	} else {
		flag.PrintDefaults()
	}
}

func detect(argv []string) {
	if len(argv) != 1 {
		flag.PrintDefaults()
		return
	}

	if _, err := os.Stat(argv[0]); os.IsNotExist(err) {
		log.Fatal(err)
	}

	var LSB lsb.Ilsb

	LSB, err := lsb.NewBMP(argv[0])
	utils.CheckError(err)

	LSB.Detect()

}

func insert(key, seed string, argv []string) {
	if len(argv) != 3 {
		flag.PrintDefaults()
		return
	}

	if _, err := os.Stat(argv[0]); os.IsNotExist(err) {
		log.Fatal(err)
	}

	err := utils.CopyFile(argv[0], argv[1])
	utils.CheckError(err)

	if _, err := os.Stat(argv[2]); os.IsNotExist(err) {
		log.Fatal(err)
	}

	message := make([]byte, 20)
	message, err = ioutil.ReadFile(argv[2])
	utils.CheckError(err)

	encrypt := key != ""
	randomise := seed != ""

	if encrypt {
		message, err = utils.RC4Encryption([]byte(key), message)
		utils.CheckError(err)
	}

	var LSB lsb.Ilsb

	LSB, err = lsb.NewBMP(argv[1])
	utils.CheckError(err)

	if randomise {
		LSB, err = lsb.NewBMPLSBRandomer(LSB.(*lsb.BMPLSB), seed)
		utils.CheckError(err)
	}
	err = LSB.InsertData(message)
	utils.CheckError(err)
	log.Println("Message insérée")
}

func retrive(key, seed string, argv []string) []byte {
	if len(argv) != 1 {
		flag.PrintDefaults()
		return nil
	}

	if _, err := os.Stat(argv[0]); os.IsNotExist(err) {
		log.Fatal(err)
	}

	encrypt := key != ""
	randomise := seed != ""

	var LSB lsb.Ilsb

	LSB, err := lsb.NewBMP(argv[0])
	utils.CheckError(err)

	if randomise {
		LSB, err = lsb.NewBMPLSBRandomer(LSB.(*lsb.BMPLSB), seed)
		utils.CheckError(err)
	}

	buf := make([]byte, 10)
	buf, err = LSB.RetriveData()
	utils.CheckError(err)

	if encrypt {
		buf, err = utils.RC4Encryption([]byte(key), buf)
		utils.CheckError(err)
	}
	return buf
}
