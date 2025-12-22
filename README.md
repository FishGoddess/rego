# 🍦 Rego

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/rego)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](_icons/coverage.svg)
![Test](https://github.com/FishGoddess/rego/actions/workflows/test.yml/badge.svg)

**Rego** 用于复用一些特定的资源，比如说网络连接。

[Read me in English](./README.en.md)

### 🍭 功能特性

* 复用资源，支持限制数量
* 支持自定义错误
* 支持查询运行指标
* 支持 context 超时机制

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

### 👥 贡献者

如果您觉得 rego 缺少您需要的功能，请不要犹豫，马上参与进来，发起一个 _**issue**_。
