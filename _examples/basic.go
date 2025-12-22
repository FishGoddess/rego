// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
