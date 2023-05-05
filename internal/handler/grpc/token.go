package grpc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"google.golang.org/grpc/metadata"
	"time"
)

func createToken(key []byte) (token string) {
	h := hmac.New(sha256.New, key)
	src := []byte(fmt.Sprint(time.Now().UnixNano()))
	h.Write(src)

	return hex.EncodeToString(h.Sum(nil)) + "-" + hex.EncodeToString(src)
}

func getOrCreateToken(ctx context.Context, key []byte) (token string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return createToken(key), false
	}

	values := md.Get("token")
	if len(values) == 0 {
		return createToken(key), false
	}

	return values[0], true
}
