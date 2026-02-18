package confetto

import (
	"os"
	"path/filepath"
	"reflect"
)

// Options configures the configuration loader.
type Options struct {
	// ConfigFile is the explicit path to a YAML configuration file.
	// If set, this takes precedence over ConfigPaths.
	ConfigFile string
	// ConfigPaths is a list of paths to search for the config file.
	// The first existing file is used. Ignored if ConfigFile is set.
	ConfigPaths []string
	// EnvPrefix is the prefix for environment variables.
	EnvPrefix string
	// Args are the command line arguments to parse.
	Args []string
	// ListSeparator is the separator for list values in strings (default: ",").
	ListSeparator string
}

// FindConfigFile searches for a configuration file in the given paths.
// Returns the path of the first existing file, or empty string if none found.
func FindConfigFile(paths []string) string {
	for _, p := range paths {
		// expand ~ to home directory
		if len(p) > 0 && p[0] == '~' {
			if home, err := os.UserHomeDir(); err == nil {
				p = filepath.Join(home, p[1:])
			}
		}
		// expand environment variables
		p = os.ExpandEnv(p)

		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// DefaultConfigPaths returns the default paths to search for config files.
// The order is: XDG config home, then system-wide /etc.
func DefaultConfigPaths(appName string) []string {
	paths := make([]string, 0, 3)

	// this folder
	paths = append(paths, "config.yaml")

	// XDG config home (or ~/.config)
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = "~/.config"
	}
	paths = append(paths, filepath.Join(xdgConfig, appName, "config.yaml"))

	// system-wide
	paths = append(paths, filepath.Join("/etc", appName, "config.yaml"))

	return paths
}

// Load loads configuration from multiple sources into the provided struct.
// The struct must contain fields that implement the Param interface.
// Sources are checked in order of priority: CLI > ENV > YAML > default.
func Load(cfg any, opts Options) error {
	if opts.ListSeparator == "" {
		opts.ListSeparator = ","
	}

	// determine config file path
	configFile := opts.ConfigFile
	if configFile == "" && len(opts.ConfigPaths) > 0 {
		configFile = FindConfigFile(opts.ConfigPaths)
	}

	// initialize sources
	cliSrc := newCLISource(opts.Args)
	envSrc := newEnvSource(opts.EnvPrefix)
	yamlSrc, err := newYAMLSource(configFile)
	if err != nil {
		return err
	}

	sources := []source{cliSrc, envSrc, yamlSrc}

	// collect all params via reflection
	params := collectParams(cfg, "")

	loadErr := &LoadError{}

	// for each param, try to load from sources in priority order
	for _, p := range params {
		loadParam(p, sources, opts, loadErr)
	}

	if loadErr.HasErrors() {
		return loadErr
	}
	return nil
}

func loadParam(p Param, sources []source, opts Options, loadErr *LoadError) {
	key := p.key()
	var value any

	for _, src := range sources {
		if v := src.get(key); v != nil {
			value = v
			break
		}
	}

	if value != nil {
		var setErr error
		if s, ok := value.(string); ok {
			setErr = p.setFromString(s, opts.ListSeparator)
		} else {
			setErr = p.setFromAny(value, opts.ListSeparator)
		}
		if setErr != nil {
			loadErr.Add(setErr)
			return
		}
	}

	if err := p.validate(); err != nil {
		loadErr.Add(err)
	}

	if p.isRequired() && !p.IsSet() && !p.hasDefault() {
		loadErr.Add(&RequiredError{Key: key})
	}
}

// collectParams walks the struct and collects all Param fields with their keys.
func collectParams(v any, prefix string) []Param {
	var params []Param

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return params
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// get the cfg tag
		tag := fieldType.Tag.Get("cfg")
		if tag == "" && !fieldType.Anonymous {
			// skip fields without cfg tag (unless embedded)
			continue
		}

		// build the key
		key := tag
		if prefix != "" && key != "" {
			key = prefix + "." + key
		} else if prefix != "" {
			key = prefix
		}

		// check if field implements Param
		if field.CanAddr() {
			if p, ok := field.Addr().Interface().(Param); ok {
				p.setKey(key)
				params = append(params, p)
				continue
			}
		}

		// recurse into nested structs
		if field.Kind() == reflect.Struct {
			nested := collectParams(field.Addr().Interface(), key)
			params = append(params, nested...)
		}
	}

	return params
}
