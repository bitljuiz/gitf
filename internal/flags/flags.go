package flags

import "fmt"

type Flag struct {
	Name          string
	ShorthandName string
	Use           string
	Value         Value
}

func (f *Flag) GetString() (string, error) {
	switch val := f.Value.(type) {
	case *StringValue:
		return val.Value(), nil
	default:
		return "", fmt.Errorf("cannot get string value of %s flag. The value should be %s",
			f.Name, f.Value.Type())
	}
}

func (f *Flag) GetBool() (bool, error) {
	switch val := f.Value.(type) {
	case *BoolValue:
		return val.Value(), nil
	default:
		return false, fmt.Errorf("cannot get string value of %s flag. The value should be %s",
			f.Name, f.Value.Type())
	}
}

func (f *Flag) GetStringSlice() ([]string, error) {
	switch val := f.Value.(type) {
	case *StringSliceValue:
		return val.Value(), nil
	default:
		return nil, fmt.Errorf("cannot get string value of %s flag. The value should be %s",
			f.Name, f.Value.Type())
	}
}

const (
	Repository = iota
	Revision
	OrderBy
	UseCommitter
	Format
	Extensions
	Languages
	Exclude
	RestrictTo
	OnlyProgramming
	flagCount

	BooleanFlag
	StringFlag
	StringSliceFlag
)

var (
	flagToName = map[int]string{
		Repository:      "repository",
		Revision:        "revision",
		OrderBy:         "order-by",
		UseCommitter:    "use-committer",
		Format:          "format",
		Extensions:      "extensions",
		Languages:       "languages",
		Exclude:         "exclude",
		RestrictTo:      "restrict-to",
		OnlyProgramming: "only-programming",
	}
	flagToShorthandName = map[int]string{
		Repository:      "p",
		Revision:        "v",
		OrderBy:         "o",
		UseCommitter:    "u",
		Format:          "f",
		Extensions:      "x",
		Languages:       "l",
		Exclude:         "c",
		RestrictTo:      "r",
		OnlyProgramming: "g",
	}
	flagToUsage = map[int]string{
		Repository:   "Returns path to Git repository. Current directory by default",
		Revision:     "Returns pointer to commit. HEAD by default",
		OrderBy:      "Sets sorting key for results. Could be: lines (default), commits, files",
		UseCommitter: "Changes author (default) by committer in commit",
		Format:       "Sets output format. Could be tabular (default), csv, json, json-lines",
		Extensions: "Reduces the number of extensions processed. " +
			"A comma separated list of extensions is accepted as input",
		Languages:       "Reduces the number of languages processed. A comma separated list of languages is accepted as input",
		Exclude:         "Reduces the number of files processed by Glob patterns, like 'foo/*, bar/*'",
		RestrictTo:      "Reduces the number of files processed, excluding those who don't satisfy any of Glob patterns",
		OnlyProgramming: "Returns information only about files written in programming language",
	}
	flagToValueType = map[int]int{
		Repository:      StringFlag,
		Revision:        StringFlag,
		OrderBy:         StringFlag,
		UseCommitter:    BooleanFlag,
		Format:          StringFlag,
		Extensions:      StringSliceFlag,
		Languages:       StringSliceFlag,
		Exclude:         StringSliceFlag,
		RestrictTo:      StringSliceFlag,
		OnlyProgramming: BooleanFlag,
	}
	flagToDefaultValue = map[int]interface{}{
		Repository:      ".",
		Revision:        "HEAD",
		OrderBy:         "lines",
		UseCommitter:    false,
		Format:          "tabular",
		Extensions:      []string{},
		Languages:       []string{},
		Exclude:         []string{},
		RestrictTo:      []string{},
		OnlyProgramming: false,
	}
)

func createValue(valueType, flagType int) Value {
	switch valueType {
	case BooleanFlag:
		return newBoolValue(flagToDefaultValue[flagType].(bool))
	case StringFlag:
		return newStringValue(flagToDefaultValue[flagType].(string))
	case StringSliceFlag:
		return newStringSliceValue(flagToDefaultValue[flagType].([]string))
	default:
		panic("unhandled default case")
	}
}

func Create() map[int]*Flag {
	flags := make(map[int]*Flag, flagCount)
	for flagIota := range flagCount {
		flags[flagIota] = &Flag{
			Name:          flagToName[flagIota],
			ShorthandName: flagToShorthandName[flagIota],
			Use:           flagToUsage[flagIota],
			Value:         createValue(flagToValueType[flagIota], flagIota),
		}
	}
	return flags
}
