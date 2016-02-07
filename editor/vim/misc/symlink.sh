#!/bin/sh
current_path=$(cd $(dirname $0) && pwd)
mkdir -p "${HOME}/.vim/ftplugin/go"
ln -sf "${current_path}/ftplugin/go/gotests.vim" "${HOME}/.vim/ftplugin/go/"
