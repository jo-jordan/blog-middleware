package common

import (
	"net/url"
	"strings"
)

func GetRepoName(s string) (string, error) {
	u, err := url.Parse(s)
	if nil != err {
		return "", err
	}

	path := u.Path
	path = path[strings.LastIndex(path, "/") + 1:]
	if strings.Contains(path, ".") {
		path = path[:strings.LastIndex(path, ".")]
	}

	return path, nil
}
