package stdin

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type HandlerFunc func(ctx context.Context, args []string) (string, error)

type Command string

const (
	HelpCommand  Command = "help"
	WorkersCount Command = "workers-count"
)

func NewCommand(name string) Command {
	return Command(name)
}

type Server struct {
	running bool
	done    chan struct{}

	handlers map[Command]HandlerFunc

	inputChan  chan string
	outputChan chan string

	numOfWorkers      int
	workersCancelFunc []context.CancelFunc

	mu sync.RWMutex
}

func NewServer(numOfWorkers int) *Server {
	if numOfWorkers < 1 {
		numOfWorkers = 1
	}

	return &Server{
		done: make(chan struct{}),

		handlers: make(map[Command]HandlerFunc),

		numOfWorkers: numOfWorkers,

		inputChan:  make(chan string, 10),
		outputChan: make(chan string),
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.AddHandler(HelpCommand, s.helpHandler)
	s.AddHandler(WorkersCount, s.setWorkersCountHandler)

	s.running = true

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.runInput(ctx, s.inputChan)
	s.runOutput(ctx, s.outputChan)

	for i := 0; i < s.numOfWorkers; i++ {
		s.startWorker(ctx, s.inputChan, s.outputChan)
	}

	s.outputChan <- "Welcome to the CLI. Type 'help' to see available commands."

	<-s.done

	s.stopWorkers()

	return nil
}

func (s *Server) startWorker(ctx context.Context, input chan string, output chan string) {
	ctx, cancel := context.WithCancel(ctx)

	s.runDispatch(ctx, input, output)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.workersCancelFunc = append(s.workersCancelFunc, cancel)
}

func (s *Server) stopWorker() {
	if len(s.workersCancelFunc) == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cancel := s.workersCancelFunc[len(s.workersCancelFunc)-1]
	go cancel()

	s.workersCancelFunc = s.workersCancelFunc[:len(s.workersCancelFunc)-1]

	return
}

func (s *Server) stopWorkers() {
	for _, cancel := range s.workersCancelFunc {
		cancel()
	}
}

func (s *Server) Stop() {
	s.running = false
	close(s.done)
}

func (s *Server) runOutput(ctx context.Context, input chan string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-input:
				fmt.Println("> " + msg)
				fmt.Println()
			}
		}
	}()
}

func (s *Server) runInput(ctx context.Context, input chan string) {
	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			select {
			case <-ctx.Done():
				close(input)
				return
			default:
				msg, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("failed to read input:", err)
					continue
				}

				msg = strings.TrimSpace(msg)

				input <- msg
			}
		}
	}()
}

func (s *Server) runDispatch(ctx context.Context, input chan string, output chan string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-input:
				output <- s.handleInput(ctx, msg)
			}
		}
	}()
}

func (s *Server) AddHandler(command Command, handler HandlerFunc) {
	// Наверное лучше будет здесь проверять флаг running и возвращать ошибку, если сервер уже запущен
	// Но в задании надо использовать мьютексы, так что пусть будет так
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers[command] = handler
}

func (s *Server) handleInput(ctx context.Context, msg string) string {
	parts := strings.Split(msg, " ")
	if len(parts) == 0 {
		return "empty input"
	}

	command := Command(parts[0])
	args := parts[1:]

	output, err := s.handleCommand(ctx, command, args)
	if err != nil {
		return err.Error()
	}

	return output
}

func (s *Server) handleCommand(ctx context.Context, command Command, args []string) (string, error) {
	s.mu.RLock()
	handler, ok := s.handlers[command]

	if !ok {
		var availableCommands []string
		for k := range s.handlers {
			availableCommands = append(availableCommands, string(k))
		}
		s.mu.RUnlock()
		return "", fmt.Errorf("command %s not found. Available commands: %v", command, availableCommands)
	}

	s.mu.RUnlock()

	return handler(ctx, args)
}

func (s *Server) helpHandler(ctx context.Context, args []string) (string, error) {
	var commands []string
	for k := range s.handlers {
		commands = append(commands, string(k))
	}

	return fmt.Sprintf("Available commands: %v", commands), nil
}

func (s *Server) setWorkersCountHandler(ctx context.Context, args []string) (string, error) {
	usage := "Usage: workers-count <count>"

	if len(args) == 0 {
		return "", fmt.Errorf("no arguments provided. %s", usage)
	}

	workersCountStr := args[0]

	workersCount, err := strconv.Atoi(workersCountStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse workers count: %w", err)
	}

	if workersCount < 1 {
		return "", fmt.Errorf("workers count should be greater than 0")
	}

	diff := s.numOfWorkers - workersCount

	if diff > 0 {
		for i := 0; i < diff; i++ {
			s.stopWorker()
		}
	} else {
		for i := 0; i < -diff; i++ {
			s.startWorker(ctx, s.inputChan, s.outputChan)
		}
	}

	s.numOfWorkers = workersCount

	return fmt.Sprintf("workers count set to %d", workersCount), nil
}
