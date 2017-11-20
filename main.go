package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Arguments
	if len(os.Args) != 3 {
		showHelp()
		os.Exit(1)
	}
	hostPort := os.Args[1]
	managerConfigJsonFile := os.Args[2]

	// Get the config
	managerConfiguration, err := getManagerConfiguration(managerConfigJsonFile)
	if err != nil {
		log.Fatal(err)
	}

	// Get connection string
	mysqlConnectionString := managerConfiguration.Admin.Name + ":" + managerConfiguration.Admin.Password + "@" + "tcp(" + hostPort + ")/mysql"

	// Connect
	db, err := sql.Open("mysql", mysqlConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// flush privileges;
	fmt.Println("-----")
	fmt.Println("Flush privileges")
	sqlExecOrFail(db, "FLUSH PRIVILEGES")
	fmt.Println("-----")

	// List existing databases
	existingDatabases := sqlStringArrayOrFail(db, "SHOW DATABASES")
	existingDatabases = arrayRemoveAll(existingDatabases, "information_schema", "mysql", "performance_schema", "sys")
	fmt.Println("Existing databases:", existingDatabases)

	// List desired databases
	fmt.Println("Desired databases:", managerConfiguration.Databases)
	fmt.Println("-----")

	// Drop extra databases
	for _, database := range existingDatabases {
		if !arrayContains(managerConfiguration.Databases, database) {
			fmt.Println("Drop database", database)
			sqlExecOrFail(db, "DROP DATABASE "+database)
		}
	}

	// Create missing databases
	for _, database := range managerConfiguration.Databases {
		if !arrayContains(existingDatabases, database) {
			fmt.Println("Create database", database)
			sqlExecOrFail(db, "CREATE DATABASE "+database)
		}
	}

	// List existing users
	rows, err := db.Query("SELECT user, host FROM user")
	defer rows.Close()
	var existingUsers []string
	for rows.Next() {
		var userName string
		var userHost string
		err = rows.Scan(&userName, &userHost)
		if err != nil {
			log.Fatal(err)
		}
		existingUsers = append(existingUsers, "'"+userName+"'@'"+userHost+"'")
	}
	fmt.Println("Existing users:", existingUsers)

	// List desired users
	var desiredUsers []string
	for _, user := range managerConfiguration.UsersToIgnore {
		desiredUsers = append(desiredUsers, "'"+user.Name+"'@'"+user.Host+"'")
	}
	for _, user := range managerConfiguration.UsersPermissions {
		desiredUsers = append(desiredUsers, "'"+user.Name+"'@'"+user.Host+"'")
	}

	// Drop extra users
	fmt.Println("Desired users:", desiredUsers)
	fmt.Println("-----")
	for _, user := range existingUsers {
		if !arrayContains(desiredUsers, user) {
			fmt.Println("Drop user", user)
			sqlExecOrFail(db, "DROP USER "+user)
		}
	}

	// Create missing users
	for _, user := range managerConfiguration.UsersPermissions {
		fullUserName := "'" + user.Name + "'@'" + user.Host + "'"
		if !arrayContains(existingUsers, fullUserName) {
			fmt.Println("Create user", fullUserName)
			sqlExecOrFail(db, "CREATE USER "+fullUserName)
		}
	}

	// Update passwords
	for _, userPermissions := range managerConfiguration.UsersPermissions {
		fullUserName := "'" + userPermissions.Name + "'@'" + userPermissions.Host + "'"
		if userPermissions.Password != "" {
			fmt.Println("Update user", fullUserName, "password")
			sqlExecOrFail(db, "ALTER USER "+fullUserName+" IDENTIFIED BY '"+userPermissions.Password+"'")
		}
	}
	fmt.Println("-----")

	// Get users grants
	for _, userPermissions := range managerConfiguration.UsersPermissions {
		fullUserName := "'" + userPermissions.Name + "'@'" + userPermissions.Host + "'"
		var desiredUserDatabases []string
		for _, databaseGrant := range userPermissions.DatabaseGrants {
			desiredUserDatabases = append(desiredUserDatabases, databaseGrant.DatabaseName)
		}

		// Remove all permissions for databases that the user does not have access
		existingDatabaseGrants := sqlStringArrayOrFail(db, "SELECT DISTINCT db FROM db WHERE user = ? AND host = ? ORDER BY db", userPermissions.Name, userPermissions.Host)
		for _, existingDatabaseGrant := range existingDatabaseGrants {
			if !arrayContains(desiredUserDatabases, existingDatabaseGrant) {
				fmt.Println("Drop all grants for user", fullUserName, "on database", existingDatabaseGrant)
				sqlExecOrFail(db, "DELETE FROM db WHERE user = ? AND host = ? AND db = ?", userPermissions.Name, userPermissions.Host, existingDatabaseGrant)
			}
		}

		// Check all existing DB grants and update them
		for _, databaseGrant := range userPermissions.DatabaseGrants {
			row := db.QueryRow("SELECT select_priv, insert_priv, update_priv, delete_priv, create_priv, drop_priv, alter_priv FROM db WHERE user = ? AND host = ? AND db = ?", userPermissions.Name, userPermissions.Host, databaseGrant.DatabaseName)
			var isSelect = "N"
			var isInsert = "N"
			var isUpdate = "N"
			var isDelete = "N"
			var isCreate = "N"
			var isDrop = "N"
			var isAlter = "N"
			err := row.Scan(&isSelect, &isInsert, &isUpdate, &isDelete, &isCreate, &isDrop, &isAlter)
			switch {
			case err == sql.ErrNoRows:
			case err != nil:
				log.Fatal(err)
			}
			grantOrRevoke(db, fullUserName, databaseGrant, "SELECT", isSelect)
			grantOrRevoke(db, fullUserName, databaseGrant, "INSERT", isInsert)
			grantOrRevoke(db, fullUserName, databaseGrant, "UPDATE", isUpdate)
			grantOrRevoke(db, fullUserName, databaseGrant, "DELETE", isDelete)
			grantOrRevoke(db, fullUserName, databaseGrant, "CREATE", isCreate)
			grantOrRevoke(db, fullUserName, databaseGrant, "DROP", isDrop)
			grantOrRevoke(db, fullUserName, databaseGrant, "ALTER", isAlter)

		}
	}

	// flush privileges;
	fmt.Println("-----")
	fmt.Println("Flush privileges")
	sqlExecOrFail(db, "FLUSH PRIVILEGES")
}

func needGrant(desiredGrants []string, grant string, currentGrant string) bool {
	return arrayContains(desiredGrants, grant) && currentGrant == "N"
}

func needRevoke(desiredGrants []string, grant string, currentGrant string) bool {
	return !arrayContains(desiredGrants, grant) && currentGrant == "Y"
}

func grantOrRevoke(db *sql.DB, fullUserName string, databaseGrant DatabaseGrants, privilege string, hasPrivilege string) {
	if needGrant(databaseGrant.Grants, privilege, hasPrivilege) {
		fmt.Println("[", fullUserName, "]", databaseGrant.DatabaseName, "grant", privilege)
		sqlExecOrFail(db, "GRANT "+privilege+" ON `"+databaseGrant.DatabaseName+"`.* TO "+fullUserName)
	}
	if needRevoke(databaseGrant.Grants, privilege, hasPrivilege) {
		fmt.Println("[", fullUserName, "]", databaseGrant.DatabaseName, "revoke", privilege)
		sqlExecOrFail(db, "REVOKE "+privilege+" ON `"+databaseGrant.DatabaseName+"`.* FROM "+fullUserName)
	}
}
