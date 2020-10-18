package client

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEmptyClientID        Error = "empty_client_id"
	ErrEmptyClientSecret    Error = "empty_client_secret"
	ErrEmptyRefreshToken    Error = "empty_refresh_token"
	ErrInvalidAuthTokenType Error = "invalid_auth_token_type"

	ErrEmptyResponse Error = "empty_response"

	ErrInvalidContactID        Error = "invalid_contact_id"
	ErrInvalidLeadID           Error = "invalid_lead_id"
	ErrInvalidPipelineID       Error = "invalid_pipeline_id"
	ErrInvalidPipelineStatusID Error = "invalid_pipeline_status_id"
	ErrInvalidTaskID           Error = "invalid_task_id"
	ErrInvalidTaskResult       Error = "invalid_task_result"
)
