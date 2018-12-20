package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("jane", "abc123")
}

type ctxKey int

const (
	ctxUserId = iota
	ctxAuthToken
)

func UserId(ctx context.Context) string {
	return ctx.Value(ctxUserId).(string)
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}

func ProcessRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserId, userId)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf("handling response for %v (%v)\n", UserId(ctx), AuthToken(ctx))
}
