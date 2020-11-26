package repository

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	netrc         []netrcLine
	readNetrcOnce sync.Once
)

// See https://github.com/golang/go/blob/master/src/cmd/go/internal/web2/web.go
// for implementation
// Temporary netrc reader until https://github.com/golang/go/issues/31334 is solved
type netrcLine struct {
	machine  string
	login    string
	password string
}

func parseNetrc(data string) []netrcLine {
	// See https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html
	// for documentation on the .netrc format.
	var nrc []netrcLine
	var l netrcLine
	inMacro := false
	for _, line := range strings.Split(data, "\n") {
		if inMacro {
			if line == "" {
				inMacro = false
			}
			continue
		}

		f := strings.Fields(line)
		i := 0
		for ; i < len(f)-1; i += 2 {
			// Reset at each "machine" token.
			// “The auto-login process searches the .netrc file for a machine token
			// that matches […]. Once a match is made, the subsequent .netrc tokens
			// are processed, stopping when the end of file is reached or another
			// machine or a default token is encountered.”
			switch f[i] {
			case "machine":
				l = netrcLine{machine: f[i+1]}
			case "login":
				l.login = f[i+1]
			case "password":
				l.password = f[i+1]
			case "macdef":
				// “A macro is defined with the specified name; its contents begin with
				// the next .netrc line and continue until a null line (consecutive
				// new-line characters) is encountered.”
				inMacro = true
			}
			if l.machine != "" && l.login != "" && l.password != "" {
				nrc = append(nrc, l)
				l = netrcLine{}
			}
		}

		if i < len(f) && f[i] == "default" {
			// “There can be only one default token, and it must be after all machine tokens.”
			break
		}
	}

	return nrc
}

func netrcPath() (string, error) {
	if env := os.Getenv("NETRC"); env != "" {
		return env, nil
	}

	dir := os.Getenv("HOME")

	base := ".netrc"
	if runtime.GOOS == "windows" {
		base = "_netrc"
	}
	return filepath.Join(dir, base), nil
}

// readNetrc parses a user's netrc file, ignoring any errors that occur.
func readNetrc() {
	path, err := netrcPath()
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	netrc = parseNetrc(string(data))
}

// addAuthFromNetrc uses basic authentication on go-get requests
// for private repositories.
func AddAuthFromNetrc(rawurl string, req *http.Request) *http.Request {
	readNetrcOnce.Do(readNetrc)
	u, err := url.Parse(rawurl)
	if err != nil {
		return req
	}

	for _, m := range netrc {
		if u.Hostname() == m.machine {
			req.SetBasicAuth(m.login, m.password)
			break
		}
	}

	return req
}
