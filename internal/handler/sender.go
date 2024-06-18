package handler

import (
	custErr "telegram-bot/internal/errors"
	"time"
)

const TIMEOUT = 5

func SendMsg(res chan<- string, msg string) error {
	select {
	case res <- msg:
		return nil
	case <-time.After(time.Second * TIMEOUT):
		return custErr.ErrBufferIsFull
	}
}
