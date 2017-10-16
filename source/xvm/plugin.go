package xvm

type Plugin struct {
	Path, Bin, Url string
}

func (p *Plugin) PrintAvailable() error {
	path := filepath.Join(p.Path, "available")

	versions, err := ioutil.ReadFile(path)
	if err == nil {
		fmt.Println(versions)
	}
	return err
}

func (p *Plugin) PrintInstalled() error {
	path := filepath.Join(p.Path, "installed")

	versions, err := ioutil.Readdirnames(path)
	if err == nil {
		fmt.Println(strings.Join(versions, "\n"))
	}
	return err
}

func (p *Plugin) Install(version string) error {
	cmd := exec.Command(p.Bin, "install", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *Plugin) Uninstall(version string) error {
	cmd := exec.Command(p.Bin, "uninstall", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *Plugin) Update() error {
	resp, err := http.Get(p.Url)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return fmt.Errorf("Failed to get %s", p.Url)
	}

	body, err := zip.OpenReader(resp.Body)
	if err != nil {
		return err
	}
	defer body.Close()

	for _, z := range r.File {
		r, err := z.Open()
		if err != nil {
			return err
		}
		defer r.Close()

		path := filepath.Join(p.Path, z.Name)
		if z.FileInfo().IsDir() {
			return os.MkdirAll(path, os.ModePerm)
		}

		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		mode := os.O_WRONLY|os.O_CREATE|os.O_TRUNC
		f, err := os.OpenFile(path, mode, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, r)
		return err
	}
}

func (p *Plugin) Remove() error {
	return os.RemoveAll(p.Path)
}
