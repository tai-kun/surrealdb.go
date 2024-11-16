package surrealdb

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	transformEndpointDefault  string = ""
	transformEndpointAuto     string = "auto"
	transformEndpointPreserve string = "preserve"
)

func processEndpoint(endpoint, transform string) (*url.URL, error) {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		err = fmt.Errorf("failed to process endpoint: %w", err)
		return nil, err
	}

	switch transform {
	case transformEndpointDefault, transformEndpointAuto:
		if !strings.HasSuffix(u.Path, "/rpc") {
			if strings.HasPrefix(u.Path, "/") {
				u.Path += "/"
			}
			u.Path += "rpc"
		}
	case transformEndpointPreserve:
		// pass
	default:
		err := fmt.Errorf(
			"failed to process endpoint: invalid transform=%s",
			strconv.Quote(transform),
		)
		return nil, err
	}

	return u, nil
}
