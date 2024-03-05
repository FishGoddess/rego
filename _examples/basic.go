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
	pool := rego.New[net.Conn](acquireConn, releaseConn, rego.WithLimit(64))
	defer pool.Close()

	conn, err := pool.Take(context.Background())
	if err != nil {
		panic(err)
	}

	defer pool.Put(conn)

	// Use the conn
	fmt.Println(conn.RemoteAddr())
}
