# Overview
util-cli is a util that can transform markdown'binary to file's image

# Installation
## Download package
firstly download and unpack the release package then execute the script
```shell script
sudo bash install.sh
```

## shell
```shell script
curl -s https://raw.githubusercontent.com/chnherb/util-cli/master/install-release.sh | sudo bash -s -- -r chnherb/util-cli 
```

# Usage
```shell script
util-cli -h     # help
util-cli -v     # show more detail
```

## imgbase64
process image
```shell script
util-cli imgbase64 -h
```

### replace
replace image of base64 encoding to file's image
e.g.
```shell script
util-cli imgbase64 replace -h
util-cli imgbase64 replace --chapter=skywalking211211 --rewrite false
```