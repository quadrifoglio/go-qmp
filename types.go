package qmp

import "io"

const (
	MessageTypeSuccess = 1
	MessageTypeError   = 2
	MessageTypeEvent   = 3
)

type JsonObject map[string]interface{}
type JsonArray []interface{}

type MessageType byte

type Session struct {
	c io.ReadWriteCloser

	Greeting      GreetingMessage
	AsyncMessages []AsyncMessage
}

type GreetingMessage struct {
	QMP struct {
		Version struct {
			Qemu struct {
				Micro int `json:"micro"`
				Minor int `json:"minor"`
				Major int `json:"major"`
			} `json:"qemu"`

			Package string `json:"package"`
		} `json:"version"`

		Capabilities []interface{} `json:"capabilities"`
	}
}

type Timestamp struct {
	Seconds      uint64 `json:"seconds"`
	Microseconds uint64 `json:"microseconds"`
}

type SuccessMessage struct {
	Return JsonObject `json:"return"`
}

type ErrorMessage struct {
	Error struct {
		Class string `json:"class"`
		Desc  string `json:"desc"`
	} `json:"error"`
}

type AsyncMessage struct {
	Event     string     `json:"event"`
	Data      JsonObject `json:"data"`
	Timestamp Timestamp  `json:"timestamp"`
}
