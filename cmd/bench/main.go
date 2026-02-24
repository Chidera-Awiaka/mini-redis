package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	addr := "localhost:6379"
	concurrency := 50
	opsPerWorker := 200

	var totalOps int64
	var failed int64

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for w := 0; w < concurrency; w++ {
		go func(worker int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", addr)
			if err != nil {
				atomic.AddInt64(&failed, 1)
				fmt.Println("Worker failed to connect:", err)
				return
			}
			defer conn.Close()

			reader := bufio.NewReader(conn)

			// Read welcome line
			_, _ = reader.ReadString('\n')

			for i := 0; i < opsPerWorker; i++ {
				key := fmt.Sprintf("k%d_%d", worker, i)

				_, _ = conn.Write([]byte(fmt.Sprintf("SET %s %d\n", key, i)))
				_, _ = reader.ReadString('\n')

				_, _ = conn.Write([]byte(fmt.Sprintf("GET %s\n", key)))
				_, _ = reader.ReadString('\n')

				atomic.AddInt64(&totalOps, 2)
			}
		}(w)
	}

	wg.Wait()

	elapsed := time.Since(start).Seconds()
	ops := atomic.LoadInt64(&totalOps)

	fmt.Printf("\nTotal ops: %d\n", ops)
	fmt.Printf("Failed workers: %d\n", atomic.LoadInt64(&failed))
	fmt.Printf("Elapsed: %.3fs\n", elapsed)

	if elapsed > 0 {
		fmt.Printf("Throughput: %.0f ops/sec\n", float64(ops)/elapsed)
	}
}
