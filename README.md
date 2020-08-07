# future-Golang-Library

This is concurrency management library similar to java Futures Interface and Scala Future, it handle the result of an asynchronous computation. Methods are provided to check if the computation is complete, to wait for its completion, and to retrieve the result of the computation. The result can only be retrieved using method get when the computation has completed, blocking if necessary until it is ready. Cancellation is performed by the cancel method. Additional methods are provided to determine if the task completed normally or was cancelled. Once a computation has completed, the computation cannot be cancelled. If you would like to use a Future for the sake of cancellability but not provide a usable result.

Although there are many ways to handle this behaviour in Golang.
This library is useful for people who got used to Java/Scala Future implementation.


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
// do something else here
// get result when needed
result := f.Get()
```

Also it is possible to use timeouts on Get
```golang
result := future.GetWithTimeout(3 * time.Second)
```

Java Futures: https://docs.oracle.com/javase/8/docs/api/index.html?java/util/concurrent/Future.html

Scala Futures: https://docs.scala-lang.org/overviews/core/futures.html
