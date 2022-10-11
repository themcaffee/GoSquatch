[_metadata_:title]:- "Github Actions"
[_metadata_:layout]:- "index"

# Live Server

GoSquatch can be used to run a local server that will rebuild your source directory on any file changes. This lets you develop you site
locally before you show it to the world.

## Installation

Run the following to install gosquatch through apt:
```
curl -s --compressed "https://www.mitchmcaffee.com/GoSquatch/ppa/KEY.gpg" | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/gosquatch.gpg > /dev/null
sudo curl -s --compressed -o /etc/apt/sources.list.d/gosquatch_list_file.list "https://www.mitchmcaffee.com/GoSquatch/ppa/gosquatch_list_file.list"
sudo apt update
sudo apt install gosquatch
```

## Usage

```
gosquatch -live-server -src-dir=./ -port=8080
```

Then visit your site at [http://localhost:8080](http://localhost:8080)

### Options

`-src-dir`: The location of your source directory

`-port`: The port the live server runs on

`-live-server`: Runs the live-server. If this is not included, GoSquatch runs in build mode.

## Updating GoSquatch

Updating your local installation of GoSquatch is just like any other apt package:

```
sudo apt update
sudo apt upgrade
```