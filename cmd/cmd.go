package cmd

import (
	"fmt"
	"flag"
	"os"
	"strconv"
	"io/ioutil"


	"go-lsb/utils"
	"go-lsb/lsb"
)

func Execute(){
	inst := flag.Bool("insert", false, "Insertion de données\nUsage : -insert $src_file $dst_file $msg_file")
	rtve := flag.Bool("retrive", false, "Recupération de données\nUsage : -retrive $src_file $msg_lenght")
	key := flag.String("key", "", "Chiffrer l'entrée utilisateur\nUsage : -key $key")
	seed := flag.String("seed", "", "Insertion aléatoire des données\nUsage : -seed $seed")

	flag.Parse()

	if *rtve == *inst{
		flag.PrintDefaults()
		return 
	}

	if *inst{
		insert(*key, *seed, flag.Args())
	} else {
		fmt.Printf("%s\n", retrive(*key, *seed, flag.Args()))
	}
}

func insert(key, seed string, argv []string){
	if len(argv) != 3 {
		flag.PrintDefaults()
		return
	}

	if _, err := os.Stat(argv[0]); os.IsNotExist(err) {
		fmt.Println(err)
        return
	}
	
	err := utils.CopyFile(argv[0], argv[1])
	utils.CheckError(err)

	if _, err := os.Stat(argv[2]); os.IsNotExist(err) {
		fmt.Println(err)
        return
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

	LSB, err = lsb.NewBMP(argv[0])
	utils.CheckError(err)

	if randomise {
		LSB, err = lsb.NewBMPLSBRandomer(LSB.(*lsb.BMPLSB), seed)
		utils.CheckError(err)
	}
	
	err = LSB.InsertData(message)
	utils.CheckError(err)
}

func retrive(key, seed string, argv []string)[]byte{
	if len(argv) != 2 {
		flag.PrintDefaults()
		return nil
	}

	if _, err := os.Stat(argv[0]); os.IsNotExist(err) {
		fmt.Println(err)
        return nil
	}

	encrypt := key != ""
	randomise := seed != ""
	msg_lenght, err := strconv.Atoi(argv[1])
	utils.CheckError(err)

	var LSB lsb.Ilsb

	LSB, err = lsb.NewBMP(argv[0])
	utils.CheckError(err)

	if randomise {
		LSB, err = lsb.NewBMPLSBRandomer(LSB.(*lsb.BMPLSB), seed)
		utils.CheckError(err)
	}

	buf := make([]byte, msg_lenght)
	buf, err = LSB.RetriveData(msg_lenght)
	utils.CheckError(err)

	if encrypt {
		buf, err = utils.RC4Encryption([]byte(key), buf)
		utils.CheckError(err)
	}
	return buf
}
