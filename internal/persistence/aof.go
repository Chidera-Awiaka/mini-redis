package persistence

import (
	"bufio"
	"os"
	"sync"
)

type AOF struct {
	mu   sync.Mutex
	file *os.File
	w    *bufio.Writer
}

func Open(path string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &AOF{
		file: f,
		w:    bufio.NewWriter(f),
	}, nil
}

func (a *AOF) Append(line string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, err := a.w.WriteString(line + "\n")
	if err != nil {
		return err
	}
	return a.w.Flush()
}

func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_ = a.w.Flush()
	return a.file.Close()
}

// Replay reads the AOF file line-by-line and calls apply(line) for each line.
// This is used at startup to rebuild the in-memory state.
func Replay(path string, apply func(line string)) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		apply(line)
	}

	return sc.Err()
}
