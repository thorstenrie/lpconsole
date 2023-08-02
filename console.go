// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package lpconsole

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/thorstenrie/tserr"
	"github.com/thorstenrie/tsfio"
	"github.com/thorstenrie/tstable"
)

type CommandFunc func(context.Context, []string) error

type Command struct {
	Key      string      // Keyword in the command line interface
	Help     string      // Help text for the command
	Function CommandFunc // Execution function
}

type Help struct {
	App, Description, Version string
}

type runner struct {
	help  *Help
	cmds  map[string]*Command
	exit  *Command
	stdin *os.File
}

var (
	help     = "help"
	helptext = "Print usage statement"
	run      = runner{cmds: make(map[string]*Command), stdin: os.Stdin}
)

func Usage(h *Help) error {
	if h.App != tsfio.Printable(h.App) {
		return tserr.NonPrintable("application name")
	}
	if h.Description != tsfio.Printable(h.Description) {
		return tserr.NonPrintable("application description")
	}
	if h.Version != tsfio.Printable(h.Version) {
		return tserr.NonPrintable("application version")
	}
	run.help = h
	Add(&Command{Key: help, Help: helptext, Function: printHelp})
	return nil
}

func Add(cmd *Command) error {
	if cmd.Function == nil {
		return tserr.NilPtr()
	}
	if cmd.Key == "" {
		return tserr.Empty("key")
	}
	if cmd.Key != tsfio.Printable(cmd.Key) {
		return tserr.NonPrintable("key")
	}
	if _, e := find(cmd.Key); e == nil {
		return tserr.Duplicate("key")
	}
	run.cmds[cmd.Key] = cmd
	return nil
}

func Exit(cmd string) error {
	c, e := find(cmd)
	if e != nil {
		return tserr.NotExistent("command")
	}
	run.exit = c
	return nil
}

func Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chin := make(chan string)
	chstop := make(chan error)

	go input(ctx, chin, chstop)

	for {
		fmt.Print("< ")

		select {
		case i, ok := <-chin:
			fmt.Print("> ")
			if !ok {
				e := tserr.Empty("input")
				fmt.Println(e)
				return e
			}
			cmd, args, e := split(i)
			if e != nil {
				fmt.Println(e) // TODO
			}
			c, e := find(cmd)
			if e != nil {
				fmt.Println(e) // TODO
				continue
			}
			e = c.Function(ctx, args)
			if e != nil {
				fmt.Println(e) // TODO
			}
			if c == run.exit {
				return e
			}
		case <-ctx.Done():
			fmt.Print("\n> ")
			e := exit(ctx)
			if e != nil {
				fmt.Println(e) // TODO
			}
			return e
		case err := <-chstop:
			if err != nil {
				fmt.Println(err) // TODO
			}
			fmt.Print("\n> ")
			e := exit(ctx)
			if e != nil {
				fmt.Println(e) // TODO
			}
			return e
		}
	}
}

// input scans Stdin and sends the input to the string channel chin. It sends an error to the error channel chstop, if any.
func input(ctx context.Context, chin chan string, chstop chan error) {
	s := bufio.NewScanner(run.stdin)
	for s.Scan() {
		chin <- s.Text()
		select {
		case <-ctx.Done():
			// Break for loop in case context is done.
			break
		default:
			// If the context is not done, scan next token.
		}
	}
	if err := s.Err(); err != nil {
		chstop <- err
	}
	close(chstop)
	close(chin)
}

func SetInput(in *os.File) error {
	if in == nil {
		return tserr.NilPtr()
	}
	run.stdin = in
	return nil
}

// split splits line l into the command name and its arguments.
func split(l string) (string, []string, error) {
	a := strings.Fields(tsfio.Printable(l))
	if len(a) == 0 {
		return "", nil, tserr.Empty("line")
	} else if len(a) == 1 {
		return a[0], nil, nil
	}
	return a[0], a[1:], nil
}

func output(format string, a ...any) error {
	_, err := fmt.Printf(format+"\n", a...)
	return err
}

func find(cmd string) (*Command, error) {
	if f, ok := run.cmds[cmd]; ok {
		return f, nil
	}
	return nil, tserr.NotExistent(cmd)
}

func printHelp(ctx context.Context, args []string) error {
	// Todo: No arguments allowed
	text := "\n"
	if run.help != nil {
		text += fmt.Sprintln(run.help.App + " " + run.help.Version + "\n")
		text += fmt.Sprintln("Description:\n\n" + "  " + run.help.Description + "\n")
	}
	text += "Usage:\n\n  command [arguments]\n\n"
	if len(run.cmds) > 0 {
		text += "Available commands:\n"
		t, e := tstable.New([]string{"command", "usage"})
		if e != nil {
			return tserr.Op(&tserr.OpArgs{Op: "New", Fn: "table", Err: e})
		}
		for k, c := range run.cmds {
			e = t.AddRow([]string{k, c.Help})
			if e != nil {
				return tserr.Op(&tserr.OpArgs{Op: "AddRow", Fn: "table", Err: e})
			}
		}
		t.SetGrid(&tstable.EmptyGrid)
		ts, e := t.Print()
		if e != nil {
			return tserr.Op(&tserr.OpArgs{Op: "Print", Fn: "table", Err: e})
		}
		fmt.Print(text + ts)
	}
	return nil
}

func exit(ctx context.Context) error {
	if run.exit == nil {
		return tserr.NotExistent("exit function")
	}
	return run.exit.Function(ctx, nil)
}
