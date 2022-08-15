#!/usr/bin/env bash
# set -x
GITHUB_REPO=""
RELEASE_VERSION=""
INSTALL_DIR=""
LINK_SOURCE=""
LINK_TARGET=""
USAGE_STR="Usage: $0 -r <github-repo> -v <release-version> -d <install-dir> -s <link-source> -t <link-target>"
while getopts "r:v:d:l:t:e:h" opt; do
    case $opt in
    h)
        echo $USAGE_STR
        exit 0
        ;;
    r)
        GITHUB_REPO="$OPTARG"
        ;;
    v)
        RELEASE_VERSION="$OPTARG"
        ;;
    d)
        INSTALL_DIR="$OPTARG"
        ;;
    s)
        LINK_SOURCE="$OPTARG"
        ;;
    t)
        LINK_TARGET="$OPTARG"
        ;;
    \?)
        echo "Invalid option: -$OPTARG" >&2
        echo $USAGE_STR
        exit
        ;;
    :)
        echo "Option -$OPTARG requires an argument." >&2
        echo $USAGE_STR
        exit 1
        ;;
    esac
done

REPO_NAME=$(echo $GITHUB_REPO | cut -d'/' -f2)
if [ x"$RELEASE_VERSION" = x"" ] ;then
  RELEASE_VERSION="latest"
fi
if [ x"$INSTALL_DIR" = x"" ]; then
    INSTALL_DIR="/opt/${SUDO_USER}/${REPO_NAME}"
fi

uNames=`uname -s`
osName=${uNames: 0: 5}
if [[ "$osName" == "Darwi" ]]; then
  BIN_NAME=$REPO_NAME"_darwin_amd64"
elif [[ "$osName" == "Linux" && $(uname -m) == "x86_64" ]]; then
  BIN_NAME=$REPO_NAME"_linux_amd64"
#elif [[ "$osName" == "MINGW" ]]; then
#  BIN_NAME=$BIN_NAME_ORIGIN"_windows_amd64"
else
  echo "failed! ${GITHUB_REPO} only supports amd64 Mac, Linux and Windows"
  exit 1
fi
if [ x"$GITHUB_REPO" = x"" ] ;then
  echo "-r is required"
  exit 1
fi

if [ x"$LINK_SOURCE" = x"" ] ;then
  LINK_SOURCE=${INSTALL_DIR}/${BIN_NAME}
fi
if [ x"$LINK_TARGET" = x"" ] ;then
  LINK_TARGET="/usr/local/bin/${REPO_NAME}"
fi

install(){
  # curl -s https://api.github.com/repos/chnherb/util-cli/releases/latest | grep browser_download_url | grep -i Darwin | cut -d'"' -f4
  download_url=$(curl -k -s https://api.github.com/repos/${GITHUB_REPO}/releases/${RELEASE_VERSION:-latest} | grep browser_download_url | grep -i ${osName} | cut -d'"' -f4)
  if [ x"$download_url" = x"" ];then
    echo "Download URL is invalid!"
    exit 1
  fi
  echo "GITHUB_REPO: $GITHUB_REPO, Version: ${RELEASE_VERSION:-latest} installing ..."
  echo
  echo "Release Download URL: $download_url"
  tmp_file_path="/tmp/$(echo $download_url | cut -d'/' -f9)"
  rm -rf $tmp_file_path
  wget -q --show-progress --progress=dot:mega $download_url -O $tmp_file_path
  echo
  if [ $? -ne 0 ] ;then
    echo "Release Download error"
    exit 1
  fi
  if test -d $INSTALL_DIR ;then
    echo "Warning: install dir $INSTALL_DIR has existed, will been override"
    rm -rf $INSTALL_DIR/*
  fi
  mkdir -p $INSTALL_DIR
  tar -zxf $tmp_file_path -C $INSTALL_DIR
  # sudo chown -R ${SUDO_USER}:${SUDO_USER} $INSTALL_DIR
  rm -rf $tmp_file_path
  echo "Install success in: $INSTALL_DIR"
  echo
}

link(){
  if [ x"$LINK_SOURCE" != x"" ] && [ x"$LINK_TARGET" != x"" ] ;then
    if test -f $LINK_TARGET || test -L $LINK_TARGET ;then
      echo "Warning: file $LINK_TARGET has exist, will been override"
      sudo rm -rf $LINK_TARGET
    fi
    sudo ln -s $LINK_SOURCE $LINK_TARGET
    if [ $? -eq 0 ] ;then
      echo "success: ln -s $LINK_SOURCE $LINK_TARGET"
    else
      echo "Error: ln -s $LINK_SOURCE $LINK_TARGET"
      exit 1
    fi
  fi
  echo
}
install
link

