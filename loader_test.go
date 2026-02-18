package confetto

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type testDBConfig struct {
	Host     StringParam   `cfg:"host"`
	Port     IntParam      `cfg:"port"`
	UseSSL   BoolParam     `cfg:"use_ssl"`
	Timeout  DurationParam `cfg:"timeout"`
	MaxConns IntParam      `cfg:"max_conns"`
}

type testServerConfig struct {
	Addr    StringParam `cfg:"addr"`
	Verbose BoolParam   `cfg:"verbose"`
}

type testConfig struct {
	DB     testDBConfig     `cfg:"db"`
	Server testServerConfig `cfg:"server"`
}

func newTestConfig() testConfig {
	return testConfig{
		DB: testDBConfig{
			Host:     String().Default("localhost").Build(),
			Port:     Int().Default(5432).Validate(Range(1, 65535)).Build(),
			UseSSL:   Bool().Default(false).Build(),
			Timeout:  Duration().Default(30 * time.Second).Build(),
			MaxConns: Int().Default(10).Build(),
		},
		Server: testServerConfig{
			Addr:    String().Default(":8080").Build(),
			Verbose: Bool().Default(false).Build(),
		},
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	cfg := newTestConfig()
	err := Load(&cfg, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DB.Host.Get() != "localhost" {
		t.Errorf("expected localhost, got %s", cfg.DB.Host.Get())
	}
	if cfg.DB.Port.Get() != 5432 {
		t.Errorf("expected 5432, got %d", cfg.DB.Port.Get())
	}
	if cfg.DB.UseSSL.Get() != false {
		t.Errorf("expected false, got %v", cfg.DB.UseSSL.Get())
	}
	if cfg.DB.Timeout.Get() != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.DB.Timeout.Get())
	}
	if cfg.Server.Addr.Get() != ":8080" {
		t.Errorf("expected :8080, got %s", cfg.Server.Addr.Get())
	}
}

func TestLoad_CLIArgs(t *testing.T) {
	cfg := newTestConfig()
	args := []string{
		"--db.host=production.db.com",
		"--db.port", "5433",
		"--db.use_ssl=true",
		"--server.verbose",
	}

	err := Load(&cfg, Options{Args: args})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DB.Host.Get() != "production.db.com" {
		t.Errorf("expected production.db.com, got %s", cfg.DB.Host.Get())
	}
	if cfg.DB.Port.Get() != 5433 {
		t.Errorf("expected 5433, got %d", cfg.DB.Port.Get())
	}
	if cfg.DB.UseSSL.Get() != true {
		t.Errorf("expected true, got %v", cfg.DB.UseSSL.Get())
	}
	if cfg.Server.Verbose.Get() != true {
		t.Errorf("expected true, got %v", cfg.Server.Verbose.Get())
	}
	if !cfg.DB.Host.IsSet() {
		t.Error("expected DB.Host.IsSet() to be true")
	}
}

func TestLoad_EnvVars(t *testing.T) {
	os.Setenv("MYAPP_DB_HOST", "env.db.com")
	os.Setenv("MYAPP_DB_PORT", "5434")
	os.Setenv("MYAPP_DB_TIMEOUT", "1m")
	defer func() {
		os.Unsetenv("MYAPP_DB_HOST")
		os.Unsetenv("MYAPP_DB_PORT")
		os.Unsetenv("MYAPP_DB_TIMEOUT")
	}()

	cfg := newTestConfig()
	err := Load(&cfg, Options{EnvPrefix: "MYAPP"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DB.Host.Get() != "env.db.com" {
		t.Errorf("expected env.db.com, got %s", cfg.DB.Host.Get())
	}
	if cfg.DB.Port.Get() != 5434 {
		t.Errorf("expected 5434, got %d", cfg.DB.Port.Get())
	}
	if cfg.DB.Timeout.Get() != time.Minute {
		t.Errorf("expected 1m, got %v", cfg.DB.Timeout.Get())
	}
}

func TestLoad_YAMLFile(t *testing.T) {
	yamlContent := `
db:
  host: yaml.db.com
  port: 5435
  use_ssl: true
  timeout: 2m
server:
  addr: ":9090"
`
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := newTestConfig()
	err := Load(&cfg, Options{ConfigFile: configFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DB.Host.Get() != "yaml.db.com" {
		t.Errorf("expected yaml.db.com, got %s", cfg.DB.Host.Get())
	}
	if cfg.DB.Port.Get() != 5435 {
		t.Errorf("expected 5435, got %d", cfg.DB.Port.Get())
	}
	if cfg.DB.UseSSL.Get() != true {
		t.Errorf("expected true, got %v", cfg.DB.UseSSL.Get())
	}
	if cfg.DB.Timeout.Get() != 2*time.Minute {
		t.Errorf("expected 2m, got %v", cfg.DB.Timeout.Get())
	}
	if cfg.Server.Addr.Get() != ":9090" {
		t.Errorf("expected :9090, got %s", cfg.Server.Addr.Get())
	}
}

func TestLoad_SourcePriority(t *testing.T) {
	// set up all three sources with different values
	yamlContent := `
db:
  host: yaml.db.com
  port: 1111
`
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("TEST_DB_HOST", "env.db.com")
	os.Setenv("TEST_DB_PORT", "2222")
	defer func() {
		os.Unsetenv("TEST_DB_HOST")
		os.Unsetenv("TEST_DB_PORT")
	}()

	args := []string{"--db.host=cli.db.com"}

	cfg := newTestConfig()
	err := Load(&cfg, Options{
		ConfigFile: configFile,
		EnvPrefix:  "TEST",
		Args:       args,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// CLI should win for host
	if cfg.DB.Host.Get() != "cli.db.com" {
		t.Errorf("expected cli.db.com (CLI priority), got %s", cfg.DB.Host.Get())
	}
	// ENV should win for port (no CLI arg)
	if cfg.DB.Port.Get() != 2222 {
		t.Errorf("expected 2222 (ENV priority), got %d", cfg.DB.Port.Get())
	}
}

func TestLoad_Validation(t *testing.T) {
	cfg := newTestConfig()
	args := []string{"--db.port=99999"} // out of range

	err := Load(&cfg, Options{Args: args})
	if err == nil {
		t.Fatal("expected validation error")
	}

	loadErr, ok := err.(*LoadError)
	if !ok {
		t.Fatalf("expected LoadError, got %T", err)
	}

	found := false
	for _, e := range loadErr.Errors {
		if _, ok := e.(*ValidationError); ok {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected ValidationError in LoadError")
	}
}

func TestLoad_Required(t *testing.T) {
	type configWithRequired struct {
		Name StringParam `cfg:"name"`
	}

	cfg := configWithRequired{
		Name: String().Required().Build(),
	}

	err := Load(&cfg, Options{})
	if err == nil {
		t.Fatal("expected required error")
	}

	loadErr, ok := err.(*LoadError)
	if !ok {
		t.Fatalf("expected LoadError, got %T", err)
	}

	found := false
	for _, e := range loadErr.Errors {
		if _, ok := e.(*RequiredError); ok {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected RequiredError in LoadError")
	}
}

func TestLoad_RequiredWithDefault(t *testing.T) {
	type configWithRequired struct {
		Name StringParam `cfg:"name"`
	}

	cfg := configWithRequired{
		Name: String().Required().Default("default-name").Build(),
	}

	err := Load(&cfg, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Name.Get() != "default-name" {
		t.Errorf("expected default-name, got %s", cfg.Name.Get())
	}
}

func TestLoad_ParseError(t *testing.T) {
	cfg := newTestConfig()
	args := []string{"--db.port=not-a-number"}

	err := Load(&cfg, Options{Args: args})
	if err == nil {
		t.Fatal("expected parse error")
	}

	loadErr, ok := err.(*LoadError)
	if !ok {
		t.Fatalf("expected LoadError, got %T", err)
	}

	found := false
	for _, e := range loadErr.Errors {
		if _, ok := e.(*ParseError); ok {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected ParseError in LoadError")
	}
}

func TestLoad_ListParams(t *testing.T) {
	type listConfig struct {
		Tags  StringListParam `cfg:"tags"`
		Ports IntListParam    `cfg:"ports"`
	}

	cfg := listConfig{
		Tags:  StringList().Default([]string{}).Build(),
		Ports: IntList().Default([]int{}).Build(),
	}

	args := []string{
		"--tags=alpha,beta,gamma",
		"--ports=8080,9090,3000",
	}

	err := Load(&cfg, Options{Args: args})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTags := []string{"alpha", "beta", "gamma"}
	if len(cfg.Tags.Get()) != len(expectedTags) {
		t.Errorf("expected %v, got %v", expectedTags, cfg.Tags.Get())
	}
	for i, tag := range cfg.Tags.Get() {
		if tag != expectedTags[i] {
			t.Errorf("expected tag %s, got %s", expectedTags[i], tag)
		}
	}

	expectedPorts := []int{8080, 9090, 3000}
	if len(cfg.Ports.Get()) != len(expectedPorts) {
		t.Errorf("expected %v, got %v", expectedPorts, cfg.Ports.Get())
	}
	for i, port := range cfg.Ports.Get() {
		if port != expectedPorts[i] {
			t.Errorf("expected port %d, got %d", expectedPorts[i], port)
		}
	}
}

func TestLoad_ListFromEnv(t *testing.T) {
	type listConfig struct {
		Tags StringListParam `cfg:"tags"`
	}

	os.Setenv("APP_TAGS", "one;two;three")
	defer os.Unsetenv("APP_TAGS")

	cfg := listConfig{
		Tags: StringList().Default([]string{}).Build(),
	}

	err := Load(&cfg, Options{EnvPrefix: "APP", ListSeparator: ";"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"one", "two", "three"}
	if len(cfg.Tags.Get()) != len(expected) {
		t.Errorf("expected %v, got %v", expected, cfg.Tags.Get())
	}
}

func TestLoad_ListFromYAML(t *testing.T) {
	type listConfig struct {
		Tags  StringListParam `cfg:"tags"`
		Ports IntListParam    `cfg:"ports"`
	}

	yamlContent := `
tags:
  - alpha
  - beta
  - gamma
ports:
  - 8080
  - 9090
`
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := listConfig{
		Tags:  StringList().Default([]string{}).Build(),
		Ports: IntList().Default([]int{}).Build(),
	}

	err := Load(&cfg, Options{ConfigFile: configFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTags := []string{"alpha", "beta", "gamma"}
	if len(cfg.Tags.Get()) != len(expectedTags) {
		t.Errorf("expected %v, got %v", expectedTags, cfg.Tags.Get())
	}

	expectedPorts := []int{8080, 9090}
	if len(cfg.Ports.Get()) != len(expectedPorts) {
		t.Errorf("expected %v, got %v", expectedPorts, cfg.Ports.Get())
	}
}

func TestLoad_AggregatedErrors(t *testing.T) {
	type multiErrorConfig struct {
		Name StringParam `cfg:"name"`
		Age  IntParam    `cfg:"age"`
		Port IntParam    `cfg:"port"`
	}

	cfg := multiErrorConfig{
		Name: String().Required().Build(),
		Age:  Int().Required().Build(),
		Port: Int().Default(8080).Validate(Range(1, 65535)).Build(),
	}

	args := []string{"--port=99999"} // validation error

	err := Load(&cfg, Options{Args: args})
	if err == nil {
		t.Fatal("expected errors")
	}

	loadErr, ok := err.(*LoadError)
	if !ok {
		t.Fatalf("expected LoadError, got %T", err)
	}

	// should have at least 3 errors: 2 required + 1 validation
	if len(loadErr.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %v", len(loadErr.Errors), loadErr.Errors)
	}
}

func TestLoad_FloatParam(t *testing.T) {
	type floatConfig struct {
		Rate FloatParam `cfg:"rate"`
	}

	cfg := floatConfig{
		Rate: Float().Default(0.5).Validate(RangeFloat(0, 1)).Build(),
	}

	args := []string{"--rate=0.75"}
	err := Load(&cfg, Options{Args: args})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Rate.Get() != 0.75 {
		t.Errorf("expected 0.75, got %f", cfg.Rate.Get())
	}
}

func TestLoad_DurationParam(t *testing.T) {
	type durationConfig struct {
		Timeout DurationParam `cfg:"timeout"`
	}

	cfg := durationConfig{
		Timeout: Duration().Default(time.Second).Validate(RangeDuration(time.Millisecond, time.Minute)).Build(),
	}

	args := []string{"--timeout=5s"}
	err := Load(&cfg, Options{Args: args})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Timeout.Get() != 5*time.Second {
		t.Errorf("expected 5s, got %v", cfg.Timeout.Get())
	}
}

func TestLoad_MissingYAMLFile(t *testing.T) {
	cfg := newTestConfig()
	err := Load(&cfg, Options{ConfigFile: "/nonexistent/config.yaml"})
	if err != nil {
		t.Fatalf("should not error for missing file: %v", err)
	}
	// should use defaults
	if cfg.DB.Host.Get() != "localhost" {
		t.Errorf("expected localhost, got %s", cfg.DB.Host.Get())
	}
}

func TestFindConfigFile(t *testing.T) {
	// create temp files
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "config1.yaml")
	file2 := filepath.Join(tmpDir, "config2.yaml")

	// only create file2
	if err := os.WriteFile(file2, []byte("test: value"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("finds first existing file", func(t *testing.T) {
		result := FindConfigFile([]string{file1, file2})
		if result != file2 {
			t.Errorf("expected %s, got %s", file2, result)
		}
	})

	t.Run("returns empty if none found", func(t *testing.T) {
		result := FindConfigFile([]string{"/nonexistent/a.yaml", "/nonexistent/b.yaml"})
		if result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	t.Run("returns empty for empty paths", func(t *testing.T) {
		result := FindConfigFile([]string{})
		if result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})
}

func TestDefaultConfigPaths(t *testing.T) {
	paths := DefaultConfigPaths("myapp")

	if len(paths) != 3 {
		t.Fatalf("expected 3 paths, got %d", len(paths))
	}

	// should contain local, xdg config path and /etc path
	foundLocal := false
	foundXdg := false
	foundEtc := false
	for _, p := range paths {
		if p == "config.yaml" {
			foundLocal = true
		} else if filepath.Base(filepath.Dir(p)) == "myapp" && filepath.Base(p) == "config.yaml" {
			if strings.HasPrefix(p, "/etc") {
				foundEtc = true
			} else {
				foundXdg = true
			}
		}
	}

	if !foundLocal {
		t.Error("expected local config path (config.yaml)")
	}
	if !foundXdg {
		t.Error("expected XDG config path")
	}
	if !foundEtc {
		t.Error("expected /etc config path")
	}
}

func TestDefaultConfigPaths_WithXDGConfigHome(t *testing.T) {
	customXDG := "/tmp/custom-xdg-config"
	os.Setenv("XDG_CONFIG_HOME", customXDG)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	paths := DefaultConfigPaths("myapp")
	expectedXDGPath := filepath.Join(customXDG, "myapp", "config.yaml")

	found := false
	for _, p := range paths {
		if p == expectedXDGPath {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected XDG path %q in paths %v", expectedXDGPath, paths)
	}
}

func TestLoad_ConfigPaths(t *testing.T) {
	// create temp config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	yamlContent := `
db:
  host: from-config-paths
`
	if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := newTestConfig()
	err := Load(&cfg, Options{
		ConfigPaths: []string{"/nonexistent/config.yaml", configFile},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DB.Host.Get() != "from-config-paths" {
		t.Errorf("expected from-config-paths, got %s", cfg.DB.Host.Get())
	}
}

func TestLoad_ConfigFileOverridesConfigPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// create two config files
	pathsConfig := filepath.Join(tmpDir, "paths.yaml")
	explicitConfig := filepath.Join(tmpDir, "explicit.yaml")

	if err := os.WriteFile(pathsConfig, []byte("db:\n  host: from-paths"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(explicitConfig, []byte("db:\n  host: from-explicit"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := newTestConfig()
	err := Load(&cfg, Options{
		ConfigFile:  explicitConfig,
		ConfigPaths: []string{pathsConfig},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ConfigFile should take precedence
	if cfg.DB.Host.Get() != "from-explicit" {
		t.Errorf("expected from-explicit, got %s", cfg.DB.Host.Get())
	}
}

func TestLoad_ComplexCommandLine(t *testing.T) {
	type myCfgConfig struct {
		Prop  StringParam `cfg:"prop"`
		Prop2 StringParam `cfg:"prop2"`
		Prop3 StringParam `cfg:"prop3"`
	}

	type complexConfig struct {
		MyCfg myCfgConfig `cfg:"my.cfg"`
	}

	args := []string{
		"-v",
		"--opt1", "optarg1",
		"--opt2",
		"mycmd",
		"-f",
		"mysubcmd",
		"--unrelated=something",
		"-fAsL",
		"-v",
		"--my.cfg.prop=thisisconfigprop",
		"--unrelated2",
		"--unrelated3", "something3",
		"--my.cfg.prop2=thisisconfigprop2",
		"--my.cfg.prop3", "space separated value",
	}

	cfg := complexConfig{
		MyCfg: myCfgConfig{
			Prop:  String().Default("default1").Build(),
			Prop2: String().Default("default2").Build(),
			Prop3: String().Default("default3").Build(),
		},
	}

	err := Load(&cfg, Options{Args: args})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// config properties should be correctly extracted
	if cfg.MyCfg.Prop.Get() != "thisisconfigprop" {
		t.Errorf("expected thisisconfigprop, got %s", cfg.MyCfg.Prop.Get())
	}
	if cfg.MyCfg.Prop2.Get() != "thisisconfigprop2" {
		t.Errorf("expected thisisconfigprop2, got %s", cfg.MyCfg.Prop2.Get())
	}
	if cfg.MyCfg.Prop3.Get() != "space separated value" {
		t.Errorf("expected 'space separated value', got %s", cfg.MyCfg.Prop3.Get())
	}
}

func TestLoad_ComplexCommandLineWithAllSources(t *testing.T) {
	type appConfig struct {
		Name    StringParam `cfg:"name"`
		Version StringParam `cfg:"version"`
		Debug   BoolParam   `cfg:"debug"`
	}

	type serverConfig struct {
		Host StringParam `cfg:"host"`
		Port IntParam    `cfg:"port"`
	}

	type complexConfig struct {
		App    appConfig    `cfg:"app"`
		Server serverConfig `cfg:"server"`
	}

	// YAML file with base config
	yamlContent := `
app:
  name: yaml-app
  version: "1.0.0"
  debug: false
server:
  host: yaml.host.com
  port: 3000
`
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	// ENV overrides some values
	os.Setenv("MYAPP_APP_VERSION", "2.0.0")
	os.Setenv("MYAPP_SERVER_PORT", "4000")
	defer func() {
		os.Unsetenv("MYAPP_APP_VERSION")
		os.Unsetenv("MYAPP_SERVER_PORT")
	}()

	// complex CLI with config props mixed in
	args := []string{
		"run",
		"--verbose",
		"-d",
		"--app.name=cli-app", // CLI overrides YAML
		"serve",
		"--workers", "4",
		"--app.debug=true", // CLI overrides YAML
		"-p", "production",
		"--server.port=5000", // CLI overrides ENV which overrides YAML
	}

	cfg := complexConfig{
		App: appConfig{
			Name:    String().Default("default-app").Build(),
			Version: String().Default("0.0.0").Build(),
			Debug:   Bool().Default(false).Build(),
		},
		Server: serverConfig{
			Host: String().Default("localhost").Build(),
			Port: Int().Default(8080).Build(),
		},
	}

	err := Load(&cfg, Options{
		ConfigFile: configFile,
		EnvPrefix:  "MYAPP",
		Args:       args,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// CLI wins for app.name
	if cfg.App.Name.Get() != "cli-app" {
		t.Errorf("expected cli-app, got %s", cfg.App.Name.Get())
	}

	// ENV wins for app.version (no CLI override)
	if cfg.App.Version.Get() != "2.0.0" {
		t.Errorf("expected 2.0.0, got %s", cfg.App.Version.Get())
	}

	// CLI wins for app.debug
	if cfg.App.Debug.Get() != true {
		t.Errorf("expected true, got %v", cfg.App.Debug.Get())
	}

	// YAML wins for server.host (no CLI or ENV override)
	if cfg.Server.Host.Get() != "yaml.host.com" {
		t.Errorf("expected yaml.host.com, got %s", cfg.Server.Host.Get())
	}

	// CLI wins for server.port
	if cfg.Server.Port.Get() != 5000 {
		t.Errorf("expected 5000, got %d", cfg.Server.Port.Get())
	}
}

func TestLoad_EdgeCasesInArgs(t *testing.T) {
	type edgeConfig struct {
		Empty      StringParam `cfg:"empty"`
		WithEquals StringParam `cfg:"with.equals"`
		WithSpaces StringParam `cfg:"with.spaces"`
		Numeric    StringParam `cfg:"numeric"`
	}

	tests := []struct {
		name     string
		args     []string
		expected map[string]string
	}{
		{
			name: "empty value with equals",
			args: []string{"--empty="},
			expected: map[string]string{
				"empty": "",
			},
		},
		{
			name: "value containing equals sign",
			args: []string{"--with.equals=key=value=extra"},
			expected: map[string]string{
				"with.equals": "key=value=extra",
			},
		},
		{
			name: "value with spaces using equals",
			args: []string{"--with.spaces=hello world with spaces"},
			expected: map[string]string{
				"with.spaces": "hello world with spaces",
			},
		},
		{
			name: "numeric string value",
			args: []string{"--numeric=12345"},
			expected: map[string]string{
				"numeric": "12345",
			},
		},
		{
			name: "mixed formats",
			args: []string{
				"--empty=",
				"--with.equals", "a=b=c",
				"--with.spaces=spaced value",
				"--numeric", "999",
			},
			expected: map[string]string{
				"empty":       "",
				"with.equals": "a=b=c",
				"with.spaces": "spaced value",
				"numeric":     "999",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := edgeConfig{
				Empty:      String().Default("not-empty").Build(),
				WithEquals: String().Default("default").Build(),
				WithSpaces: String().Default("default").Build(),
				Numeric:    String().Default("default").Build(),
			}

			err := Load(&cfg, Options{Args: tt.args})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if v, ok := tt.expected["empty"]; ok {
				if cfg.Empty.Get() != v {
					t.Errorf("empty: expected %q, got %q", v, cfg.Empty.Get())
				}
			}
			if v, ok := tt.expected["with.equals"]; ok {
				if cfg.WithEquals.Get() != v {
					t.Errorf("with.equals: expected %q, got %q", v, cfg.WithEquals.Get())
				}
			}
			if v, ok := tt.expected["with.spaces"]; ok {
				if cfg.WithSpaces.Get() != v {
					t.Errorf("with.spaces: expected %q, got %q", v, cfg.WithSpaces.Get())
				}
			}
			if v, ok := tt.expected["numeric"]; ok {
				if cfg.Numeric.Get() != v {
					t.Errorf("numeric: expected %q, got %q", v, cfg.Numeric.Get())
				}
			}
		})
	}
}

func TestLoad_BoolListParam(t *testing.T) {
	type boolListConfig struct {
		Flags BoolListParam `cfg:"flags"`
	}

	t.Run("CLI", func(t *testing.T) {
		cfg := boolListConfig{
			Flags: BoolList().Default([]bool{}).Build(),
		}
		err := Load(&cfg, Options{Args: []string{"--flags=true,false,true"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []bool{true, false, true}
		got := cfg.Flags.Get()
		if len(got) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("index %d: expected %v, got %v", i, expected[i], got[i])
			}
		}
	})

	t.Run("YAML", func(t *testing.T) {
		yamlContent := "flags:\n  - true\n  - false\n"
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := boolListConfig{
			Flags: BoolList().Default([]bool{}).Build(),
		}
		err := Load(&cfg, Options{ConfigFile: configFile})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []bool{true, false}
		got := cfg.Flags.Get()
		if len(got) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("index %d: expected %v, got %v", i, expected[i], got[i])
			}
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		cfg := boolListConfig{
			Flags: BoolList().Default([]bool{}).Build(),
		}
		err := Load(&cfg, Options{Args: []string{"--flags=true,notbool"}})
		if err == nil {
			t.Fatal("expected parse error")
		}
		loadErr, ok := err.(*LoadError)
		if !ok {
			t.Fatalf("expected LoadError, got %T", err)
		}
		found := false
		for _, e := range loadErr.Errors {
			if _, ok := e.(*ParseError); ok {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected ParseError in LoadError")
		}
	})
}

func TestLoad_FloatListParam(t *testing.T) {
	type floatListConfig struct {
		Rates FloatListParam `cfg:"rates"`
	}

	t.Run("CLI", func(t *testing.T) {
		cfg := floatListConfig{
			Rates: FloatList().Default([]float64{}).Build(),
		}
		err := Load(&cfg, Options{Args: []string{"--rates=1.5,2.7,3.14"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []float64{1.5, 2.7, 3.14}
		got := cfg.Rates.Get()
		if len(got) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("index %d: expected %v, got %v", i, expected[i], got[i])
			}
		}
	})

	t.Run("YAML", func(t *testing.T) {
		yamlContent := "rates:\n  - 1.5\n  - 2.7\n"
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := floatListConfig{
			Rates: FloatList().Default([]float64{}).Build(),
		}
		err := Load(&cfg, Options{ConfigFile: configFile})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []float64{1.5, 2.7}
		got := cfg.Rates.Get()
		if len(got) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("index %d: expected %v, got %v", i, expected[i], got[i])
			}
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		cfg := floatListConfig{
			Rates: FloatList().Default([]float64{}).Build(),
		}
		err := Load(&cfg, Options{Args: []string{"--rates=1.5,abc"}})
		if err == nil {
			t.Fatal("expected parse error")
		}
		loadErr, ok := err.(*LoadError)
		if !ok {
			t.Fatalf("expected LoadError, got %T", err)
		}
		found := false
		for _, e := range loadErr.Errors {
			if _, ok := e.(*ParseError); ok {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected ParseError in LoadError")
		}
	})
}

func TestLoad_DurationListParam(t *testing.T) {
	type durationListConfig struct {
		Timeouts DurationListParam `cfg:"timeouts"`
	}

	t.Run("CLI", func(t *testing.T) {
		cfg := durationListConfig{
			Timeouts: DurationList().Default([]time.Duration{}).Build(),
		}
		err := Load(&cfg, Options{Args: []string{"--timeouts=1s,2m,500ms"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []time.Duration{time.Second, 2 * time.Minute, 500 * time.Millisecond}
		got := cfg.Timeouts.Get()
		if len(got) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("index %d: expected %v, got %v", i, expected[i], got[i])
			}
		}
	})

	t.Run("YAML", func(t *testing.T) {
		yamlContent := "timeouts:\n  - \"1s\"\n  - \"2m\"\n"
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := durationListConfig{
			Timeouts: DurationList().Default([]time.Duration{}).Build(),
		}
		err := Load(&cfg, Options{ConfigFile: configFile})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []time.Duration{time.Second, 2 * time.Minute}
		got := cfg.Timeouts.Get()
		if len(got) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("index %d: expected %v, got %v", i, expected[i], got[i])
			}
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		cfg := durationListConfig{
			Timeouts: DurationList().Default([]time.Duration{}).Build(),
		}
		err := Load(&cfg, Options{Args: []string{"--timeouts=1s,invalid"}})
		if err == nil {
			t.Fatal("expected parse error")
		}
		loadErr, ok := err.(*LoadError)
		if !ok {
			t.Fatalf("expected LoadError, got %T", err)
		}
		found := false
		for _, e := range loadErr.Errors {
			if _, ok := e.(*ParseError); ok {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected ParseError in LoadError")
		}
	})
}
