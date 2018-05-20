package socket

import (
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// Parse produces a struct suitable for net.Dial
// given a string representing a socket address to bind or connect to
func Parse(value string) (string, string) {
	value = strings.TrimSpace(value)

	// If value begins with "unix:" then we are a Unix domain socket
	if strings.Index(value, "unix:") == 0 {
		return "unix", strings.TrimSpace(value[5:])
	}

	// If value is a port number, prepend a colon
	if _, err := strconv.Atoi(value); err == nil {
		return "tcp", ":" + value
	}

	// If the value is a host with a port, return exact.
	if host, port, err := net.SplitHostPort(value); err == nil {
		if host == "*" {
			host = ""
		}
		return "tcp", net.JoinHostPort(host, port)
	}

	// Otherwise, assume that the input is just a hostname or IP
	// However, if it is an IPv6 address, we need to remove brackets for JoinHostPort to work.
	value = strings.Trim(value, "[]")

	// Use default port of 80
	return "tcp", net.JoinHostPort(value, "80")
}

// Listen takes the same arguments as the net.Listen function but includes
// the mode with which to set the file permissions if creating a Unix socket.
//
// The socket path is appended with ".tmp" before it is created and later renamed.
//
// This function calls syscall.Umask before the socket is created in order to
// completely restrict the socket before it is restored again with os.Chmod.
// Therefore, this function is not thread-safe as Umask is non-reentrant.
//
// Additionally, this function will make sure to remove the temporary socket path
// before creating it.
func Listen(network, address string, mode os.FileMode) (net.Listener, error) {
	switch network {
	case "unix", "unixgram", "unixpacket":
		// Create temporary socket with no permissions
		addressTmp := address + ".tmp"
		listener, err := func() (net.Listener, error) {
			os.Remove(addressTmp)
			oldmask := syscall.Umask(0777)
			defer syscall.Umask(oldmask)
			return net.Listen(network, addressTmp)
		}()
		if err != nil {
			return nil, err
		}

		// Set desired socket permissions
		if err := os.Chmod(addressTmp, mode); err != nil {
			listener.Close()
			return nil, err
		}

		// Move socket to its final resting place
		if err := os.Rename(addressTmp, address); err != nil {
			listener.Close()
			return nil, err
		}

		return listener, err
	}

	return net.Listen(network, address)
}
