package network

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isNotFoundErr(err error) bool {
	s, ok := status.FromError(err)
	if ok {
		switch s.Code() {
		case codes.NotFound:
			return true
		case codes.Unknown:
			if strings.Contains(s.Message(), "not found") {
				return true
			}
			return false
		}
	}
	return false
}
