# Configo

Package for parse and merge config into struct from:
- env
- file (json, yaml)
- default params (struct fields with tag `default`)

## Usage:

For usage, you need declare your config struct with field tags `config` and if you want `default`

```go
type Config struct {
	Http struct {
		Address string `configo:"address" default:"127.0.0.1"`
		Port    string `configo:"port" default:"80"`
    }`config:"http"`
}
```

then you can use Option struct for parametrize and parse config:

```go
config, err := configo.Parse[Config](configo.Option{})
```
after that you can you that var:

```go
http.ListenAndServe(net.JoinHostPort(config.Http.Address, config.Http.Port), nil)
```