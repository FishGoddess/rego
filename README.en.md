# 🍦 Rego

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/rego)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](_icons/coverage.svg)
![Test](https://github.com/FishGoddess/rego/actions/workflows/test.yml/badge.svg)

**Rego** is for reusing some resources like network connections.

[阅读中文版的文档](./README.md)

### 🍭 Features

* Reuse resources by limiting the quantity
* Error handling supports
* Pool statistics supports
* Context timeout supports

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

// releaseConn releases the conn, and returns an error if failed.
func releaseConn(ctx context.Context, conn net.Conn) error {
	return conn.Close()
}

func main() {
	// Create a pool which type is net.Conn and limit is 64.
	ctx := context.Background()

	pool := rego.New(64, acquireConn, releaseConn)
	defer pool.Close(ctx)

	// Acquire a conn from pool.
	conn, err := pool.Acquire(ctx)
	if err != nil {
		panic(err)
	}

    // Remember releasing the conn after using.
	defer pool.Release(ctx, conn)

	// Use the conn
	fmt.Println(conn.RemoteAddr())
}
```

### 👥 Contributing

If you find that something is not working as expected please open an _**issue**_.
