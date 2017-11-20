# About

This is an application to update the databases and the users permissions by applying the config from a file.

# Local Usage


## Compile

`./create-local-release.sh`

## Configure

```
TMPDIR=$(mktemp -d)

# Create configuration file
vi $TMPDIR/manager-config.json
{
	"admin" : {
		"name" : "root",
		"password" : "qwerty"
	},
	"databases" : [
		"first_db",
		"second_db"
	],
	"usersToIgnore" : [
		{
			"name" : "root",
			"host" : "localhost"
		},
		{
			"name" : "debian-sys-maint",
			"host" : "localhost"
		},
		{
			"name" : "mysql.session",
			"host" : "localhost"
		},
		{
			"name" : "mysql.sys",
			"host" : "localhost"
		}
	],
	"usersPermissions" : [
		{
			"name" : "alice",
			"host" : "%",
			"password" : "qwerty",
			"databaseGrants" : [
				{
					"databaseName" : "first_db",
					"grants" : ["CREATE", "ALTER", "DROP"]
				}
			]
		},
		{
			"name" : "bob",
			"host" : "%",
			"password" : "*qwerty",
			"databaseGrants" : [
				{
					"databaseName" : "first_db",
					"grants" : ["SELECT", "INSERT", "UPDATE", "DELETE"]
				}
			]
		}
	]
}

```

## Execute

To see the help:
`./build/gopath/build/out/mysql-manager`

To execute:
`./build/gopath/build/out/mysql-manager 127.0.0.1:3306 $TMPDIR/manager-config.json`

## Important points

* You can ignore some users. This is mostly to ignore the admin (root) account that you manage the password in another way.
* Any databases that is not listed will be dropped.
* Any user that is not specifically ignored or listed in `usersPermissions` will be dropped.
* There is no check that the databases in `usersPermissions` are listed in `databases`.
* There is no check that the `databaseGrants` for a user is not using multiple times the same database. The last one listed for a specific database is what the user will get.

# Create release

`./create-public-release.sh`

# Use with debian

```bash
echo "deb https://dl.bintray.com/foilen/debian stable main" | sudo tee /etc/apt/sources.list.d/foilen.list
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61
sudo apt update
sudo apt install mysql-manager
```
