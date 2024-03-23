package error

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

func New(errorCode string, message string) error {
	return &ApplicationError{
		ErrorCode: errorCode,
		Message:   message,
	}
}

type ApplicationError struct {
	ErrorCode string
	Message   string
}

func (e *ApplicationError) Error() string {
	return e.Message
}

// handling timeout from http and g rpc
func IsTimeout(err error) (timeout bool) {
	timeout = os.IsTimeout(err)
	if timeout {
		return
	}

	st, ok := status.FromError(err)
	if !ok {
		return
	}

	if st.Code() == codes.DeadlineExceeded {
		timeout = true
	}

	return
}
