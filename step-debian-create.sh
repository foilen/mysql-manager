#!/bin/bash

set -e

RUN_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $RUN_PATH

echo ----[ Create .deb ]----
DEB_FILE=mysql-manager_${VERSION}_amd64.deb
DEB_PATH=$RUN_PATH/build/debian_out/mysql-manager
rm -rf $DEB_PATH
mkdir -p $DEB_PATH $DEB_PATH/DEBIAN/ $DEB_PATH/usr/local/bin/

cat > $DEB_PATH/DEBIAN/control << _EOF
Package: mysql-manager
Version: $VERSION
Maintainer: Foilen
Architecture: amd64
Description: This is an application to update the databases and the users permissions by applying the config from a file.
_EOF

cp -rv DEBIAN $DEB_PATH/
cp -rv build/bin/* $DEB_PATH/usr/local/bin/

cd $DEB_PATH/..
dpkg-deb --no-uniform-compression --build mysql-manager
mv mysql-manager.deb $DEB_FILE
