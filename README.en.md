# 🍦 Rego

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/rego)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](_icons/coverage.svg)
![Test](https://github.com/FishGoddess/rego/actions/workflows/test.yml/badge.svg)

**Rego** is a resource pool library which is used for controlling and reusing some resources like network connection.

[阅读中文版的文档](./README.md)

### 🍭 Features

* Based resource pool, which can limit the count of resources.

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
func acquireConn() (net.Conn, error) {
	// Guess this ip is from which websites?
	return net.Dial("tcp", "20.205.243.166:80")
}

// releaseConn releases the given conn, and returns an error if failed.
func releaseConn(conn net.Conn) error {
	return conn.Close()
}

func main() {
	// Create a resource pool which type is net.Conn and limit is 64.
	pool := rego.New[net.Conn](acquireConn, releaseConn, rego.WithLimit(64))
	defer pool.Close()

	// Take a resource from pool.
	conn, err := pool.Take(context.Background())
	if err != nil {
		panic(err)
	}

	// Remember put the client to pool when your using is done.
	// This is why we call the resource in pool is reusable.
	// We recommend you to do this job in a defer function.
	defer pool.Put(conn)

	// Use the conn
	fmt.Println(conn.RemoteAddr())
}
```

### 👥 Contributing

If you find that something is not working as expected please open an _**issue**_.
