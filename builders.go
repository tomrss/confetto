package confetto

import "time"

// StringBuilder builds a StringParam.
type StringBuilder struct {
	p StringParam
}

// String returns a new StringBuilder.
func String() *StringBuilder {
	return &StringBuilder{}
}

func (b *StringBuilder) Default(v string) *StringBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *StringBuilder) Required() *StringBuilder {
	b.p.required = true
	return b
}

func (b *StringBuilder) Desc(d string) *StringBuilder {
	b.p.desc = d
	return b
}

func (b *StringBuilder) Validate(fn func(string) error) *StringBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *StringBuilder) Build() StringParam {
	return b.p
}

// IntBuilder builds an IntParam.
type IntBuilder struct {
	p IntParam
}

// Int returns a new IntBuilder.
func Int() *IntBuilder {
	return &IntBuilder{}
}

func (b *IntBuilder) Default(v int) *IntBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *IntBuilder) Required() *IntBuilder {
	b.p.required = true
	return b
}

func (b *IntBuilder) Desc(d string) *IntBuilder {
	b.p.desc = d
	return b
}

func (b *IntBuilder) Validate(fn func(int) error) *IntBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *IntBuilder) Build() IntParam {
	return b.p
}

// BoolBuilder builds a BoolParam.
type BoolBuilder struct {
	p BoolParam
}

// Bool returns a new BoolBuilder.
func Bool() *BoolBuilder {
	return &BoolBuilder{}
}

func (b *BoolBuilder) Default(v bool) *BoolBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *BoolBuilder) Required() *BoolBuilder {
	b.p.required = true
	return b
}

func (b *BoolBuilder) Desc(d string) *BoolBuilder {
	b.p.desc = d
	return b
}

func (b *BoolBuilder) Validate(fn func(bool) error) *BoolBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *BoolBuilder) Build() BoolParam {
	return b.p
}

// FloatBuilder builds a FloatParam.
type FloatBuilder struct {
	p FloatParam
}

// Float returns a new FloatBuilder.
func Float() *FloatBuilder {
	return &FloatBuilder{}
}

func (b *FloatBuilder) Default(v float64) *FloatBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *FloatBuilder) Required() *FloatBuilder {
	b.p.required = true
	return b
}

func (b *FloatBuilder) Desc(d string) *FloatBuilder {
	b.p.desc = d
	return b
}

func (b *FloatBuilder) Validate(fn func(float64) error) *FloatBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *FloatBuilder) Build() FloatParam {
	return b.p
}

// DurationBuilder builds a DurationParam.
type DurationBuilder struct {
	p DurationParam
}

// Duration returns a new DurationBuilder.
func Duration() *DurationBuilder {
	return &DurationBuilder{}
}

func (b *DurationBuilder) Default(v time.Duration) *DurationBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *DurationBuilder) Required() *DurationBuilder {
	b.p.required = true
	return b
}

func (b *DurationBuilder) Desc(d string) *DurationBuilder {
	b.p.desc = d
	return b
}

func (b *DurationBuilder) Validate(fn func(time.Duration) error) *DurationBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *DurationBuilder) Build() DurationParam {
	return b.p
}

// StringListBuilder builds a StringListParam.
type StringListBuilder struct {
	p StringListParam
}

// StringList returns a new StringListBuilder.
func StringList() *StringListBuilder {
	return &StringListBuilder{}
}

func (b *StringListBuilder) Default(v []string) *StringListBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *StringListBuilder) Required() *StringListBuilder {
	b.p.required = true
	return b
}

func (b *StringListBuilder) Desc(d string) *StringListBuilder {
	b.p.desc = d
	return b
}

func (b *StringListBuilder) Validate(fn func([]string) error) *StringListBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *StringListBuilder) Build() StringListParam {
	return b.p
}

// IntListBuilder builds an IntListParam.
type IntListBuilder struct {
	p IntListParam
}

// IntList returns a new IntListBuilder.
func IntList() *IntListBuilder {
	return &IntListBuilder{}
}

func (b *IntListBuilder) Default(v []int) *IntListBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *IntListBuilder) Required() *IntListBuilder {
	b.p.required = true
	return b
}

func (b *IntListBuilder) Desc(d string) *IntListBuilder {
	b.p.desc = d
	return b
}

func (b *IntListBuilder) Validate(fn func([]int) error) *IntListBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *IntListBuilder) Build() IntListParam {
	return b.p
}

// BoolListBuilder builds a BoolListParam.
type BoolListBuilder struct {
	p BoolListParam
}

// BoolList returns a new BoolListBuilder.
func BoolList() *BoolListBuilder {
	return &BoolListBuilder{}
}

func (b *BoolListBuilder) Default(v []bool) *BoolListBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *BoolListBuilder) Required() *BoolListBuilder {
	b.p.required = true
	return b
}

func (b *BoolListBuilder) Desc(d string) *BoolListBuilder {
	b.p.desc = d
	return b
}

func (b *BoolListBuilder) Validate(fn func([]bool) error) *BoolListBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *BoolListBuilder) Build() BoolListParam {
	return b.p
}

// FloatListBuilder builds a FloatListParam.
type FloatListBuilder struct {
	p FloatListParam
}

// FloatList returns a new FloatListBuilder.
func FloatList() *FloatListBuilder {
	return &FloatListBuilder{}
}

func (b *FloatListBuilder) Default(v []float64) *FloatListBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *FloatListBuilder) Required() *FloatListBuilder {
	b.p.required = true
	return b
}

func (b *FloatListBuilder) Desc(d string) *FloatListBuilder {
	b.p.desc = d
	return b
}

func (b *FloatListBuilder) Validate(fn func([]float64) error) *FloatListBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *FloatListBuilder) Build() FloatListParam {
	return b.p
}

// DurationListBuilder builds a DurationListParam.
type DurationListBuilder struct {
	p DurationListParam
}

// DurationList returns a new DurationListBuilder.
func DurationList() *DurationListBuilder {
	return &DurationListBuilder{}
}

func (b *DurationListBuilder) Default(v []time.Duration) *DurationListBuilder {
	b.p.defaultVal = v
	b.p.value = v
	b.p.hasDefVal = true
	return b
}

func (b *DurationListBuilder) Required() *DurationListBuilder {
	b.p.required = true
	return b
}

func (b *DurationListBuilder) Desc(d string) *DurationListBuilder {
	b.p.desc = d
	return b
}

func (b *DurationListBuilder) Validate(fn func([]time.Duration) error) *DurationListBuilder {
	b.p.validators = append(b.p.validators, fn)
	return b
}

func (b *DurationListBuilder) Build() DurationListParam {
	return b.p
}
