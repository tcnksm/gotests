# Gotg

`gotg` generates Go program test code from the given source code.

## Install

To install, use `go get`:

```bash
$ go get -d github.com/tcnksm/gotg
```

## Example

```bash
$ cat test-fixtures/basic.go
package basic

func DoSomething() error {
    return nil
}

type User struct {
    Name, Password string
}

func (u *User) Validate() error {
    return nil
}

$ cat test-fixtures/basic_test.go
package basic

import "testing"

func TestDoSomething(t *testing.T) {}

$ ./bin/goatg -diff test-fixtures/basic.go
diff basic_test.go
--- /var/folders/hk/t6xpt_j974gfnv_nd3t79x3h2b8vpv/T/gotg350053338      2016-01-19 14:20:59.000000000 +0900
+++ /var/folders/hk/t6xpt_j974gfnv_nd3t79x3h2b8vpv/T/gotg937594737      2016-01-19 14:20:59.000000000 +0900
@@ -3,3 +3,5 @@
import "testing"

func TestDoSomething(t *testing.T) {}
+func TestValidate(t *testing.T) {
+}
```

## Contribution

1. Fork ([https://github.com/tcnksm/gotg/fork](https://github.com/tcnksm/gotg/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[tcnksm](https://github.com/tcnksm)
