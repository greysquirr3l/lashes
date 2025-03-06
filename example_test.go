package lashes_test

import (
	"context"
	"fmt"
	"log"

	"github.com/greysquirr3l/lashes"
)

func ExampleNew() {
	rotator, err := lashes.New(lashes.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = rotator.AddProxy(ctx, "http://example.com:8080", lashes.HTTP)
	if err != nil {
		log.Fatal(err)
	}

	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using proxy: %s\n", proxy.URL)
	// Output: Using proxy: http://example.com:8080
}
