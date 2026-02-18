package confetto

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// source represents a configuration source.
type source interface {
	// get returns the value for a key, or nil if not found.
	get(key string) any
}

// cliSource parses command line arguments.
type cliSource struct {
	values map[string]string
}

func newCLISource(args []string) *cliSource {
	s := &cliSource{values: make(map[string]string)}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		arg = strings.TrimPrefix(arg, "--")

		// handle --key=value format
		if idx := strings.Index(arg, "="); idx != -1 {
			key := arg[:idx]
			value := arg[idx+1:]
			s.values[key] = value
			continue
		}

		// handle --key value format
		key := arg
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			s.values[key] = args[i+1]
			i++
		} else {
			// flag without value (boolean)
			s.values[key] = "true"
		}
	}
	return s
}

func (s *cliSource) get(key string) any {
	if v, ok := s.values[key]; ok {
		return v
	}
	return nil
}

// envSource reads from environment variables.
type envSource struct {
	prefix string
}

func newEnvSource(prefix string) *envSource {
	return &envSource{prefix: prefix}
}

func (s *envSource) get(key string) any {
	// convert key to env var format: db.host -> DB_HOST
	envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	if s.prefix != "" {
		envKey = s.prefix + "_" + envKey
	}
	if v, ok := os.LookupEnv(envKey); ok {
		return v
	}
	return nil
}

// yamlSource reads from a YAML file.
type yamlSource struct {
	data map[string]any
}

func newYAMLSource(filename string) (*yamlSource, error) {
	s := &yamlSource{data: make(map[string]any)}
	if filename == "" {
		return s, nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(content, &s.data); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *yamlSource) get(key string) any {
	parts := strings.Split(key, ".")
	var current any = s.data

	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current, ok = m[part]
		if !ok {
			return nil
		}
	}

	return current
}
