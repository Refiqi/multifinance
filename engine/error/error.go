package error

// TransactionError handles custom error response when loan limit is exceeded
type TransactionError struct {
	Message string
	Code    int
	Limit   float64
	Current float64
	New     float64
}

func (e *TransactionError) Error() string {
	return e.Message
}