package amocrm

import "fmt"

func (e *AmoError) Error() string {
	return fmt.Sprintf("%s: %s", AmoErrorTypeMap[e.ErrorCode], e.ErrorDetail)
}
