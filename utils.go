package log

import (
	"net/url"
	"os"
	"path"
)

// Return the name of the program.
func prog() string {
	return path.Base(os.Args[0])
}

// Intelligently return the transport and address values from a URI.
func uriaddr(uri string) (string, string, error) {
	if uri, err := url.Parse(uri); err == nil {
		host := ""
		if uri.Scheme == "tcp" || uri.Scheme == "udp" {
			host = uri.Host
		} else {
			host = uri.Path
		}
		return uri.Scheme, host, nil
	} else {
		return "", "", err
	}
}
