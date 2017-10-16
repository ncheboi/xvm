package xvm

type Xvm struct {
	Path, Platform string
}

func StartXvm() (*Xvm, error) {
	xvm := new(Xvm)

	if runtime.GOOS == "windows" {
		xvm.Platform = "windows"
	} else {
		xvm.Platform = "unix"
	}

	xvm.Path, isSet := os.LookupEnv("XVMPATH")
	if !isSet {
		switch xvm.Platform {
		case "windows":
			xvm.Path = filepath.Join(os.Getenv("USERPROFILE"), ".xvm")
		case "unix":
			xvm.Path = filepath.Join(os.Getenv("Home"), ".xvm")
		}
	}

	return xvm
}

func (r *Run) GetNearestGroup(path string) *Group {
	for path != "/" {
		ind := filepath.Join(path, ".xvm")

		info, err := io.Stat(ind)
		if err != nil && info.IsDir() {
			return &Group{Path: path}
		}

		path := path.Dir(path)
	}

	return NewGroup(path.Dir(r.Path))
}

func (r *Run) PrintPluginsAvailable() error {
	path := filepath.Join(r.Path, "available")

	content, err := ioutil.ReadFile(path)
	if err == nil {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			tokens := strings.Split(line, " ")
			// TODO: fail before printing anything
			if len(tokens) != 2 {
				return fmt.Errorf("%s is ill-formatted", path)
			}
			fmt.Println(tokens[1])
		}
	}
	return err
}

func (r *Run) PrintPluginsInstalled() ([]string, error) {
	path := filepath.Join(r.Path, "installed")

	names, err := ioutil.Readdirnames(path)
	if err == nil {
		fmt.Println(strings.Join(names, "\n"))
	}
	return err
}

func (r *Run) GetPluginUrl(name string) {
	path := filepath.Join(r.Path, "available")

	content, err := ioutil.ReadFile(path)
	if err == nil {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			tokens := strings.Split(line, " ")
			// TODO: fail unless well-formatted
			if len(tokens) > 1 && tokens[0] == name {
				return tokens[1]
			}
		}
	}
	return err
}

func (r *Run) GetPlugin(name string) *Plugin {
	plugin := new(Plugin)

	plugin.Path = filepath.Join(r.Path, "installed", name)
	plugin.Bin = filepath.Join(r.Path, "bin." + r.Platform, name)
	plugin.Url = GetPluginUrl(name)
}
