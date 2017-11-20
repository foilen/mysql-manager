package main

import (
	"fmt"
)

func showHelp() {
	fmt.Println("Arguments:")
	fmt.Println("  hostname:port : The host and port to connect to")
	fmt.Println("  admin username : Most likely 'root'")
	fmt.Println("  admin password file : The file with the password")
	fmt.Println("  databases.txt : The list of database names to keep or create")
	fmt.Println("  userPermissions.json : The list of users with their database access")
}
