# nanoshlib
Simple and stable shell interface.

Support timeout and service management.

# How to use

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
func Exec0(myCmdStr string) (chan int, error)
```

Exec0 is an extension of Exec. When Exec0 is executed, it will return immediately. Exec0 does not care about the result, or even if the process is still running, it is only responsible for making sure the process can be started/killed successfully. It can be used to start/kill a service.
Specifically, Exec0 will return killCahn, err:

- killChan, you can use killChan<-0(any number) to kill the process.
- err, command start error.

**Example**

```golang
func ExampleExec0_normal() {
	killChan, err := Exec0("emulator -avd Nexus_5_API_25")
	if err != nil {
		fmt.Println("start command failed")
	} else {
		time.Sleep(10*time.Second)
		fmt.Println("the emulator will be killed after 10s")
		killChan<-0
	}
}
```

