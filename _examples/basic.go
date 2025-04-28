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
