#!/bin/bash
#!/usr/bin/env bash

# set -x
BIN_NAME_ORIGIN="util-cli"
BIN_NAME=""

#if [[ "$(uname)" == "Darwin" && $(uname -m) == "x86_64" ]]; then
#  BIN_NAME=$BIN_NAME_ORIGIN"_darwin_amd64"
#elif [[ "$(expr substr $(uname -s) 1 5)" == "Linux" && $(uname -m) == "x86_64" ]]; then
#  BIN_NAME=$BIN_NAME_ORIGIN"_linux_amd64"
#elif [[ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]]; then
#  BIN_NAME=$BIN_NAME_ORIGIN"_windows_amd64"
uNames=`uname -s`
osName=${uNames: 0: 5}
install_path="/usr/local/bin/"
if [[ "$osName" == "Darwi" ]]; then
  BIN_NAME=$BIN_NAME_ORIGIN"_darwin_amd64"
elif [[ "$osName" == "Linux" && $(uname -m) == "x86_64" ]]; then
  BIN_NAME=$BIN_NAME_ORIGIN"_linux_amd64"
#elif [[ "$osName" == "MINGW" ]]; then
#  BIN_NAME=$BIN_NAME_ORIGIN"_windows_amd64"
#  install_path="C:/util-cli/"
#  sudo mkdir -p $install_path
else
  echo "failed! ${BIN_NAME_ORIGIN} only supports amd64 Mac, Linux and Windows"
  exit 1
fi

echo "install path: ${install_path}${BIN_NAME_ORIGIN}"
sudo cp -rf "$(pwd)/${BIN_NAME}" ${install_path}${BIN_NAME_ORIGIN}
# sudo rm -rf "/usr/local/bin/$BIN_NAME"
# sudo ln -s "$(pwd)/${BIN_NAME}" "/usr/local/bin/${BIN_NAME_ORIGIN}"
