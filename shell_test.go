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
	killChan, err := Exec0("emulator -avd Nexus_5_API_25")
	if err != nil {
		fmt.Println("start command failed")
	} else {
		time.Sleep(10*time.Second)
		fmt.Println("the emulator will be killed after 10s")
		killChan<-0
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

//func TestExec0_Normal(t *testing.T) {
//	killChan, err := Exec0("emulator -avd Nexus_5_API_25")
//	if err != nil {
//		t.Fatal("start command failed")
//	} else {
//		time.Sleep(10*time.Second)
//		t.Log("the emulator will be killed after 10s")
//		killChan<-0
//	}
//}

func TestExec0_Error(t *testing.T) {
	killChan, err := Exec0("top")
	if err != nil {
		t.Fatal("top must fail")
	} else {
		t.Log("kill top")
		killChan<-0
	}
}
