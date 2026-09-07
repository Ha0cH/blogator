package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmdLookUp map[string]func(*state, command) error
}

func (cmds *commands) run(s *state, cmd command) error {
	c, ok := cmds.cmdLookUp[cmd.name]
	if !ok {
		return fmt.Errorf("Command not found")
	}

	return c(s, cmd)
}

func (cmds *commands) register(name string, f func(*state, command) error) {
	cmds.cmdLookUp[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Login command expects a single argument: Username")
	}

	userName := cmd.args[0]
	err := s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Println("The user has been set.")
	return nil
}
