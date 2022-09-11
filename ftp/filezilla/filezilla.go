package filezilla

import (
	"fmt"
	"log"
	"os/user"
)

func display() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	username := currentUser.Username

	fmt.Printf("Username is: %s\n", username)
}
