package terraform

// TODO: all of this should probably be put in format.go

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// CommandArguments is meant to describe how a terraform or terragrunt command's arguments are converted to the command
// string representation.
type CommandArguments interface {
	// UnmarshalArgs should be implemented by structs whose command argument fields are annotated by `tfarg`
	MarshalTfArgs() (string, error)
}

// commandArgField is a struct that holds the information for each field of a CommandArguments struct
// This is needed because of how reflect works - we need to use reflect.TypeOf to get the tag values and
// field types, and then reflect.ValueOf to get the values of the struct passed in.
type commandArgField struct {
	name      string       // the name of the field, so that we can lookup the value of a field
	tagValue  string       // the value of the `tfarg` tag; the corresponding command line argument for a field
	kind      reflect.Kind // the type of the field. This determines how a value should be converted into the command string
	omittrue  bool         // whether to exclude arguments whose value is true; only good for bool args that are true by default, like -refresh
	omitempty bool         // whether to exclude arguments with a zero-value field when generating the string
}

var (
	ErrInvalidTfArgMap      = errors.New("invalid map type for tf argument")
	ErrInvalidTfArgSlice    = errors.New("invalid slice type for tf argument")
	ErrUnsupportedFieldType = errors.New("the type of this field is not supported by the default MarshalTfArgs function; consider writing one for this struct")
)

// a function that parses structs with the tfarg tag and, based on the content of the tagged field and the tag, we generate a command line
// TODO: consider returning []string instead
func MarshalTfArgs(argStruct any) (string, error) {
	// TODO: allow commands to define their own struct that can unmarshal itself to a "command line string"
	if args, hasCustomUnmarshalFunc := argStruct.(CommandArguments); hasCustomUnmarshalFunc {
		return args.MarshalTfArgs()
	}

	// go through each field of `argStruct`, look if it has the tfarg struct tag, parse it based on type(f)
	// limited set of supported types: these are already defined in format.go so we can just use those functions
	fields := reflect.TypeOf(argStruct).Elem()
	taggedFields := []*commandArgField{}
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)

		// "options" for how the field is encoded to a command string are supported, like "omitempty" for JSON encoding,
		// can be appended to tags and separated by commas
		tag := field.Tag.Get("tfarg")
		if len(tag) == 0 {
			continue
		}
		fullTagString := strings.Split(tag, ",")
		tagValue := fullTagString[0]
		caf := &commandArgField{
			name:     field.Name,
			kind:     field.Type.Kind(),
			tagValue: tagValue,
		}
		if len(fullTagString) > 1 {
			options := fullTagString[1:]
			for _, o := range options {
				if o == "omittrue" {
					caf.omittrue = true
				}
				if o == "omitempty" {
					caf.omitempty = true
				}
				// other supported options go here
			}
		}
		taggedFields = append(taggedFields, caf)
	}

	var commandString []string
	var err error
	// these "field" names are a nightmare
	fieldValues := reflect.ValueOf(argStruct).Elem()
	// TODO: is this a good idea?
	// defer func() {
	// 	if panicked := recover(); panicked != nil {
	// 		err = errors.WithMessage(ErrCannotUnmarshalTfArgs, panicked)
	// 	}
	// }()

	// iterate through the fields that had tags that were set and generate a command string
	for _, taggedField := range taggedFields {
		fieldValue := fieldValues.FieldByName(taggedField.name)

		switch taggedField.kind {
		// if the field is a bool type, we expect the user to have provided the flag name in the tfarg tag
		// if fieldValue.Bool() false, use the tag value =false; otherwise just use the tag value
		case reflect.Bool:
			b := fieldValue.Bool()
			if b {
				if !taggedField.omittrue {
					commandString = append(commandString, taggedField.tagValue)
				}
			} else {
				if !taggedField.omitempty {
					commandString = append(commandString, fmt.Sprintf("%s=false", taggedField.tagValue))
				}
			}
		case reflect.Map:
			// TODO: there's probably a better way to do this, I was just being lazy and wanted to use formatTerraformArgs()
			if m, ok := fieldValue.Interface().(map[string]any); ok {
				if len(m) > 0 {
					commandString = append(commandString, formatTerraformArgs(m, taggedField.tagValue, false)...)
				}
			} else {
				return "", ErrInvalidTfArgMap
			}
		case reflect.Slice:
			if s, ok := fieldValue.Interface().([]string); ok {
				if len(s) > 0 {
					commandString = append(commandString, FormatTerraformArgs(taggedField.tagValue, s)...)
				}
			} else {
				return "", errors.WithMessage(ErrInvalidTfArgSlice, "should have been a slice of strings")
			}
		case reflect.String:
			if argValue := fieldValue.String(); len(argValue) > 0 {
				commandString = append(commandString, fmt.Sprintf("%s=%s", taggedField.tagValue, argValue))
			}
		case reflect.Int:
			commandString = append(commandString, fmt.Sprintf("%s=%d", taggedField.tagValue, fieldValue.Int()))
		case reflect.Float32:
			commandString = append(commandString, fmt.Sprintf("%s=%f", taggedField.tagValue, fieldValue.Float()))
		// try to allow nested "TfArg" structs
		// TODO: this didn't work because those structs are a field of Options, and the are not themselves tagged - just their fields are
		// case reflect.Struct:
		// 	if argString, err := MarshalTfArgs(fieldValue.Interface()); err != nil {
		// 		return "", errors.Wrap(err, fmt.Sprintf("couldn't marshal field %s", taggedField.name))
		// 	} else {
		// 		commandString = append(commandString, argString)
		// 	}
		// case reflect.Pointer:
		// 	if argString, err := MarshalTfArgs(fieldValue.Elem().Interface()); err != nil {
		// 		return "", errors.Wrap(err, fmt.Sprintf("couldn't marshal field %s", taggedField.name))
		// 	} else {
		// 		commandString = append(commandString, argString)
		// 	}
		default:
			return "", errors.WithMessagef(ErrUnsupportedFieldType, "type: %s", taggedField.kind)
		}
	}

	return strings.Join(commandString, " "), err
}
