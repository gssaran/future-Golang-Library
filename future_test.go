package future_test

import (
	"future"
	"reflect"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func function() (string, int) {
	defer wg.Done()
	name := "gouri"
	id := 100
	time.Sleep(1 * time.Second)
	return name, id
}

func TestStartGo(t *testing.T) {
	var expectedResult []interface{}
	expectedResult = append(expectedResult, "gouri", 100)
	startTime := time.Now()
	wg.Add(1)
	f := future.StartGo(function)
	wg.Wait()

	timeDuration := time.Since(startTime)

	t.Logf("function took %s", timeDuration)

	if timeDuration.Milliseconds() < (1 * time.Second).Milliseconds() {
		t.Errorf("returned before expected time, either routine not started or terminated")
	}

	if reflect.DeepEqual(expectedResult, f.Result) == false {
		t.Errorf("Expecting %v, but get %v", expectedResult, f.Result)
	}

	if f.Start == false || f.Success == false {
		t.Errorf("routine not started or not successfully completed")
	}
}

func TestCancel(t *testing.T) {
	//test Cancel for completed ot cancelled task
	wg.Add(1)
	f := future.StartGo(function)
	f.Get() // waiting to finish the task

	if f.Success {
		t.Logf("Task is Finished at %v", time.Now())
	}

	if f.Cancel(false) == true {
		t.Errorf("Completed or Cancelled task can't be Cancelled ")
	}

	if f.Cancel(true) == true {
		t.Errorf("Completed or Cancelled task can't be Interrupted ")
	}

	// test Cancel, before Finished the task
	f = future.StartGo(func() (string, int) {
		name := "gouri"
		id := 100
		time.Sleep(1 * time.Second)
		return name, id
	})

	if f.Start == true && f.Cancel(false) == true {
		t.Errorf("Task can't be cancelled without Interrupted one it started")
	}

	if f.Cancelled == false && f.Cancel(true) == false {
		t.Errorf("Task is not cancelled")
	}

	if f.Cancelled == true && f.Success == true {
		t.Errorf("Task is cancelled but showing Successed")
	}
}

func TestIsDone(t *testing.T) {
	//If Completed
	wg.Add(1)
	f := future.StartGo(function)
	wg.Wait()
	t.Logf("Task is Finished at %v", time.Now())
	if f.IsDone() == false {
		t.Errorf("Task is finished but showing executing")
	}

	// If Cancelled
	f = future.StartGo(func() (string, int) {
		name := "gouri"
		id := 100
		time.Sleep(1 * time.Second)
		return name, id
	})
	if f.Cancel(true) {
		t.Logf("Task is Cancelled at %v", time.Now())
		if f.IsDone() == false {
			t.Errorf("Task is Cancelled but showing executing")
		}
	}
}

func TestGet(t *testing.T) {
	var expectedResult []interface{}
	expectedResult = append(expectedResult, "gouri", 100)
	//if task not cancelled
	startTime := time.Now()
	wg.Add(1)
	f := future.StartGo(function)
	result, _ := f.Get()
	timeDuration := time.Since(startTime)

	t.Logf("function took %s", timeDuration)

	if timeDuration.Milliseconds() < (1 * time.Second).Milliseconds() {
		t.Errorf("returned before expected time, either routine not started or terminated")
	}

	if reflect.DeepEqual(expectedResult, result) == false {
		t.Errorf("Expecting %v, but found %v", expectedResult, result)
	}
	//it task is cancelled
	wg.Add(1)
	f = future.StartGo(function)
	if f.Cancel(true) {
		_, err := f.Get()
		if err == nil {
			t.Errorf("Expecting %v, but found %v", future.CancellationError{"task is cancelled by user", time.Now()}, err)
		}
	}
}

func TestGetWithTimeout(t *testing.T) {
	var expectedResult []interface{}
	expectedResult = append(expectedResult, "gouri", 100)
	//if task not cancelled
	startTime := time.Now()
	wg.Add(1)
	f := future.StartGo(function)
	result, err := f.GetWithTimeout(100000000)
	if err == nil {
		t.Errorf("Expecting %v, but found %v", future.TimeoutError{"Time Out, result not found", time.Now()}, err)
	}

	result, err = f.GetWithTimeout(2000000000)
	timeDuration := time.Since(startTime)

	t.Logf("function took %s", timeDuration)

	if timeDuration.Milliseconds() < (1 * time.Second).Milliseconds() {
		t.Errorf("returned before expected time, either routine not started or terminated")
	}

	if err == nil && reflect.DeepEqual(expectedResult, result) == false {
		t.Errorf("Expecting %v, but get %v", expectedResult, result)
	}

	//it task is cancelled
	wg.Add(1)
	f = future.StartGo(function)
	if f.Cancel(true) {
		result, err = f.GetWithTimeout(2000000000)
		if err == nil {
			t.Errorf("Expecting %v, but get %v", future.CancellationError{"task is cancelled by user", time.Now()}, err)
		}
	}
}
