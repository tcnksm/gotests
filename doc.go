// DON"T EDIT THIS FILE
// THIS IS GENERATED VIA GO GENERATE

/*
gotests is tool to generate Go test functions from
the given source code. 

https://github.com/tcnksm/gotests

Usage:

  gotests [options] PATH ...

Options:

  -diff, -d      Display diffs instead of rewriting files.

  -write, -w     Write result to target file instead of stdout.
                 For example, if source file name is 'A.go',
                 target file will be 'A_test.go'.

  -list, -l      List test files to be updated/generated.

  -i             Include unexport function/method for generating target.

  -reverse, -r   (experimental) Allow to provide test file instead of source file.
                 By default, gotests expects source file PATH provided.
                 With this flag, the test file can be given. 
                 For example, you can provide 'A_test.go' instead of 'A.go'.
                 This flag is useful for editor integration.

*/
package main
