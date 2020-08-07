package future

import (
	"fmt"
	"reflect"
	"time"
)

/*
    Modifier and Type      	Method and Description

    bool                 	Cancel(boolean mayInterruptIfRunning)
                         		Attempts to cancel execution of this task.

    []interface{}, error    Get(boolean mayInterruptIfRunning)
    							Waits if necessary for the computation to complete, and then retrieves its result.

	[]interface{}, error    GetWithTimeout(timeout time.Duration)
  								Waits if necessary for at most the given time for the computation to complete, and then retrieves its result, if available.

	bool 					IsCancelled()
								Returns true if this task was cancelled before it completed normally.

	bool 					IsDone()
								Returns true if this task is completed.
*/
// Future type holds Result and state
type Future struct {
	Start         bool               // True if Routine started
	Success       bool               // True if Routine finished successfully
	Cancelled     bool               // True if Routine Cancelled
	Err           error              // Holds CancellationError or any other Exception
	Result        []interface{}      // Return Value of the goroutine, User have to manually Type cast the variable, Advice to use reflect
	ErrorChannel  chan error         // Transmit CancellationError or any other Exception
	ResultChannel chan []interface{} // Transmit Routine return values
	KillChannel   chan bool          // Transmit cancel call
}

//Error when routine is cancelled
type CancellationError struct {
	Message   string
	ErrorTime time.Time
}

//Timeout Error when Result is not getting within given time
type TimeoutError struct {
	Message   string
	ErrorTime time.Time
}

//CancellationError formatting
func (e *CancellationError) Error() string {
	return fmt.Sprintf("at %v, %s", e.ErrorTime, e.Message)
}

//TimeoutError formatting
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("at %v, %s", e.ErrorTime, e.Message)
}

/*
* Attempts to cancel execution of this task. This attempt will fail if the task has already completed, has already been cancelled.
* If successful, and this task has not started when cancel is called, this task should never run.
* If the task has already started, then the mayInterruptIfRunning parameter determines whether the thread executing this task should be interrupted in an attempt to stop the task.
 */
func (f *Future) Cancel(mayInterruptIfRunning bool) bool {
	if f.Cancelled || f.Success || f.Err != nil {
		return false
	} else if f.Start {
		if mayInterruptIfRunning {
			f.KillChannel <- true
			f.Cancelled = true
			return true
		} else {
			return false
		}
	} else {
		f.KillChannel <- true
		f.Cancelled = true
		return true
	}
	return false
}

/*
* Returns true if this task was cancelled before it completed normally
 */
func (f *Future) IsCancelled() bool {
	return f.Cancelled
}

/*
* Returns true if this task completed. Completion may be due to normal termination, an exception, or cancellation -- in all of these cases, this method will return true.
 */
func (f *Future) IsDone() bool {
	return f.Cancelled || f.Success || f.Err != nil
}

/*
* Waits if necessary for the computation to complete, and then retrieves its result.
 */
func (f *Future) Get() ([]interface{}, error) {

	if f.Success && f.Cancelled == false {
		return f.Result, nil
	}

	if f.Cancelled {
		return nil, &CancellationError{
			"task is cancelled by user",
			time.Now(),
		}
	}
	if f.Err != nil {
		return nil, f.Err
	}

	select {
	case res := <-f.ResultChannel:
		f.Result = res
		f.Success = true
		return f.Result, nil
	case f.Err = <-f.ErrorChannel:
		f.Result = nil
		f.Success = false
		return nil, f.Err
	}
	return f.Result, f.Err
}

/*
* Waits if necessary for at most the given time for the computation to complete, and then retrieves its result, if available.
 */
func (f *Future) GetWithTimeout(timeout time.Duration) ([]interface{}, error) {

	if f.Success && f.Cancelled == false {
		return f.Result, nil
	}

	if f.Cancelled {
		return nil, &CancellationError{
			"task is cancelled by user",
			time.Now(),
		}
	}
	if f.Err != nil {
		return nil, f.Err
	}

	timeoutChannel := time.After(timeout)
	select {
	case res := <-f.ResultChannel:
		f.Result = res
		f.Success = true
		return f.Result, nil
	case f.Err = <-f.ErrorChannel:
		f.Result = nil
		f.Success = false
		return nil, f.Err
	case <-timeoutChannel:
		f.Result = nil
		f.Success = false
		return nil, &TimeoutError{
			"Time Out, result not found",
			time.Now(),
		}
	}
	return f.Result, f.Err
}

/*
*	Wrapper to start a function in a goroutine and get the result from a channel
*	First parameter is function and second is interfaces for arguments
 */
func StartGo(functionImplementation interface{}, args ...interface{}) *Future {
	arguments := make([]reflect.Value, len(args), len(args))
	function := reflect.ValueOf(functionImplementation)
	for i, v := range args {
		arguments[i] = reflect.ValueOf(v)
	}
	resultChannel := make(chan []interface{}) // result will be pushed to this channel
	errorChannel := make(chan error)          // send error
	quit := make(chan bool)                   // receive cancel call
	future := &Future{
		Start:         false,
		Success:       false,
		Cancelled:     false,
		Err:           nil,
		Result:        nil,
		ErrorChannel:  errorChannel,
		ResultChannel: resultChannel,
		KillChannel:   quit,
	}
	var convertTOInterface []interface{}
	go func() {
		for {
			select {
			case <-quit:
				errorChannel <- &CancellationError{"task is cancelled by user", time.Now()}
				return
			default:
				future.Start = true
				res := function.Call(arguments)
				for i := 0; i < len(res); i++ {
					convertTOInterface = append(convertTOInterface, res[i].Interface())
				}
				future.Result = convertTOInterface
				future.Success = true
				fmt.Println(convertTOInterface)
				resultChannel <- convertTOInterface
				return
			}
		}
	}()
	return future
}
