package qmp

import "io"

const (
	MessageTypeSuccess = 1
	MessageTypeError   = 2
	MessageTypeEvent   = 3
)

type JsonValue interface{}
type JsonObject map[string]interface{}
type JsonArray []interface{}

type MessageType byte

// Session represents a connection to a QMP instance
type Session struct {
	c io.ReadWriteCloser

	Greeting      GreetingMessage
	AsyncMessages []AsyncMessage
}

// GreetingMessage is the first message sent by QMP
// It contains information about the QEMU instance
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

// Timestamp is the timestamp representation in QMP messages
type Timestamp struct {
	Seconds      uint64 `json:"seconds"`
	Microseconds uint64 `json:"microseconds"`
}

// SuccessMessage is a successfull response from QMP
type SuccessMessage struct {
	Return JsonValue `json:"return"`
}

// ErrorMessage represents an error sent by QMP
type ErrorMessage struct {
	Error struct {
		Class string `json:"class"`
		Desc  string `json:"desc"`
	} `json:"error"`
}

// AsyncMessage is an asynchronous event sent by QMP
type AsyncMessage struct {
	Event     string     `json:"event"`
	Data      JsonObject `json:"data"`
	Timestamp Timestamp  `json:"timestamp"`
}
