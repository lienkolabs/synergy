package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/lienkolabs/swell/crypto"
)

func createNewServer() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not retrieve USER home dir: %v\n", err)
	}

	var files []fs.DirEntry
	path := filepath.Join(homeDir, ".synergy")
	if files, err = os.ReadDir(path); err != nil {
		if err := os.Mkdir(path, fs.ModePerm); err != nil {
			log.Fatalf("could not create directort: %v\n", err)
		}
		if files, err = os.ReadDir(path); err != nil {
			log.Fatalf("unexpected error: %v\n", err)
		}
	}
	var instructionGateway, protocolGateway string
	if len(files) > 0 {
		return
	}
	token, _ := crypto.RandomAsymetricKey()
	fmt.Printf("You must grant power of attorney to the application key\n%v\n", token)

	fmt.Println("instruction gateway:")
	fmt.Scanln(&instructionGateway)
	fmt.Printf("Ok, connected to %v gateway\n", instructionGateway)
	fmt.Println("protocol gateway:")
	fmt.Scanln(&protocolGateway)
	fmt.Printf("Ok, connected to %v gateway\n", protocolGateway)

}

func main() {
	if len(os.Args) == 1 {
		// start server
		return
	}
	if os.Args[1] == "init" {
		createNewServer()
	}
}
