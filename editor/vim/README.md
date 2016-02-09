# Vim plugin for gotests

This is [Vim](http://www.vim.org/) plugin for `gotests`.

## Install & Configuration

To use this plugin, you need to install `gotests` binary itself first. To install, run `go get`:

```sh
$ go get github.com/tcnksm/gotests
```

To install, use well known plugin manager. You can add the following lines and execute the plugin's install:

- [vim-plug](https://github.com/junegunn/vim-plug)
    - `Plug 'tcnksm/gotests', { 'rtp': 'editor/vim' }`
- [NeoBundle](https://github.com/Shougo/neobundle.vim)
    - `NeoBundle 'tcnksm/gotests', { 'rtp': 'editor/vim' }`
- [Vundle](https://github.com/gmarik/vundle)
    - `Plugin 'tcnksm/gotests', { 'rtp': 'editor/vim' }` 

If you don't use any plugin manager, run the following command and link to the plugin path,

```sh
$ sh $GOPATH/src/github.com/tcnksm/gotests/editor/vim/symlink.sh
```

## Author

[@htm](https://github.com/hfm)
