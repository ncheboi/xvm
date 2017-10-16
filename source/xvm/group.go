package xvm

type Group struct {
	Path string
}

func (g *Group) PrintAllVersions() error {
	versions, err := g.GetAllVersions()
	if err == nil {
		for plugin, version := range versions {
			fmt.Printf("%s %s\n", plugin, version)
		}
	}
	return err
}

func (g *Group) GetAllVersions() (map[string]string, error) {
	path := filepath.Join(g.Path, "versions")
	plugins, err := ioutil.Readdirnames(path)
	if err != nil {
		return nil, err
	}

	versions := make(map[string]string)

	for _, plugin := range plugins {
		version, err := g.GetVersion(plugin)
		if err != nil {
			return nil, err
		}

		versions[plugin] = version
	}

	return versions, nil
}

func (g *Group) GetVersion(plugin string) (string, error) {
	path := filepath.Join(g.Path, "versions", plugin)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	n := bytes.Count(content, []byte{'\n'})
	if n > 1 {
		return nil, fmt.Errorf("%s is ill-formatted", path)
	}

	return strings.Trim(string(content), "\n"), nil
}

func (g *Group) SetVersion(plugin string, version string) error {
	path := filepath.Join(g.Path, "versions", plugin)
	return ioutil.WriteFile(path, []byte(version), os.ModePerm)
}

func (g *Group) UnsetVersion(plugin string) error {
	path := filepath.Join(g.Path, "versions", plugin)
	return os.RemoveAll(path)
}

func (g *Group) Remove() error {
	return os.RemoveAll(g.Path)
}
