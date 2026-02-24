# MiniRedis (Go)

A lightweight Redis-inspired in-memory key-value database built in Go.

Supports concurrent clients over TCP, TTL expiration, LRU eviction, append-only persistence (AOF), and runtime statistics.

---

## Features

- Concurrent TCP server (goroutine per connection)
- SET / GET / DEL commands
- TTL expiration: `SET key value EX <seconds>`
- LRU eviction with fixed max capacity
- Append-only persistence (AOF) with replay on startup
- STATS command for runtime introspection
- Graceful shutdown handling
- Benchmark tool for throughput testing
- Unit tests for core store logic

---

## Architecture Overview

Client  
→ TCP Server  
→ Command Parser  
→ Store Engine  

### Store Engine

- In-memory hashmap for values
- TTL map for expirations with background cleanup goroutine
- LRU cache using hashmap + doubly linked list

### Persistence

- Append-only file (AOF)
- Replay on startup to rebuild in-memory state

---

## Quick Start

### Run the server

```bash
go run ./cmd/server

Server listens on:

localhost:6379
Example Client Session (PowerShell)
$client = New-Object System.Net.Sockets.TcpClient("localhost", 6379)
$stream = $client.GetStream()
$writer = New-Object System.IO.StreamWriter($stream)
$reader = New-Object System.IO.StreamReader($stream)
$writer.AutoFlush = $true

$reader.ReadLine() | Out-Null

$writer.WriteLine("SET name Chidera")
$reader.ReadLine()

$writer.WriteLine("GET name")
$reader.ReadLine()

$writer.WriteLine("SET temp hello EX 5")
$reader.ReadLine()

$writer.WriteLine("STATS")
$reader.ReadLine()

$writer.WriteLine("QUIT")
$reader.ReadLine()

$client.Close()

Supported Commands
| Command                      | Description                |
| ---------------------------- | -------------------------- |
| `PING`                       | Returns `PONG`             |
| `SET key value`              | Stores value               |
| `SET key value EX <seconds>` | Stores value with TTL      |
| `GET key`                    | Returns value or `NULL`    |
| `DEL key`                    | Deletes key                |
| `STATS`                      | Returns runtime statistics |
| `QUIT`                       | Closes connection          |

Benchmark

Run server in one terminal:

go run ./cmd/server

Run benchmark in another terminal:

go run ./cmd/bench
Sample Result (Local Machine)

Throughput: ~266,795 ops/sec

Workers: 50

Total operations: 20,000 (SET + GET)

Tests
go test ./...

Note: Go race detector requires CGO toolchain on Windows.

Skills Demonstrated

Concurrent systems design

TCP networking

Data structures (hashmap + doubly linked list)

Memory management strategies (TTL + LRU)

Persistence design (append-only log + replay)

Benchmarking and performance measurement

Clean project structure and modular design