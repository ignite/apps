package network

import (
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
			if s.Message() == "not found: unknown request" {
				return true
			}
			return false
		}
	}
	return false
}
