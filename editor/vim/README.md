# Vim plugin for gotests

This is [Vim](http://www.vim.org/) plugin for `gotests`.

## Setup

### Manual installation

You can create the symlink using `symlinks.sh` like the following:

```console
$ go get -u github.com/tcnksm/gotests
$ sh $GOPATH/src/github.com/tcnksm/gotests/editor/vim/misc/symlink.sh
```

### Using plugin manager

You can add the following line to your .vimrc:

```vim
" vim-plug
Plug 'tcnksm/gotests', { 'rtp': 'editor/vim' }

" Vundle
Plugin 'tcnksm/gotests', { 'rtp': 'editor/vim' }

" NeoBundle
NeoBundle 'tcnksm/gotests', { 'rtp': 'editor/vim' }
```

And then update your packages by running install command.
