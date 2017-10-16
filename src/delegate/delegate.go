package delegate

import (
	"fmt"
)

func Main(args []string) error {
	for _, arg := range args {
		fmt.Println(arg)
	}

	return nil
}
