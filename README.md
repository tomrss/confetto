# Confetto

**Sugar-coated configuration for Go applications.**

Confetto loads configuration from multiple sources — YAML files, environment variables, and CLI flags — with a clear priority order and zero boilerplate. 

Define your config as a struct, set defaults and validators with a fluent API, and let Confetto do the rest.

## Features

- **Multiple sources**: YAML files, environment variables, CLI flags
- **Clear priority**: CLI > ENV > YAML > default
- **Type-safe**: `string`, `int`, `bool`, `float64`, `time.Duration`, and list variants
- **Validation**: built-in validators (`Range`, `OneOf`, `NotEmpty`, ...) or bring your own — fail fast!
- **Single dependency**: only `gopkg.in/yaml.v3`

## Usage

Confetto works in two steps:

1. **Define** your configuration as a Go struct. Each field's `cfg` tag determines its configuration key. Nested structs produce dotted keys (e.g. a `Host` field tagged `cfg:"host"` inside a struct tagged `cfg:"db"` becomes `db.host`). Use the fluent builder API to set defaults, validators, and required constraints.
2. **Load** with a single call. Confetto resolves each key by merging multiple sources in priority order — CLI flags > environment variables > YAML file > defaults — and validates everything at once.

### Install

```bash
go get github.com/tomrss/confetto
```

### Define your configuration

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/tomrss/confetto"
)

type DBConfig struct {
    Host    confetto.StringParam   `cfg:"host"`
    Port    confetto.IntParam      `cfg:"port"`
    Name    confetto.StringParam   `cfg:"name"`
    Timeout confetto.DurationParam `cfg:"timeout"`
}

type ServerConfig struct {
    Addr    confetto.StringParam `cfg:"addr"`
    Verbose confetto.BoolParam   `cfg:"verbose"`
}

type Config struct {
    DB     DBConfig     `cfg:"db"`
    Server ServerConfig `cfg:"server"`
}

func newConfig() Config {
    return Config{
        DB: DBConfig{
            Host:    confetto.String().Default("localhost").Build(),
            Port:    confetto.Int().Default(5432).Validate(confetto.Range(1, 65535)).Build(),
            Name:    confetto.String().Required().Validate(confetto.NotEmpty()).Build(),
            Timeout: confetto.Duration().Default(30 * time.Second).Build(),
        },
        Server: ServerConfig{
            Addr:    confetto.String().Default(":8080").Build(),
            Verbose: confetto.Bool().Default(false).Build(),
        },
    }
}

func main() {
    cfg := newConfig()
    err := confetto.Load(&cfg, confetto.Options{
        ConfigPaths: confetto.DefaultConfigPaths("myapp"),
        EnvPrefix:   "MYAPP",
        Args:        os.Args[1:],
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("DB host:", cfg.DB.Host.Get())
    fmt.Println("DB port:", cfg.DB.Port.Get())
    fmt.Println("Server addr:", cfg.Server.Addr.Get())

    // IsSet() returns true only if the value came from a source, not from the default
    if cfg.Server.Verbose.IsSet() {
        fmt.Println("verbose mode was explicitly enabled")
    }
}
```

### YAML file

The YAML structure mirrors the nesting of your `cfg` tags. A key like `db.host` maps to:

```yaml
# config.yaml
db:
  host: production.db.com
  port: 5433
  name: mydb
  timeout: 1m
server:
  addr: ":9090"
```

Each level of dot-separation in the key (`db.host`, `server.addr`) becomes a level of YAML nesting.

Where the config file will be searched can be configured.

`DefaultConfigPaths` returns conventional paths for a given app name:

```go
confetto.DefaultConfigPaths("myapp")
// -> ["config.yaml", "~/.config/myapp/config.yaml", "/etc/myapp/config.yaml"]
```

Or specify exact paths — the first existing file wins:

```go
confetto.Options{
    ConfigPaths: []string{"./config.yaml", "/etc/myapp/config.yaml"},
}
```

### Environment variables

Environment variables are derived from the key by uppercasing and replacing `.` with `_`, then prepending the configured prefix:

```bash
export MYAPP_DB_HOST=env.db.com
export MYAPP_DB_PORT=5434
export MYAPP_SERVER_ADDR=":3000"
```

### CLI flags

Flags use the dotted key directly, with `--key=value` or `--key value` syntax:

```bash
./myapp --db.host=cli.db.com --db.port 5435 --server.verbose
```

Confetto only picks up flags matching your `cfg` keys — it won't interfere with other flags or subcommands in your CLI.

### Source priority

When the same key is set in multiple sources, the highest-priority source wins:

1. **CLI flags** (highest)
2. **Environment variables**
3. **YAML file**
4. **Default value** (lowest)

### List parameters

List types (`StringListParam`, `IntListParam`, etc.) are supported. In YAML they map to arrays; in ENV/CLI they are split by a configurable separator (default `,`):

```go
type Config struct {
    Tags  confetto.StringListParam `cfg:"tags"`
    Ports confetto.IntListParam    `cfg:"ports"`
}

cfg := Config{
    Tags:  confetto.StringList().Default([]string{}).Build(),
    Ports: confetto.IntList().Default([]int{}).Build(),
}
```

```bash
# CLI
./myapp --tags=alpha,beta,gamma --ports=8080,9090

# ENV
export MYAPP_TAGS="alpha,beta,gamma"
```

The default list separator is `,`. You can change it with:

```go
confetto.Options{ListSeparator: ";"}
```

```yaml
# YAML
tags:
  - alpha
  - beta
  - gamma
ports:
  - 8080
  - 9090
```

### Validation

Use built-in validators or pass any `func(T) error`:

```go
confetto.Int().Validate(confetto.Range(1, 65535)).Build()
confetto.Float().Validate(confetto.RangeFloat(0, 1)).Build()
confetto.String().Validate(confetto.OneOf("dev", "staging", "prod")).Build()
confetto.String().Validate(confetto.NotEmpty()).Build()
confetto.String().Validate(confetto.MinLen(3)).Build()
confetto.StringList().Validate(confetto.MinItems[string](1)).Build()
confetto.Int().Validate(confetto.Positive()).Build()

// custom validator
confetto.String().Validate(func(s string) error {
    if s[0] == '/' {
        return fmt.Errorf("path must not be absolute")
    }
    return nil
}).Build()
```

### Error handling

All errors (parse, validation, required) are collected into a single `LoadError`:

```
3 configuration errors:
  - required parameter "db.name" is not set
  - failed to parse "abc" as int for key "db.port": strconv.Atoi: parsing "abc": invalid syntax
  - validation failed for "db.max_conns" (value: 99999): validation error: value 99999 is not in range [1, 65535]
```

## Development

Prerequisites: Go 1.25+

```bash
# run tests
make test

# run linter
make lint

# run vulnerability check
make vulncheck

# run all checks
make all

# auto-fix lint issues
make fix

# coverage report
make cover-report
```

## License

TODO
