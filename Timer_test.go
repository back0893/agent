package main

import (
	"agent/src"
	"context"
	"testing"
)

func TestTimer(t *testing.T) {
	ctx := context.Background()
	timer := src.NewTimer(ctx)
}
