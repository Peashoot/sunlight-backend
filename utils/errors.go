package utils

type UserNotFoundError struct {
	Message string
}

func (e UserNotFoundError) Error() string {
	return e.Message
}

type DBExecuteError struct {
	Message string
}

func (e DBExecuteError) Error() string {
	return e.Message
}

type RecordExistsFoundError struct {
	Message string
}

func (e RecordExistsFoundError) Error() string {
	return e.Message
}
