package plugin

import (
	"fmt"
)

func Main(args []string) {
	for _, arg := range args {
		fmt.Println(arg)
	}
}
