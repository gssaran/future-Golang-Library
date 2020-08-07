# future-Golang-Library

This is concurrency management library similar to java Futures Interface and Scala Future, it handle the result of an asynchronous computation. Methods are provided to check if the computation is complete, to wait for its completion, and to retrieve the result of the computation. The result can only be retrieved using method get when the computation has completed, blocking if necessary until it is ready. Cancellation is performed by the cancel method. Additional methods are provided to determine if the task completed normally or was cancelled. Once a computation has completed, the computation cannot be cancelled. If you would like to use a Future for the sake of cancellability but not provide a usable result.

Although there are many ways to handle this behaviour in Golang.
This library is useful for people who got used to Java/Scala Future implementation.


### Method Description:
```golang
Modifier and Type      	Method and Description

bool                 	Cancel(boolean mayInterruptIfRunning)
                        	Attempts to cancel execution of this task.
[]interface{}, error    Get(boolean mayInterruptIfRunning)
    				Waits if necessary for the computation to complete, and then retrieves its result.

[]interface{}, error    GetWithTimeout(timeout time.Duration)
  				Waits if necessary for at most the given time for the computation to complete, and then retrieves its result, if available.

bool 			IsCancelled()
				Returns true if this task was cancelled before it completed normally.

bool 			IsDone()
				Returns true if this task is completed.
```
#### Import:
```golang
import gofuture "github.com/gssaran/future-Golang-Library"
```

#### Usage:

```golang
f := future.StartGo(func () (string, int) {
	name := "gouri"
	id := 100
	time.Sleep(1 * time.Second)
	return name, id
})
// do Something
result := f.Get()
```

Also it is possible to use timeouts on Get
```golang
result := future.GetWithTimeout(3 * time.Second)
```

Java Futures: https://docs.oracle.com/javase/8/docs/api/index.html?java/util/concurrent/Future.html

Scala Futures: https://docs.scala-lang.org/overviews/core/futures.html
