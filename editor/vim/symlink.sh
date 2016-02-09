#!/bin/sh
# This script is used for linking plugin to plugin direcoty

DIR=$(pwd)

# Create plugin directory
mkdir -p "${HOME}/.vim/ftplugin/go"

# Link to plugin directory
echo "--> Link ${DIR}/ftplugin/go/gotests.vim to ${HOME}/.vim/ftplugin/go/"
ln -sf "${DIR}/ftplugin/go/gotests.vim" "${HOME}/.vim/ftplugin/go/"
