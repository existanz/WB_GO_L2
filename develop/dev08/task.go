package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	sh := NewShell(os.Stdout, os.Stdin)
	if err := sh.Run(); err != nil {
		log.Fatal(err)
	}
}

type Shell struct {
	prompt []byte
	w      io.Writer
	r      io.Reader
}

func NewShell(w io.Writer, r io.Reader) *Shell {
	return &Shell{
		prompt: []byte("\033[32mmysh>\033[0m "),
		w:      w,
		r:      r,
	}
}

func (sh *Shell) Run() error {
	for {
		sh.w.Write(sh.prompt)
		input, err := bufio.NewReader(sh.r).ReadString('\n')
		if err != nil {
			sh.writeln("Error reading input:", err)
			continue
		}

		commands := strings.Split(input, "|")
		var processes []*exec.Cmd

		for _, command := range commands {
			args := strings.Fields(command)
			if len(args) == 0 {
				continue
			}
			if args[0] == "q" || args[0] == "quit" || args[0] == "e" || args[0] == "exit" {
				os.Exit(0)
			}
			if args[0] == "cd" {
				if len(args) < 2 {
					sh.writeln("cd: missing argument")
				} else {
					if err := os.Chdir(args[1]); err != nil {
						sh.writeln("cd:", err)
					}
				}
				continue
			}
			cmd := exec.Command(args[0], args[1:]...)
			if len(processes) > 0 {
				cmd.Stdin, _ = processes[len(processes)-1].StdoutPipe()
			}
			processes = append(processes, cmd)
		}

		sh.executeCommands(processes)
	}
}

func (sh *Shell) executeCommands(commands []*exec.Cmd) {
	if len(commands) == 0 {
		return
	}

	commands[len(commands)-1].Stdout = sh.w
	commands[len(commands)-1].Stderr = sh.w

	for _, cmd := range commands {
		if err := cmd.Start(); err != nil {
			sh.writeln("Error starting command:", err)
			break
		}
	}

	for _, cmd := range commands {
		if err := cmd.Wait(); err != nil {
			sh.writeln("Error waiting for command:", err)
		}
	}
}

func (sh *Shell) writeln(output ...any) {
	line := fmt.Sprint(output...)
	sh.w.Write([]byte(line + "\n"))
}
