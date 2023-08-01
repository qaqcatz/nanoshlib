package nanoshlib

import (
	"testing"
	"time"
)

func TestExecNormal(t *testing.T) {
	outStream, errStream, err := Exec("echo helloworld", 1000)
	if err != nil {
		t.Fatal("[err]\n" + err.Error())
	} else {
		t.Log("[out stream]\n" + outStream)
		t.Log("[err stream]\n" + errStream)
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
			t.Log("[out stream]\n" + outStream)
			t.Log("[err stream]\n" + errStream)
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

func TestExec0_Error(t *testing.T) {
	_, killChan, err := Exec0("top")
	if err != nil {
		t.Fatal("top must succeed")
	} else {
		t.Log("kill top")
		killChan<-0
	}
}

// you should install an avd named 'test' first.
func TestExec0s_createSession(t *testing.T) {
	_, _, err := Exec0s("/home/hzy/Android/Sdk/emulator/emulator -avd test")
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(3*time.Second)
}

func TestExecStd(t *testing.T) {
	err := ExecStd("echo hello && sleep 1s && echo hello && sleep 1s && echo hello", 5000)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecStdX(t *testing.T) {
	outStream, errStream, err := ExecStdX("echo hello && sleep 1s && echo hello && sleep 1s && echo hello", 5000)
	if err != nil {
		t.Fatal(err, errStream)
	}
	t.Log("out:", outStream)
	t.Log("err:", errStream)
}

func TestExecStdX2(t *testing.T) {
	outStream, errStream, err := ExecStdX("abcxyz", 5000)
	if err == nil {
		t.Fatal("expect an error")
	} else {
		t.Log("out:", outStream)
		t.Log("err:", errStream)
	}
}