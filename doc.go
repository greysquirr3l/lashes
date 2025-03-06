// Package lashes provides a proxy rotation library for Go applications.
//
// It supports multiple proxy types (HTTP, SOCKS4, SOCKS5), configurable rotation
// strategies, and optional database persistence. The library is designed to be
// zero-dependency for core functionality with optional database support.
//
// Basic usage:
//
//	rotator, err := lashes.New(lashes.DefaultOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	err = rotator.AddProxy(ctx, "http://proxy1.example.com:8080", lashes.HTTP)
//	client, err := rotator.Client(context.Background())
//	resp, err := client.Get("https://api.example.com")
//
// With database storage:
//
//	opts := lashes.Options{
//	    Storage: &lashes.StorageOptions{
//	        Type: lashes.SQLite,
//	        Database: lashes.StorageOptions{
//	            DSN: "file:proxies.db",
//	        },
//	    },
//	}
//	rotator, err := lashes.New(opts)
package lashes
