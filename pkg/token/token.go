// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import "context"

type Token struct{}

type TokenBucket struct {
	tokens chan Token
}

func NewBucket(limit uint64) *TokenBucket {
	bucket := &TokenBucket{
		tokens: make(chan Token, limit),
	}

	for range limit {
		bucket.tokens <- Token{}
	}

	return bucket
}

// ConsumeToken consumes a token from bucket and waits util context done if there is no token.
func (tb *TokenBucket) ConsumeToken(ctx context.Context) error {
	select {
	case <-tb.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ProduceToken produces a token to bucket.
func (tb *TokenBucket) ProduceToken() {
	select {
	case tb.tokens <- Token{}:
		return
	default:
		return
	}
}

// Close closes the bucket.
func (tb *TokenBucket) Close() error {
	close(tb.tokens)
	return nil
}
