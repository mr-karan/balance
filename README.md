<p align="center">
<img src="./.github/logo.png" alt="logo" width="15%" />
</p>

# balance

A minimal Golang library for implementing weighted round-robin load balancing for a given set of items.

## Installation

```bash
go get github.com/mr-karan/balance
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/mr-karan/balance"
    
    // Create a new load balancer.
    b := balance.NewBalance()

    // Add items to the load balancer with their corresponding weights.
    b.Add("a", 5)
    b.Add("b", 3)
    b.Add("c", 2)    
    
    // Get the next item from the load balancer.
    fmt.Println(b.Next())

    // For 10 requests, the output sequence will be: [a b c a a b a c b a]
)
```

## Algorithm

The algorithm is based on the [Smooth Weighted Round Robin](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35) used by NGINX.

> Algorithm is as follows: on each peer selection we increase current_weight
of each eligible peer by its weight, select peer with greatest current_weight
and reduce its current_weight by total number of weight points distributed
among peers.

## Examples

### Round Robin

For implementing an equal weighted round-robin load balancer for a set of servers, you can use the following config:

```go
b.Add("server1", 1)
b.Add("server2", 1)
b.Add("server3", 1)
```

Since the weights of all 3 servers are equal, the load balancer will distribute the load equally among all 3 servers.

### Weighted Round Robin

For implementing a weighted round-robin load balancer for a set of servers, you can use the following config:

```go
b.Add("server1", 5)
b.Add("server2", 3)
b.Add("server3", 2)
```

The load balancer will distribute the load in the ratio of 5:3:2 among the 3 servers.

## Benchmark

```bash
go test -v -failfast -bench=. -benchmem -run=^$
goos: linux
goarch: amd64
pkg: github.com/mr-karan/balance
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkBalance
BenchmarkBalance/items-10
BenchmarkBalance/items-10-8             18100609                65.84 ns/op            0 B/op          0 allocs/op
BenchmarkBalance/items-100
BenchmarkBalance/items-100-8             7224061               164.3 ns/op             0 B/op          0 allocs/op
BenchmarkBalance/items-1000
BenchmarkBalance/items-1000-8            1000000              1089 ns/op               0 B/op          0 allocs/op
BenchmarkBalance/items-10000
BenchmarkBalance/items-10000-8             88384             13776 ns/op               0 B/op          0 allocs/op
BenchmarkBalance/items-100000
BenchmarkBalance/items-100000-8             7112            166452 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/mr-karan/balance     6.367s
```
