# socket

Go package to parse socket names described in [Nginx convention][nginx-listen].

[![GoDoc](https://godoc.org/github.com/oftn-oswg/socket?status.svg)](https://godoc.org/github.com/oftn-oswg/socket)

## Installation and usage

```sh
go get github.com/oftn-oswg/socket
```

```go
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/oftn-oswg/socket"
)

// This is a simple echo server which creates a socket
// based on the arguments provided to the command. The socket can be
// described with a syntax that nginx, the web server, accepts for its listen directive.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <socket> [<mode>]\n", os.Args[0])
		os.Exit(1)
	}

	mode := os.FileMode(0660)
	if len(os.Args) >= 3 {
		num, err := strconv.ParseInt(os.Args[2], 0, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file mode: %s\n", err)
			os.Exit(1)
		}
		mode = os.FileMode(num)
	}

	network, address := socket.Parse(os.Args[1])
	listener, err := socket.Listen(network, address, mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating socket: %s\n", err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		go func() {
			defer conn.Close()
			_, err := io.Copy(conn, conn)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()
	}
}

```

[nginx-listen]: http://nginx.org/en/docs/http/ngx_http_core_module.html#listen
