// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/FishGoddess/rego"
)

// runServer runs a test server for printing some messages from client.
func runServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("server:", string(bs))
		time.Sleep(time.Second)
	})

	if err := http.ListenAndServe("127.0.0.1:9876", nil); err != nil {
		panic(err)
	}
}

// acquireClient acquires a new http client, and returns an error if failed.
func acquireClient(ctx context.Context) (*http.Client, error) {
	fmt.Println("acquire client...")
	return &http.Client{}, nil
}

// releaseClient releases the client, and returns an error if failed.
func releaseClient(ctx context.Context, client *http.Client) error {
	fmt.Println("release client...")
	return nil
}

func availableClient(ctx context.Context, client *http.Client) bool {
	fmt.Println("available client...")
	return true
}

func poolClosedErr(ctx context.Context) error {
	return errors.New("_example: http client pool is closed")
}

func main() {
	// Run a server for test.
	ctx := context.Background()

	go runServer()
	time.Sleep(time.Second)

	// Create a pool which type is *http.Client.
	pool := rego.New(4, acquireClient, releaseClient).WithAvailableFunc(availableClient).WithPoolClosedErrFunc(poolClosedErr)
	defer pool.Close(ctx)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()

			// Acquire a client from the pool.
			client, err := pool.Acquire(ctx)
			if err != nil {
				panic(err)
			}

			// Remember releasing the client after using.
			defer pool.Release(ctx, client)

			// Use the client whatever you want.
			body := strings.NewReader(strconv.Itoa(ii))

			_, err = client.Post("http://127.0.0.1:9876", "", body)
			if err != nil {
				panic(err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("pool status: %+v\n", pool.Status())

	pool.Close(ctx)
	_, err := pool.Acquire(ctx)
	fmt.Printf("pool acquire err: %+v\n", err)
}
