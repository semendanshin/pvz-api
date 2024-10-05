package stdin

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type HandlerFunc func(ctx context.Context, args []string) (string, error)

type Command string

const (
	HelpCommand Command = "help"
)

func NewCommand(name string) Command {
	return Command(name)
}

type Server struct {
	running      bool
	done         chan struct{}
	handlers     map[Command]HandlerFunc
	numOfWorkers int
}

func NewServer(numOfWorkers int) *Server {
	if numOfWorkers < 1 {
		numOfWorkers = 1
	}

	return &Server{
		done:         make(chan struct{}),
		handlers:     make(map[Command]HandlerFunc),
		numOfWorkers: numOfWorkers,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.AddHandler(HelpCommand, s.helpHandler)

	s.running = true

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	input := make(chan string, s.numOfWorkers)
	output := make(chan string)

	s.runInput(ctx, input)
	s.runOutput(ctx, output)

	for i := 0; i < s.numOfWorkers; i++ {
		s.runDispatch(ctx, input, output)
	}

	output <- "Welcome to the CLI. Type 'help' to see available commands."

	<-s.done

	return nil
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

func (s *Server) runInput(ctx context.Context, intake chan string) {
	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			select {
			case <-ctx.Done():
				close(intake)
				return
			default:
				msg, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("failed to read input:", err)
					continue
				}

				msg = strings.TrimSpace(msg)

				intake <- msg
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
				go func() {
					output <- s.handleInput(ctx, msg)
				}()
			}
		}
	}()
}

func (s *Server) AddHandler(command Command, handler HandlerFunc) {
	s.handlers[command] = handler
}

func (s *Server) handleCommand(ctx context.Context, command Command, args []string) (string, error) {
	handler, ok := s.handlers[command]
	if !ok {
		var availableCommands []string
		for k := range s.handlers {
			availableCommands = append(availableCommands, string(k))
		}
		return "", fmt.Errorf("command %s not found. Available commands: %v", command, availableCommands)
	}

	return handler(ctx, args)
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

func (s *Server) helpHandler(ctx context.Context, args []string) (string, error) {
	var commands []string
	for k := range s.handlers {
		commands = append(commands, string(k))
	}

	return fmt.Sprintf("Available commands: %v", commands), nil
}
