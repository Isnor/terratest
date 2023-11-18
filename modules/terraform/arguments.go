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
	UnmarshalArgs() (string, error)
}

// ExampleCommandArgs are command argument examples of different types of arguments from the existing commands
type ExampleCommandArgs struct {
	BackendConfig map[string]any `tfarg:"-backend-config"`                      // The vars to pass to the terraform init command for extra configuration for the backend
	MigrateState  bool           `tfarg:"-migrate-state -force-copy,omitempty"` // Set the -migrate-state and -force-copy (suppress 'yes' answer prompt) flag to the terraform init command
	PluginDir     string         `tfarg:"-plugin-dir"`                          // The path of downloaded plugins to pass to the terraform init command (-plugin-dir)
	// reconfigure is also a bool type, but it's encoded slightly differently than MigrateState
	Reconfigure bool `tfarg:"-reconfigure,omitempty"` // Set the -reconfigure flag to the terraform init command
	// refresh is another odd bool case where the default is actually true, so when this is supplied as false we need to write `-refresh=false`
	Refresh bool `tfarg:"-refresh,omittrue"`
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

var ErrInvalidTfArgMap = errors.New("invalid map type for tf argument")
var ErrUnsupportedFieldType = errors.New("the type of this field is not supported by the default MarshalTfArgs function; consider writing one for this struct")

// a function that parses structs with the tfarg tag and, based on the content of the tagged field and the tag, we generate a command line
func MarshalTfArgs[T any](v *T) (string, error) {
	// TODO: allow commands to define their own struct that can unmarshal itself to a "command line string"
	// if args, hasCustomUnmarshalFunc := v.(CommandArguments); hasCustomUnmarshalFunc {
	// 	return args.UnmarshalArgs()
	// }

	// go through each field of `v`, look if it has the tfarg struct tag, parse it based on type(f)
	// limited set of supported types: these are already defined in format.go so we can just use those functions
	fields := reflect.TypeOf(*v)
	taggedFields := []*commandArgField{}
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)

		// "options" for how the field is encoded to a command string are supported, like "omitempty" for JSON encoding
		fullTagString := strings.Split(field.Tag.Get("tfarg"), ",")
		if len(fullTagString) == 0 {
			continue
		}
		tagValue := fullTagString[0]
		caf := &commandArgField{
			name:     field.Name,
			kind:     field.Type.Kind(),
			tagValue: tagValue,
		}
		if len(fullTagString) > 1 {
			options := fullTagString[1:]
			for _, o := range options {
				// other supported options go here
				if o == "omittrue" {
					caf.omittrue = true
				}
				if o == "omitempty" {
					caf.omitempty = true
				}
			}
		}
		taggedFields = append(taggedFields, caf)
	}

	var commandString []string
	var err error
	fieldValues := reflect.ValueOf(*v)
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
			if taggedField.omitempty && !fieldValue.Bool() {
				continue
			}
			if taggedField.omittrue && fieldValue.Bool() {
				continue
			}
			if fieldValue.Bool() {
				commandString = append(commandString, taggedField.tagValue)
			} else {
				commandString = append(commandString, fmt.Sprintf("%s=false", taggedField.tagValue))
			}
			// if !fieldValue.Bool() {
			// 	if !taggedField.omitempty {
			// 		commandString = append(commandString, fmt.Sprintf("%s=false", taggedField.tagValue))
			// 		break
			// 	}
			// }
			// //fieldValue = true
			// // don't include
			// if !taggedField.omittrue {
			// 	commandString = append(commandString, taggedField.tagValue)
			// }
		case reflect.Map:
			if m, ok := fieldValue.Interface().(map[string]any); ok {
				if len(m) > 0 {
					commandString = append(commandString, formatTerraformArgs(m, taggedField.tagValue, false)...)
				}
			} else {
				return "", ErrInvalidTfArgMap
			}
			// format map function
		case reflect.Slice:
			// format slice function
		case reflect.String:
			if argValue := fieldValue.String(); len(argValue) > 0 {
				commandString = append(commandString, fmt.Sprintf("%s=%s", taggedField.tagValue, argValue))
			}
		}
	}

	return strings.Join(commandString, " "), err
}
