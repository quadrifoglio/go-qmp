# Go-QMP - Golang interface to QMP

The QEMU Machine Protocol (QMP) is a JSON-based protocol which allows applications to control a QEMU instance.

This library offers a simple interface to QMP for the Go programming language.

## Basic Usage

```go
import "github.com/quadrifoglio/go-qmp"

// Connection to QMP
c, err := qmp.Open("unix", "/tmp/qmp.sock")
if err != nil {
	log.Fatal(err)
}

defer c.Close()

// Execute simple QMP command
result, err = c.Command("query-status", nil)
if err != nil {
	log.Fatal(err)
}

fmt.Println(result)

// Execute QMP command with arguments

args := map[string]string {
	"device": "ide1-cd0"
}

result, err = c.Command("eject", args)
if err != nil {
	log.Fatal(err)
}

// Execute HMP command
result, err = c.HumanMonitorCommand("savevm checkpoint")
if err != nil {
	log.Fatal(err)
}
```
