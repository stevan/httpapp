# Simple HTTP Applications in Go

This is a small bit of Go code that builds upon the
standard net/http framework to help in writing HTTP
applications in the style of WSGI using middleware.

## Installation

With Go and git installed you can do:

    go get github.com/stevan/httpapp

and `go` will download, compile, and install the
package into your `$GOROOT` directory hierarchy.

Or you can just directly import it in your project:

    import "github.com/stevan/httpapp"

and then when you run `go build`, the package will
be downloaded and installed automatically.

