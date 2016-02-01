# Emacs package for gotests

This is [Emacs](http://www.gnu.org/software/emacs/) package for `gotests`.

*NOTE*: I'm not good at emacs package development. So this plugin should not well written. If you are good at emacs plugin please send PR 🙇

## Install

To install this,

```bash
$ go get -d github.com/tcnksm/gotests
```

## Configuration

Add the following line in your `init.el` file,

```lisp
(load-file (concat (getenv "GOPATH") "/src/github.com/tcnksm/gotests/editor/emacs/gotests.el"))
```
