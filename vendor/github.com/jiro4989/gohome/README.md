# gohome

gohome is provide functions to return path of user home directory.

## usage

`GetHomeDir()` returns `$HOME` or `$USERPROFILE`.

`GetConfigDir()` returns `$HOME/.config` or `$USERPROFILE/.config`.

```go
homeDir, err := gohome.GetHomeDir()

configDir, err := gohome.GetConfigDir()
```

## install package

```bash
go get github.com/jiro4989/gohome
```
