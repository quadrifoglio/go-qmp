# Go-QMP - Golang interface to QMP

## Usage

```go
import "github.com/quadrifoglio/go-qmp"

// Connection to QMP
c, err := qmp.Open("unix", "/tmp/qmp.sock")
if err != nil {
	return err
}

defer c.Close() 

// Execute simple QMP command
result, err = c.Command("query-status", nil)
if err != nil {
	return err
}

fmt.Println(result)

// Execute QMP command with arguments

args := map[string]string {
	"device": "ide1-cd0"
}

result, err = c.Command("eject", args)
if err != nil {
	return err
}

// Execute HMP command
result, err = c.HumanMonitorCommand("savevm checkpoint")
if err != nil {
	return err
}
```
