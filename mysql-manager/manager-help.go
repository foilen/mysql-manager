package main

import (
	"fmt"
)

func showHelp() {
	fmt.Println("Arguments:")
	fmt.Println("  hostname:port : The host and port to connect to")
	fmt.Println("  manager-config.json : The admin credentials, the databases and the users permissions")
}
