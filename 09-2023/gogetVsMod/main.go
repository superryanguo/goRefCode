package main

import (
	"fmt"
	"os"

	"golang.org/x/text/language"
)

func main() {
	for _, arg := range os.Args[1:] {
		tag, err := language.Parse(arg)
		if err != nil {
			fmt.Printf("%s: error: %v\n", arg, err)
		} else if tag == language.Und {
			fmt.Printf("%s: undefined\n", arg)
		} else {
			fmt.Printf("%s: tag %s\n", arg, tag)
		}
	}
}

//govulncheck ./...
//在go1.18之后，在gomod模式下可以用go get来指定某个包的版本

//Overview
//Starting in Go 1.17, installing executables with go get is deprecated. go install may be used instead.
//In Go 1.18, go get will no longer build packages; it will only be used to add, update, or remove dependencies in go.mod. Specifically, go get will always act as if the -d flag were enabled.

//From <https://go.dev/doc/go-get-install-deprecation>
//go: cannot determine module path for source directory//只需要放在gopath里面就不会有这个问题
///repo/eggorri/5g/mywork/vulncheck (outside GOPATH, module path must be specified)//需要的话就给个模块名

//Example usage:
//'go mod init example.com/m' to initialize a v0 or v1 module
//'go mod init example.com/m/v2' to initialize a v2 module

//Run 'go help mod init' for more information.
//go install example.com/cmd
//To install an executable while ignoring the current module, use go install with a version suffix like @v1.2.3 or @latest, as below. When used with a version suffix, go install does not read or update the go.mod file in the current directory or a parent directory.
//# Install a specific version.
//go install example.com/cmd@v1.2.3
//# Install the highest available version.
//go install example.com/cmd@latest

//From <https://go.dev/doc/go-get-install-deprecation>
