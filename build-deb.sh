#!/bin/bash

# This script is used to build the debian package for the project
# It is assumed that the project is already built and the binaries are in the bin folder

# The script takes the following arguments
# 1. The version of the project

packageName="gosquatch"
architecture="amd64"
maintainer="Mitch McAffee <squatch@mitchmcaffee.com>"
Description="Convert markdown to a static site"
gpgEmail="gpg@mitchmcaffee.com"

packageDir="gosquatch_$1-1_amd64" # The name of the debian package directory

go build -o gosquatch .

# Create the internal folder structure
mkdir $packageDir
mkdir -p $packageDir/usr/local/bin
cp gosquatch $packageDir/usr/local/bin

# Create the control file
mkdir -p $packageDir/DEBIAN
touch $packageDir/DEBIAN/control
echo "Package: $packageName" >> $packageDir/DEBIAN/control
echo "Version: $1" >> $packageDir/DEBIAN/control
echo "Architecture: $architecture" >> $packageDir/DEBIAN/control
echo "Maintainer: $maintainer" >> $packageDir/DEBIAN/control
echo "Description: $Description" >> $packageDir/DEBIAN/control

# Build the deb package
dpkg-deb --build --root-owner-group $packageDir

# Cleanup the build folder
rm -rf $packageDir

# Cleanup the dist folder
rm docs/ppa/Packages
rm docs/ppa/Packages.gz
rm docs/ppa/Release
rm docs/ppa/Release.gpg
rm docs/ppa/InRelease

# Move to the docs folder for distribution
mv $packageDir.deb docs/ppa/$packageDir.deb

cd docs/ppa

# Create the signature
dpkg-scanpackages --multiversion . > Packages
gzip -k -f Packages

# Sign the package
apt-ftparchive release . > Release
gpg --default-key "${gpgEmail}" -abs -o - Release > Release.gpg
gpg --default-key "${gpgEmail}" --clearsign -o - Release > InRelease