package socket

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		Input   string
		Network string
		Address string
	}{
		// Examples from http://nginx.org/en/docs/http/ngx_http_core_module.html#listen
		// listen 127.0.0.1:8000;
		// listen 127.0.0.1;
		// listen 8000;
		// listen *:8000;
		// listen localhost:8000;
		// listen [::]:8000;
		// listen [::1];
		// listen unix:/var/run/nginx.sock;

		{"127.0.0.1:8000", "tcp", "127.0.0.1:8000"},
		{"127.0.0.1", "tcp", "127.0.0.1:80"}, // "If only address is given, the port 80 is used."
		{"8000", "tcp", ":8000"},
		{"*:8000", "tcp", ":8000"},
		{"localhost:8000", "tcp", "localhost:8000"},
		{"localhost", "tcp", "localhost:80"},
		{"[::]:8000", "tcp", "[::]:8000"},
		{"[::1]", "tcp", "[::1]:80"},
		{"unix:/var/run/nginx.sock", "unix", "/var/run/nginx.sock"},
	}

	for _, test := range tests {
		network, address := Parse(test.Input)
		if network != test.Network || address != test.Address {
			t.Errorf("For input %q: Expected %q, %q but got %q, %q",
				test.Input, test.Network, test.Address, network, address)
		}
	}
}
