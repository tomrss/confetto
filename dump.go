package confetto

import (
	"strings"
)

const maskedValue = "****"

// Dump returns a string representation of all configuration parameters
// in the provided struct. Secret parameters are masked with "****".
// Parameters that are not set and have no default value are shown as "<not set>".
func Dump(cfg any) string {
	return dumpParams(collectParams(cfg, ""))
}

// Dump returns a string representation of all configuration parameters
// across all registered configs. Secret parameters are masked with "****".
// Parameters that are not set and have no default value are shown as "<not set>".
func (l *Loader) Dump() string {
	return dumpParams(l.collectAllParams())
}

func dumpParams(params []Param) string {
	var b strings.Builder
	for i, p := range params {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(p.key())
		b.WriteString(" = ")
		if p.isSecret() {
			b.WriteString(maskedValue)
		} else if !p.IsSet() && !p.hasDefault() {
			b.WriteString("<not set>")
		} else {
			b.WriteString(p.stringValue())
		}
	}
	return b.String()
}
