package client

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEmptyLogin    Error = "empty_login"
	ErrEmptyAPIHash  Error = "empty_api_hash"
	ErrEmptyResponse Error = "empty_response"
)
