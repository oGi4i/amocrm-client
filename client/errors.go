package client

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEmptyLogin              Error = "empty_login"
	ErrEmptyAPIHash            Error = "empty_api_hash"
	ErrEmptyResponse           Error = "empty_response"
	ErrInvalidContactID        Error = "invalid_contact_id"
	ErrInvalidLeadID           Error = "invalid_lead_id"
	ErrInvalidPipelineID       Error = "invalid_pipeline_id"
	ErrInvalidPipelineStatusID Error = "invalid_pipeline_status_id"
	ErrInvalidTaskID           Error = "invalid_task_id"
	ErrInvalidTaskResult       Error = "invalid_task_result"
)
