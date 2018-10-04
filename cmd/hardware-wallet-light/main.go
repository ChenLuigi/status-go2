package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ebfe/scard"
	"github.com/status-im/status-go/smartcard/lightwallet"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s CAP_FILE_PATH\n", os.Args[0])
	}

	capFilePath := os.Args[1]

	ctx, err := scard.EstablishContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Release()

	readers, err := ctx.ListReaders()
	if err != nil {
		log.Fatal(err)
	}

	if len(readers) == 0 {
		log.Fatal("couldn't find any reader")
	}

	if len(readers) > 1 {
		log.Fatal("too many readers found")
	}

	reader := readers[0]
	fmt.Printf("using reader %s:\n", reader)
	fmt.Printf("connecting to card in %s\n", reader)
	card, err := ctx.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		log.Fatal(err)
	}
	defer card.Disconnect(scard.ResetCard)

	status, err := card.Status()
	if err != nil {
		log.Fatal(err)
	}

	switch status.ActiveProtocol {
	case scard.ProtocolT0:
		fmt.Printf("Protocol T0\n")
	case scard.ProtocolT1:
		fmt.Printf("Protocol T1\n")
	default:
		fmt.Printf("Unknown protocol\n")
	}

	f, err := os.Open(capFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	i := lightwallet.NewInstaller(card)
	secrets, err := i.Install(f, false)
	if err != nil {
		log.Fatal("installation error: ", err)
	}

	fmt.Printf("PUK %s\n", secrets.Puk())
	fmt.Printf("Pairing password: %s\n", secrets.PairingPass())
}
