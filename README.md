# nanoshlib

[![Go Reference](https://pkg.go.dev/badge/github.com/qaqcatz/nanoshlib.svg)](https://pkg.go.dev/github.com/qaqcatz/nanoshlib)

Simple and stable shell interface.

Support timeout and service management.

# How to use

## import

```golang
// go.mod:
require github.com/qaqcatz/nanoshlib v1.5.0
// xxx.go:
import "github.com/qaqcatz/nanoshlib"
```

## Exec

```golang
func Exec(cmdStr string, timeoutMS int) (string, string, error)
```

Exec use /bin/bash -c to execute cmdStr, wait for the result, or timeout, return out stream, error stream, and an error, which can be nil, normal error or *TimeoutError.
timeoutMS <= 0 means timeoutMS = inf

**Example**

```golang
func TestExecNormal(t *testing.T) {
	outStream, errStream, err := Exec("echo helloworld", 1000)
	if err != nil {
		t.Fatal("[err]\n" + err.Error())
	} else {
		t.Log("[out stream]\n" + outStream)
		t.Log("[err stream]\n" + errStream)
	}
}
```

## Exec0

```golang
func Exec0(myCmdStr string) (chan error, chan int, error)
```

Exec0 is an extension of Exec. When Exec0 is executed, it will return immediately. Exec0 does not care about the result, it is only responsible for making sure the process can be started/killed successfully and checking if the process is still running. It can be used to start/monitor/kill a service.

Specifically, Exec0 will return doneChan, killChan, err:

- doneChan, you can use select case: <-errChan default: ... to check if the process is still running
- killChan, you can use killChan<-0(any number) to kill the process.
- err, command start error.

**Example**

```golang
// you should install an avd named 'test' first.
func TestExec0_Normal(t *testing.T) {
	doneChan, killChan, err := Exec0("/home/hzy/Android/Sdk/emulator/emulator -avd test")
	if err != nil {
		t.Fatal("start command failed")
	} else {
		t.Log("wait 10s for start")
		time.Sleep(10*time.Second)
		select {
		case err := <-doneChan:
			errStr := "no err"
			if err != nil {
				errStr = err.Error()
			}
			t.Fatal("the emulator is closed! " + errStr)
		default:
			t.Log("the emulator is running")
		}
		t.Log("kill the emulator")
		killChan<-0
		t.Log("wait 3s for kill")
		time.Sleep(3*time.Second)
		select {
		case err := <-doneChan:
			errStr := "no err"
			if err != nil {
				errStr = err.Error()
			}
			t.Log("the emulator is closed. " + errStr)
		default:
			t.Fatal("the emulator is still running!")
		}
	}
}
```

## Exec0s

```golang
func Exec0s(myCmdStr string) (chan error, chan int, error)
```

Exec0s is an extension of Exec0. We use cmd.SysProcAttr = &syscall.SysProcAttr{Setsid:true} to avoid the process being killed when the program ends.

```golang
// you should install an avd named 'test' first.
func TestExec0s_createSession(t *testing.T) {
	_, _, err := Exec0s("/home/hzy/Android/Sdk/emulator/emulator -avd test")
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(3*time.Second)
}
```
## ExecStd

Likes Exec, but write the result to standard output stream and standard error stream.

## ExecStdX

likes ExecStdX, but it also returns the outStream and errStream.