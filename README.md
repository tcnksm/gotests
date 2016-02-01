# gotests

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[license]: /LICENSE
[godocs]: http://godoc.org/github.com/tcnksm/gotests

`gotests` generates Go test functions from the given source code.

Given `A.go` file, it analyzes test functions in `A_test.go` and adds functions which are not defined in that file. For example, if a function `DoSomething()` is defined in `A.go` but `TestDoSomething()` is not in `A_test.go`, it adds that function to `A_test.go`. By default, it only checks the exported functions (its name starts with upper case). Given a directory, it operates on all `*.go` files in that directory. By default, `gotests` prints the updated test sources to standard output.

I hope this tool would be a new friend of Gophers like `gofmt` or `gorename`. 

## Editor

`gotests` works well with your favorite editor like `gofmt` does. The following demo shows using `gotests` from Emacs. The left display shows the source codes (`A.go`) and the right shows test source codes (`A_test.go`). It generates and adds test functions on the right codes,

![demo](/gotests.gif)

`gotests.el` used by this demo is available in [`editor/emacs`](/editor/emacs) directory (I'm not good at emacs plugin development. So this plugin should not well written. If you are good at emacs plugin please send PR 🙇 ).

A plugin PR for the other editor is welcome.

## Install

To install, use `go get`:

```bash
$ go get -d github.com/tcnksm/gotests
```

## Usage

The usage:

```bash
gotests [options] PATH ...
```

Available options:

```bash
-diff, -d    Display diffs instead of rewriting files.

-write, -w   Write result to target file instead of stdout. For example, if source file name is 'A.go',             
             target file would be 'A_test.go'.

-list, -l    List test files to be updated/generated.

-i           Include unexport function/method for generating target.
```

## Rule

The followings are the current basic rule of how to generate test functions, 

- Test functions for `A.go` are added to `A_test.go`
- Generating test function name of function is `Test{{ title .Name }}` 
- Generating test function name of method is `Test{{ title .ReceiverTypeName }}_{{ title .Name }}`

(In future, this rule will be configurable from command line option)

## Contribution

1. Fork ([https://github.com/tcnksm/gotests/fork](https://github.com/tcnksm/gotests/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[Taichi Nakashima](https://github.com/tcnksm)
