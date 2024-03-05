// Copyright 2024 FishGoddess. All rights reserved.
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
