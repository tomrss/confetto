package confetto

import (
	"testing"
	"time"
)

func TestDump(t *testing.T) {
	type DB struct {
		Host     StringParam `cfg:"host"`
		Password StringParam `cfg:"password"`
		Name     StringParam `cfg:"name"`
	}
	type Config struct {
		Host string      `cfg:"host"`
		Port IntParam    `cfg:"port"`
		DB   DB          `cfg:"db"`
	}

	cfg := Config{
		Port: Int().Default(8080).Build(),
		DB: DB{
			Host:     String().Default("dbserver").Build(),
			Password: String().Secret().Required().Build(),
			Name:     String().Default("mydb").Build(),
		},
	}

	// simulate loading: set password value
	cfg.DB.Password.value = "s3cret"
	cfg.DB.Password.set = true

	// set keys via collectParams (same as Load does)
	collectParams(&cfg, "")

	got := Dump(&cfg)
	expected := "port = 8080\ndb.host = dbserver\ndb.password = ****\ndb.name = mydb"
	if got != expected {
		t.Errorf("Dump() =\n%s\nwant:\n%s", got, expected)
	}
}

func TestDumpNotSet(t *testing.T) {
	type Config struct {
		Host StringParam `cfg:"host"`
		Port IntParam    `cfg:"port"`
	}

	cfg := Config{
		Port: Int().Default(3000).Build(),
	}
	collectParams(&cfg, "")

	got := Dump(&cfg)
	expected := "host = <not set>\nport = 3000"
	if got != expected {
		t.Errorf("Dump() =\n%s\nwant:\n%s", got, expected)
	}
}

func TestDumpAllTypes(t *testing.T) {
	type Config struct {
		S  StringParam       `cfg:"s"`
		I  IntParam          `cfg:"i"`
		B  BoolParam         `cfg:"b"`
		F  FloatParam        `cfg:"f"`
		D  DurationParam     `cfg:"d"`
		SL StringListParam   `cfg:"sl"`
		IL IntListParam      `cfg:"il"`
		BL BoolListParam     `cfg:"bl"`
		FL FloatListParam    `cfg:"fl"`
		DL DurationListParam `cfg:"dl"`
	}

	cfg := Config{
		S:  String().Default("hello").Build(),
		I:  Int().Default(42).Build(),
		B:  Bool().Default(true).Build(),
		F:  Float().Default(3.14).Build(),
		D:  Duration().Default(5 * time.Second).Build(),
		SL: StringList().Default([]string{"a", "b"}).Build(),
		IL: IntList().Default([]int{1, 2}).Build(),
		BL: BoolList().Default([]bool{true, false}).Build(),
		FL: FloatList().Default([]float64{1.1, 2.2}).Build(),
		DL: DurationList().Default([]time.Duration{time.Second, time.Minute}).Build(),
	}
	collectParams(&cfg, "")

	got := Dump(&cfg)
	expected := "s = hello\ni = 42\nb = true\nf = 3.14\nd = 5s\nsl = [a b]\nil = [1 2]\nbl = [true false]\nfl = [1.1 2.2]\ndl = [1s 1m0s]"
	if got != expected {
		t.Errorf("Dump() =\n%s\nwant:\n%s", got, expected)
	}
}

func TestDumpSecretMaskedEvenIfNotSet(t *testing.T) {
	type Config struct {
		Token StringParam `cfg:"token"`
	}

	cfg := Config{
		Token: String().Secret().Build(),
	}
	collectParams(&cfg, "")

	got := Dump(&cfg)
	expected := "token = ****"
	if got != expected {
		t.Errorf("Dump() = %q, want %q", got, expected)
	}
}

func TestDumpEmpty(t *testing.T) {
	type Config struct{}
	cfg := Config{}
	got := Dump(&cfg)
	if got != "" {
		t.Errorf("Dump() = %q, want empty string", got)
	}
}
