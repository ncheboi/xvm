package plugin

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"

	"../utils"
)

func add(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("Too few arguments")
	}

	remotes, err := getRemotePlugins()
	if err != nil {
		return err
	}

	name := args[3]
	url, ok := remotes[name]
	if !ok {
		return fmt.Errorf("Plugin %s unknown", name)
	}

	dir, err := utils.GetXvmSubDir("plugins")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, name)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-stop
		fmt.Fprintf(os.Stderr, "Cancelling...\n")
		cancel()
	}()

	fmt.Printf("Installing %s. Push Ctrl-C to cancel.\n", name)

	_, err = git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}

	return nil
}
