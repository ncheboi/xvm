# xvm
The X Version Manager: solve for X.

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

And add `.xvm\bin` to your path:

```powershell
[Environment]::SetEnvironmentVariable("Path", "$env:USERPROFILE\.xvm\win32;" + $env:Path, [System.EnvironmentVariableTarget]::Machine)
```

### Unix (Mac OSX, Linux or BSD)

Clone the source into your home directory:

```bash
git clone https://github.com/skotchpine/xvm .xvm
```

And add the following to your `~/.bashrc`:

```bash
echo -e 'PATH=~/.xvm/bin.unix:$PATH\n' >> ~/.bashrc
```
