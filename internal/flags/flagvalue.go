package flags

type Value interface {
	Type() string
}

type BoolValue struct {
	value        bool
	defaultValue bool
}

func newBoolValue(defaultValue bool) *BoolValue {
	b := BoolValue{
		defaultValue: defaultValue,
	}
	return &b
}

func (b *BoolValue) Type() string {
	return "bool"
}

func (b *BoolValue) Pointer() *bool { return &b.value }

func (b BoolValue) Value() bool { return b.value }

func (b BoolValue) DefaultValue() bool { return b.defaultValue }

type StringValue struct {
	value        string
	defaultValue string
}

func newStringValue(defaultValue string) *StringValue {
	s := StringValue{
		defaultValue: defaultValue,
	}
	return &s
}
func (sv *StringValue) Type() string { return "string" }

func (sv *StringValue) Pointer() *string {
	return &sv.value
}

func (sv StringValue) Value() string {
	return sv.value
}

func (sv StringValue) DefaultValue() string {
	return sv.defaultValue
}

type StringSliceValue struct {
	value        []string
	defaultValue []string
}

func newStringSliceValue(defaultValue []string) *StringSliceValue {
	ssv := StringSliceValue{
		defaultValue: defaultValue,
	}
	return &ssv
}

func (ssv *StringSliceValue) Type() string {
	return "stringSlice"
}

func (ssv *StringSliceValue) Pointer() *[]string {
	return &ssv.value
}

func (ssv StringSliceValue) Value() []string {
	return ssv.value
}

func (ssv StringSliceValue) DefaultValue() []string {
	return ssv.defaultValue
}
