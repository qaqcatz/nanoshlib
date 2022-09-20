# nanoshlib

[![Go Reference](https://pkg.go.dev/badge/github.com/qaqcatz/nanoshlib.svg)](https://pkg.go.dev/github.com/qaqcatz/nanoshlib)

Simple and stable shell interface.

Support timeout and service management.

# How to use

## import

```golang
// go.mod:
require github.com/qaqcatz/nanoshlib v1.0.0
// xxx.go:
import "github.com/qaqcatz/nanoshlib"
```

## Exec

```golang
func Exec(cmdStr string, timeoutMS int) ([]byte, []byte, error)
```

Exec use /bin/bash -c to execute cmdStr, wait for the result, or timeout, return out stream, error stream, and an error, which can be nil, normal error or *TimeoutError.
timeoutMS <= 0 means timeoutMS = inf

**Example**

```golang
func ExampleExec_normal() {
	outStream, errStream, err := Exec("echo helloworld", 1000)
	if err != nil {
		fmt.Println("[err]\n" + err.Error())
	} else {
		fmt.Println("[out stream]\n" + string(outStream))
		fmt.Println("[err stream]\n" + string(errStream))
	}
	// Output:
	// [out stream]
	// helloworld
	//
	// [err stream]
}

func ExampleExec_timeout() {
	_, _, err := Exec("sleep 3s", 1000)
	if err != nil {
		switch err.(type) {
		case *TimeoutError:
			fmt.Println(err.Error())
		default:
			fmt.Println("sleep 3s must timeout")
		}
	} else {
		fmt.Println("sleep 3s must fail")
	}
	// Output:
	// time out error
}

func ExampleExec_error() {
	_, _, err := Exec("sleep 3s", 0)
	if err != nil {
		fmt.Println("sleep 3s must succeed")
	} else {
		fmt.Println("succeed")
	}
	// Output:
	// succeed
}
```

## Exec0

```golang
func Exec0(myCmdStr string, createSession bool) (chan error, chan int, error)
```

Exec0 is an extension of Exec. When Exec0 is executed, it will return immediately. Exec0 does not care about the result, it is only responsible for making sure the process can be started/killed successfully and checking if the process is still running. It can be used to start/monitor/kill a service.

You can set createSession to true to avoid the process being killed when the program ends.

Specifically, Exec0 will return doneChan, killChan, err:

- doneChan, you can use select case: <-errChan default: ... to check if the process is still running
- killChan, you can use killChan<-0(any number) to kill the process.
- err, command start error.

**Example**

```golang
func ExampleExec0_normal() {
	doneChan, killChan, err := Exec0("/home/hzy/Android/Sdk/emulator/emulator -avd test")
	if err != nil {
		fmt.Println("start command failed")
	} else {
		fmt.Println("wait 10s for start")
		time.Sleep(10*time.Second)
		select {
		case err := <-doneChan:
			errStr := "no err"
			if err != nil {
				errStr = err.Error()
			}
			fmt.Println("the emulator is closed! " + errStr)
			return
		default:
			fmt.Println("the emulator is running")
		}
		fmt.Println("kill the emulator")
		killChan<-0
		fmt.Println("wait 3s for kill")
		time.Sleep(3*time.Second)
		select {
		case err := <-doneChan:
			errStr := "no err"
			if err != nil {
				errStr = err.Error()
			}
			fmt.Println("the emulator is closed. " + errStr)
		default:
			fmt.Println("the emulator is still running!")
		}
	}
}

func ExampleExec0_createSession() {
	_, _, err := Exec0("/home/hzy/Android/Sdk/emulator/emulator -avd test", true)
	if err != nil {
		log.Fatal(err.Error())
	}
	time.Sleep(3*time.Second)
	// The emulator will not be killed when the program ends.
	// If you change createSession:true to false, the emulator will be killed when the program ends.
}
```

