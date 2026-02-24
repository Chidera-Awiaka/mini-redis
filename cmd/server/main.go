package main

import (
	"log"
	"mini-redis/internal/persistence"
	"mini-redis/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	aofPath := "data.aof"

	aof, err := persistence.Open(aofPath)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(":6379", aof)

	// Replay existing commands from disk
	err = persistence.Replay(aofPath, func(line string) {
		s.Apply(line)
	})
	if err != nil {
		log.Fatal(err)
	}

	// Handle Ctrl+C / termination
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Shutting down MiniRedis...")
		_ = s.Close()
		os.Exit(0)
	}()

	log.Println("MiniRedis listening on :6379")
	log.Fatal(s.Start())
}
