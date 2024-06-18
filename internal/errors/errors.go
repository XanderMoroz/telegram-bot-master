package errors

type Error string

func (e Error) Error() string {
	return string(e)
}

const ErrBufferIsFull = Error("Buffer is fulled")
