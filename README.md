# About

This is an application to update the databases and the users permissions by applying the config from a file.

# Local Usage


## Compile

`./create-local-release.sh`

The file is then in `build/bin/mysql-manager`

## Create a mariadb database

```
INSTANCE=test-manager

docker run \
  --name $INSTANCE \
  -e MYSQL_ROOT_PASSWORD=qwerty \
  -e DBNAME=noname \
  -p 3306:3306 \
  -d mariadb:10.4.8
```

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
			"name" : "root",
			"host" : "%"
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
`./build/bin/mysql-manager`

To execute:
`./build/bin/mysql-manager 127.0.0.1:3306 $TMPDIR/manager-config.json`

## Important points

* You can ignore some users. This is mostly to ignore the admin (root) account that you manage the password in another way.
* Any databases that is not listed will be dropped.
* Any user that is not specifically ignored or listed in `usersPermissions` will be dropped.
* There is no check that the databases in `usersPermissions` are listed in `databases`.
* There is no check that the `databaseGrants` for a user is not using multiple times the same database. The last one listed for a specific database is what the user will get.

## Delete the mariadb database

```
docker rm -f $INSTANCE
```

# Create release

`./create-public-release.sh`

That will show the latest created version. Then, you can choose one and execute:
`./create-public-release.sh X.X.X`

# Use with debian

Get the version you want from https://deploy.foilen.com/mysql-manager/ .

```bash
wget https://deploy.foilen.com/mysql-manager/mysql-manager_X.X.X_amd64.deb
sudo dpkg -i mysql-manager_X.X.X_amd64.deb
```
