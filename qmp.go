package qmp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// Open creates a connection to QMP at the specified address using the named protocol
// Protocol can be tcp or unix
func Open(proto, addr string) (Session, error) {
	var s Session

	c, err := net.Dial(proto, addr)
	if err != nil {
		return s, err
	}

	s.c = c

	var g GreetingMessage

	err = json.NewDecoder(c).Decode(&g)
	if err != nil {
		return s, fmt.Errorf("qmp: decode json greeting: %s", err)
	}

	s.Greeting = g

	_, err = s.Command("qmp_capabilities", nil)
	if err != nil {
		return s, err
	}

	return s, nil
}

// Close closes the connection
func (s Session) Close() error {
	return s.c.Close()
}

// Command sends a command and returns the response from QEMU.
func (s Session) Command(command string, arguments map[string]string) (JsonValue, error) {
	cmd := make(JsonObject)
	cmd["execute"] = command

	if arguments != nil {
		cmd["arguments"] = arguments
	}

	err := json.NewEncoder(s.c).Encode(cmd)
	if err != nil {
		return nil, fmt.Errorf("qmp: encode json: %s", err)
	}

	t, data, err := s.read()
	if err != nil {
		return nil, err
	}

	if t == MessageTypeSuccess {
		var m SuccessMessage

		err := json.Unmarshal(data, &m)
		if err != nil {
			return nil, fmt.Errorf("qmp: decode json error: %s", err)
		}

		return m.Return, nil
	} else if t == MessageTypeError {
		var e ErrorMessage

		err := json.Unmarshal(data, &e)
		if err != nil {
			return nil, fmt.Errorf("qmp: decode json error: %s", err)
		}

		return nil, fmt.Errorf("qmp error %s: %s", e.Error.Class, e.Error.Desc)
	}

	return nil, fmt.Errorf("qmp: unknown message")
}

// HumanMonitorCommand sends a HMP command to QEMU via the QMP protocol
func (s Session) HumanMonitorCommand(command string) (JsonValue, error) {
	return s.Command("human-monitor-command", map[string]string{"command-line": command})
}

func (s *Session) read() (MessageType, []byte, error) {
	scanner := bufio.NewScanner(s.c)

	for scanner.Scan() {
		str := scanner.Text()

		if strings.Contains(str, "\"return\"") {
			return MessageTypeSuccess, []byte(str), nil
		} else if strings.Contains(str, "\"error\"") {
			return MessageTypeError, []byte(str), nil
		} else if strings.Contains(str, "\"event\"") {
			var e AsyncMessage

			err := json.Unmarshal([]byte(str), &e)
			if err != nil {
				return 0, nil, fmt.Errorf("qmp: json: failed to decode event: %s", err)
			}

			s.AsyncMessages = append(s.AsyncMessages, e)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, nil, fmt.Errorf("qmp: failed to read line: %s", err)
	}

	return 0, nil, fmt.Errorf("qmp: invalid response")
}
