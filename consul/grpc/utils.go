package grpc

import (
	"errors"
	"regexp"
)

var (
	regexConsul, _ = regexp.Compile("^([A-z0-9._-]+)(:[0-9]{1,5})?/([A-z_]+)$")
)

func parseTarget(target string) (host, port, name string, err error) {
	if target == "" {
		return "", "", "", errors.New("consul resolver: missing address")
	}

	if !regexConsul.MatchString(target) {
		return "", "", "", errors.New("consul resolver: invalid uri")
	}

	groups := regexConsul.FindStringSubmatch(target)
	host = groups[1]
	port = groups[2]
	name = groups[3]
	if port == "" {
		port = "8500"
	}
	return host, port, name, nil
}
