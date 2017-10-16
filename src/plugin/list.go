package plugin

import (
	"fmt"
)

func list(args []string) error {
	remote, local := true, true

	if len(args) == 4 {
		switch args[3] {
		case "remote":
			local = false
		case "local":
			remote = false
		default:
			return fmt.Errorf("Expected local or remote, but got %s", args[3])
		}
	} else if len(args) > 4 {
		return fmt.Errorf("Too many arguments")
	}

	if remote {
		if local {
			fmt.Println("---- remote")
		}

		remotes, err := getRemotePlugins()
		if err != nil {
			return err
		}

		for remote, _ := range remotes {
			fmt.Println(remote)
		}
	}

	if local {
		if remote {
			fmt.Println("---- local")
		}

		locals, err := getLocalPlugins()
		if err != nil {
			return err
		}

		for local, _ := range locals {
			fmt.Println(local)
		}
	}

	return nil
}
