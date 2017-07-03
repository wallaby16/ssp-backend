package gluster

import (
	"errors"
	"log"
)

func executeCommandsLocally(commands []string) error {
	log.Println("Got new commands to execute:")
	for _, c := range commands {
		out, err := ExecRunner.Run("bash", "-c", c)
		if err != nil {
			log.Println("Error executing command: ", c, err.Error(), string(out))
			return errors.New(commandExecutionError)
		}
		log.Printf("Cmd: %v | StdOut: %v", c, string(out))
	}

	return nil
}
