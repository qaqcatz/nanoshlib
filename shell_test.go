package nanoshlib

import (
	"fmt"
	"time"
)
import "testing"

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

func TestExecNormal(t *testing.T) {
	outStream, errStream, err := Exec("echo helloworld", 1000)
	if err != nil {
		t.Fatal("[err]\n" + err.Error())
	} else {
		t.Log("[out stream]\n" + string(outStream))
		t.Log("[err stream]\n" + string(errStream))
	}
}

func TestExecTimeout(t *testing.T) {
	_, _, err := Exec("sleep 3s", 1000)
	if err != nil {
		switch err.(type) {
		case *TimeoutError:
			t.Log(err.Error())
		default:
			t.Fatal("sleep 3s must timeout")
		}
	} else {
		t.Fatal("sleep 3s must fail")
	}
}

func TestExecError(t *testing.T) {
	outStream, errStream, err := Exec("top", 1000)
	if err != nil {
		switch err.(type) {
		case *TimeoutError:
			t.Fatal("top must fail but not timeout")
		default:
			t.Log("[out stream]\n" + string(outStream))
			t.Log("[err stream]\n" + string(errStream))
		}
	} else {
		t.Fatal("top must fail")
	}
}

func TestExecTimeout0(t *testing.T) {
	_, _, err := Exec("sleep 3s", 0)
	if err != nil {
		t.Fatal("sleep 3s must succeed")
	} else {
		t.Log("succeed")
	}
}

// you should install an avd named 'test' first.
//func TestExec0_Normal(t *testing.T) {
//	doneChan, killChan, err := Exec0("/home/hzy/Android/Sdk/emulator/emulator -avd test")
//	if err != nil {
//		t.Fatal("start command failed")
//	} else {
//		t.Log("wait 10s for start")
//		time.Sleep(10*time.Second)
//		select {
//		case err := <-doneChan:
//			errStr := "no err"
//			if err != nil {
//				errStr = err.Error()
//			}
//			t.Fatal("the emulator is closed! " + errStr)
//		default:
//			t.Log("the emulator is running")
//		}
//		t.Log("kill the emulator")
//		killChan<-0
//		t.Log("wait 3s for kill")
//		time.Sleep(3*time.Second)
//		select {
//		case err := <-doneChan:
//			errStr := "no err"
//			if err != nil {
//				errStr = err.Error()
//			}
//			t.Log("the emulator is closed. " + errStr)
//		default:
//			t.Fatal("the emulator is still running!")
//		}
//	}
//}

func TestExec0_Error(t *testing.T) {
	_, killChan, err := Exec0("top")
	if err != nil {
		t.Fatal("top must fail")
	} else {
		t.Log("kill top")
		killChan<-0
	}
}
