// Copyright 2024 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
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
	})

	if err := http.ListenAndServe("127.0.0.1:9876", nil); err != nil {
		panic(err)
	}
}

// acquireClient acquires a new http client, and returns an error if failed.
func acquireClient() (*http.Client, error) {
	fmt.Println("acquire client...")
	return &http.Client{}, nil
}

// releaseClient releases the given client, and returns an error if failed.
func releaseClient(client *http.Client) error {
	fmt.Println("release client...")
	return nil
}

func main() {
	// Prepare some backend resources.
	ctx := context.Background()
	go runServer()
	time.Sleep(time.Second)

	// Create a resource pool which type is *http.Client.
	// You should prepare two functions: acquire and release.
	// The acquire function is for acquiring a new resource, and you can do some setups for your resource.
	// The release function is for releasing the given resource, and you can destroy everything of your resource.
	// Also, you can specify some options to change the default settings.
	pool := rego.New[*http.Client](acquireClient, releaseClient, rego.WithLimit(8))
	defer pool.Close()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()

			// Take a client from the pool.
			// The pool will maintain the count of clients.
			client, err := pool.Take(ctx)
			if err != nil {
				panic(err)
			}

			// Remember put the client to pool when your using is done.
			// This is why we call the resource in pool is reusable.
			// We recommend you to do this job in a defer function.
			defer pool.Put(client)

			// Use the client whatever you want.
			body := strings.NewReader(strconv.Itoa(ii))
			_, err = client.Post("http://127.0.0.1:9876", "", body)
			if err != nil {
				panic(err)
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(time.Second)
}
