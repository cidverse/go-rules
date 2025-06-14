package expr

import (
	"os/exec"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

var (
	stringListType      = reflect.TypeOf([]string{})
	additionalFunctions = []cel.EnvOption{
		cel.Function(overloads.Contains,
			cel.Overload("string_contains_string",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.BoolType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					return types.Bool(strings.Contains(string(lhs.(types.String)), string(rhs.(types.String))))
				}),
			),
			cel.Overload("stringslice_contains_string",
				[]*cel.Type{cel.ListType(cel.StringType), cel.StringType},
				cel.BoolType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					list, err := lhs.ConvertToNative(stringListType)
					if err != nil {
						return types.NewErr(err.Error())
					}
					return types.Bool(slices.Contains(list.([]string), string(rhs.(types.String))))
				}),
			),
		),
		cel.Function("containsKey",
			cel.Overload("containsKey_map",
				[]*cel.Type{cel.MapType(cel.StringType, cel.StringType), cel.StringType},
				cel.StringType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					mapVal, err := lhs.ConvertToNative(reflect.TypeOf(map[string]string{}))
					if err != nil {
						return types.NewErr(err.Error())
					}
					if _, ok := mapVal.(map[string]string)[string(rhs.(types.String))]; ok {
						return types.Bool(true)
					}
					return types.Bool(false)
				}),
			),
		),
		cel.Function("getMapValue",
			cel.Overload("getMapValue_map",
				[]*cel.Type{cel.MapType(cel.StringType, cel.StringType), cel.StringType},
				cel.StringType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					mapVal, err := lhs.ConvertToNative(reflect.TypeOf(map[string]string{}))
					if err != nil {
						return types.NewErr(err.Error())
					}
					if value, ok := mapVal.(map[string]string)[string(rhs.(types.String))]; ok {
						return types.String(value)
					} else {
						return types.String("")
					}
				}),
			),
		),
		cel.Function("hasPrefix",
			cel.Overload("hasPrefix_string",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.BoolType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					return types.Bool(strings.HasPrefix(string(lhs.(types.String)), string(rhs.(types.String))))
				}),
			),
		),
		cel.Function("inPath",
			cel.Overload("inPath",
				[]*cel.Type{cel.StringType},
				cel.BoolType,
				cel.UnaryBinding(func(key ref.Val) ref.Val {
					_, err := exec.LookPath(string(key.(types.String)))
					if err == nil {
						return types.Bool(true)
					}

					return types.Bool(false)
				}),
			),
		),
		cel.Function("regex",
			cel.Overload("regex_string",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.BoolType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					input := string(lhs.(types.String))
					pattern := string(rhs.(types.String))
					matched, err := regexp.MatchString(pattern, input)
					if err != nil {
						return types.NewErr(err.Error())
					}

					return types.Bool(matched)
				}),
			),
		),
	}
)
