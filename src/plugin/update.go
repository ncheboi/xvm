package plugin

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
)

func update(args []string) error {
	if len(args) < 4 {
		fmt.Errorf("Too fmt arguments")
	} else if len(args) > 4 {
		fmt.Errorf("Too many arguments")
	}

	name := args[3]
	path := filepath.Join("plugins", name)
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-stop
		fmt.Println("Cancelling...")
		cancel()
	}()

	fmt.Printf("Updating %s. Push Ctrl-C to cancel.\n", name)

	err = tree.PullContext(ctx, &git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}

	return nil
}
