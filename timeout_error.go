package nanoshlib

type TimeoutError struct {

}

func (TimeoutError *TimeoutError) Error() string {
	return "time out error"
}
