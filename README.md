# 🍦 Rego

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/rego)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](_icons/coverage.svg)
![Test](https://github.com/FishGoddess/rego/actions/workflows/test.yml/badge.svg)

**Rego** 是一个简单的资源池库，用于控制、复用一些特定的资源，比如说网络连接。

[Read me in English](./README.en.md)

### 🍭 功能特性

* 简单的资源池，可以控制资源数量

_历史版本的特性请查看 [HISTORY.md](./HISTORY.md)。未来版本的新特性和计划请查看 [FUTURE.md](./FUTURE.md)。_

### 💡 使用方式

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
	pool := rego.New(64, acquireConn, releaseConn)
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

### 👥 贡献者

如果您觉得 rego 缺少您需要的功能，请不要犹豫，马上参与进来，发起一个 _**issue**_。
