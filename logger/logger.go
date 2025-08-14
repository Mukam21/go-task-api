package logger

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Logger struct {
	ch     chan string
	wg     sync.WaitGroup
	mu     sync.RWMutex
	closed bool
}

func NewLogger(bufferSize int) *Logger {
	return &Logger{
		ch: make(chan string, bufferSize),
	}
}

func (l *Logger) Start() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for msg := range l.ch {
			log.Printf("%s %s", time.Now().UTC().Format(time.RFC3339), msg)
		}
	}()
}

func (l *Logger) Log(message string) {
	l.mu.RLock()
	if l.closed {
		l.mu.RUnlock()
		return
	}
	l.mu.RUnlock()

	select {
	case l.ch <- message:
	default:
		fmt.Println(time.Now().UTC().Format(time.RFC3339), "logger: channel full, dropping:", message)
	}
}

func (l *Logger) Close() {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return
	}
	l.closed = true
	close(l.ch)
	l.mu.Unlock()

	l.wg.Wait()
}
