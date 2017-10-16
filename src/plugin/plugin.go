package plugin

import (
	"fmt"
)

func Main(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("Too few arguments")
	}

	switch args[2] {
	case "list":
		if err := list(args); err != nil {
			return err
		}
	case "add":
		if err := add(args); err != nil {
			return err
		}
	case "update":
		if err := update(args); err != nil {
			return err
		}
	case "remove":
		if err := remove(args); err != nil {
			return err
		}
	default:
		if err := list(args); err != nil {
			return err
		}
	}

	return nil
}
