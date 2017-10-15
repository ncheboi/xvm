# xvm
The X Version Manager: where X is whatever you need.

## Disclaimer

The library is at a very early stage in development.
You will find bugs, catastrophic and irreversible bugs.
If you aren't comfortable with losing all the things and building
your system from scratch, then please wait for a beta release.

## Installation

### Windows

In powershell, clone the source into your user profile directory:

```powershell
git clone https://github.com/skotchpine/xvm "$env:USERPROFILE/.xvm"
```

Add `.xvm\bin` to your path:

```powershell
[Environment]::SetEnvironmentVariable(
	"Path",
	$env:Path + ";$env:USERPROFILE\.xvm\win32",
	[System.EnvironmentVariableTarget]::Machine
)
```

### Unix (Mac OSX, Linux or BSD)

Clone the source into your home directory:

```bash
git clone https://github.com/skotchpine/xvm .xvm
```

Add the following to your `~/.bashrc`:

```bash
echo -e 'PATH=$PATH:~/.xvm/unix\n' >> ~/.bashrc
```

## Usage

### Group

```bash
xvm group init [dir]
```

```bash
xvm group status [dir]
```

```bash
xvm group remove [dir]
```

### Plugins

```bash
xvm plugin list
```

```bash
xvm plugin add js
```

```bash
xvm plugin update js
```

```bash
xvm plugin remove js
```

### Versions

```bash
xvm list js
```

```bash
xvm install js 8.7.0
```

```bash
xvm uninstall js 8.7.0
```

```bash
xvm set js [dir] 8.7.0
```

```bash
xvm unset js [dir]
```

### The `.xvm/config` File

Add the following lines to `.xvm/config` as needed.

- `gitignore`
- `no-confirm-remove`

### Common Things

- Keeping `.xvm` out of source control.

	```bash
	echo -e '.xvm\n' >> .gitignore
	```

## Development

### Hooks

### The Plugin Interface

```bash
my-plugin list
```

```bash
my-plugin install
```
