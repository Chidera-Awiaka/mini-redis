package server

import (
	"bufio"
	"fmt"
	"mini-redis/internal/persistence"
	"mini-redis/internal/store"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	addr  string
	store *store.Store
	aof   *persistence.AOF
}

func New(addr string, aof *persistence.AOF) *Server {
	return &Server{
		addr:  addr,
		store: store.New(),
		aof:   aof,
	}
}

func (s *Server) Close() error {
	if s.aof != nil {
		return s.aof.Close()
	}
	return nil
}

func (s *Server) Apply(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	parts := strings.Split(line, " ")
	switch strings.ToUpper(parts[0]) {
	case "SET":
		if len(parts) == 3 {
			s.store.Set(parts[1], parts[2])
			return
		}
		if len(parts) == 5 && strings.ToUpper(parts[3]) == "EX" {
			ttl, err := strconv.Atoi(parts[4])
			if err != nil {
				return
			}
			s.store.SetWithTTL(parts[1], parts[2], ttl)
			return
		}
	case "DEL":
		if len(parts) == 2 {
			s.store.Delete(parts[1])
			return
		}
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Fprintln(conn, "Welcome to MiniRedis v0.5 (Graceful shutdown). Commands: PING, SET, GET, DEL, STATS, QUIT")

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		switch strings.ToUpper(parts[0]) {

		case "PING":
			fmt.Fprintln(conn, "PONG")

		case "QUIT":
			fmt.Fprintln(conn, "BYE")
			return

		case "SET":
			if len(parts) == 3 {
				s.store.Set(parts[1], parts[2])
				if s.aof != nil {
					_ = s.aof.Append(line)
				}
				fmt.Fprintln(conn, "OK")
				continue
			}

			if len(parts) == 5 && strings.ToUpper(parts[3]) == "EX" {
				ttl, err := strconv.Atoi(parts[4])
				if err != nil {
					fmt.Fprintln(conn, "ERR invalid TTL")
					continue
				}
				s.store.SetWithTTL(parts[1], parts[2], ttl)
				if s.aof != nil {
					_ = s.aof.Append(line)
				}
				fmt.Fprintln(conn, "OK")
				continue
			}

			fmt.Fprintln(conn, "ERR usage: SET key value [EX seconds]")

		case "GET":
			if len(parts) != 2 {
				fmt.Fprintln(conn, "ERR usage: GET key")
				continue
			}
			val, ok := s.store.Get(parts[1])
			if !ok {
				fmt.Fprintln(conn, "NULL")
				continue
			}
			fmt.Fprintln(conn, val)

		case "DEL":
			if len(parts) != 2 {
				fmt.Fprintln(conn, "ERR usage: DEL key")
				continue
			}
			s.store.Delete(parts[1])
			if s.aof != nil {
				_ = s.aof.Append(line)
			}
			fmt.Fprintln(conn, "OK")

		case "STATS":
			fmt.Fprintln(conn, s.store.Stats())

		default:
			fmt.Fprintln(conn, "ERR unknown command")
		}
	}
}
