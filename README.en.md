# 🍦 Rego

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/rego)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](_icons/coverage.svg)
![Test](https://github.com/FishGoddess/rego/actions/workflows/test.yml/badge.svg)

**Rego** is a resource pool used for reusing some resources like network connections.

[阅读中文版的文档](./README.md)

### 🍭 Features

* Reuse resources by using tokens to limit the quantity of resources
* Error handling callback for different errors
* Check pool statistics like active and idle quantity of resources
* Passing context to callbacks and supporting timeout.

_Check [HISTORY.md](./HISTORY.md) and [FUTURE.md](./FUTURE.md) to know about more information._

### 💡 How to use

```shell
$ go get -u github.com/FishGoddess/rego
```

```go
package main

import (
	"context"
	"fmt"
	"net"

	"github.com/FishGoddess/rego"
)

// acquireConn acquires a new conn, and returns an error if failed.
func acquireConn(ctx context.Context) (net.Conn, error) {
	// Guess this ip is from which websites?
	var dialer net.Dialer
	return dialer.DialContext(ctx, "tcp", "20.205.243.166:80")
}

// releaseConn releases the given conn, and returns an error if failed.
func releaseConn(ctx context.Context, conn net.Conn) error {
	return conn.Close()
}

func main() {
	// Create a resource pool which type is net.Conn and limit is 64.
	ctx := context.Background()

	pool := rego.New(64, acquireConn, releaseConn)
	defer pool.Close(ctx)

	// Take a resource from pool.
	conn, err := pool.Take(ctx)
	if err != nil {
		panic(err)
	}

	// Remember put the client to pool when your using is done.
	// This is why we call the resource in pool is reusable.
	// We recommend you to do this job in a defer function.
	defer pool.Put(ctx, conn)

	// Use the conn
	fmt.Println(conn.RemoteAddr())
}
```

### 🚄 Benchmarks

```shell
$ make bench
```

```shell
goos: linux
goarch: amd64
cpu: AMD EPYC 7K62 48-Core Processor

BenchmarkPool-2          3024594               378.2 ns/op            48 B/op          1 allocs/op
```

> Benchmarks: _examples/performance_test.go

### 👥 Contributing

If you find that something is not working as expected please open an _**issue**_.
