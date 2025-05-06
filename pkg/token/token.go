// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import "context"

type Token struct{}

type Bucket struct {
	tokens chan Token
}

func NewBucket(limit uint64) *Bucket {
	bucket := &Bucket{
		tokens: make(chan Token, limit),
	}

	for range limit {
		bucket.tokens <- Token{}
	}

	return bucket
}

// ConsumeToken consumes a token from bucket and waits util context done if there is no token.
func (b *Bucket) ConsumeToken(ctx context.Context) error {
	select {
	case <-b.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ProduceToken produces a token to bucket.
func (b *Bucket) ProduceToken() {
	select {
	case b.tokens <- Token{}:
		return
	default:
		return
	}
}

// Close closes the bucket.
func (b *Bucket) Close() error {
	close(b.tokens)
	return nil
}
