package chat

import "bytes"

const (
	START = "/start"
	STOP = "/stop"
)

func IsStartMessage(message []byte) bool {
	return equal(message, []byte(START))
}

func IsStopMessage(message []byte) bool {
	return equal(message, []byte(STOP))
}

func equal(message, expected []byte) bool {
	return bytes.Equal(message, expected)
}
