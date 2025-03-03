package config

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// MustHydrate hydrates the given struct with values from the environment.
//
// If a field has a tag "env" then the value of the environment variable with the same
// name will be used. The "env" value will be prefixed with the given prefix from the
// value of the `env` tag on the parent struct.
//
// If a field has a tag "toml" then the value will be fetched from the relevant path from
// the global toml config file, provided by in the root of the project.
func MustHydrate(ctx context.Context, cfg interface{}) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer && v.Elem().Kind() != reflect.Struct {
		panic("cfg must be a pointer to a struct")
	}

	t := v.Elem().Type()
	v = v.Elem()

	hydrateStructFromEnv(ctx, v, t, "")
	// hydrateStructFromTOML(ctx, v, t, "")
}

func hydrateStructFromEnv(ctx context.Context, v reflect.Value, t reflect.Type, prefix string) {
	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)
		envTag := fieldType.Tag.Get("env")

		if envTag == "" {
			continue
		}

		if prefix != "" {
			envTag = prefix + "_" + envTag
		}

		switch field.Kind() {
		case reflect.Ptr:
			panic("pointers are not supported")
		case reflect.Struct:
			hydrateStructFromEnv(ctx, field, fieldType.Type, envTag)
		case reflect.Slice:
			panic("slices are not supported")
		default:
			envValue := os.Getenv(envTag)

			if envValue != "" {
				if field.CanSet() {
					switch field.Kind() {
					case reflect.String:
						field.SetString(envValue)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						value, err := strconv.ParseInt(envValue, 10, field.Type().Bits())
						if err != nil {
							panic(fmt.Sprintf("failed to parse %s as int: %v", envTag, err))
						}
						field.SetInt(value)
					case reflect.Bool:
						value, err := strconv.ParseBool(envValue)
						if err != nil {
							panic(fmt.Sprintf("failed to parse %s as bool: %v", envTag, err))
						}
						field.SetBool(value)
					case reflect.Float32, reflect.Float64:
						value, err := strconv.ParseFloat(envValue, field.Type().Bits())
						if err != nil {
							panic(fmt.Sprintf("failed to parse %s as float: %v", envTag, err))
						}
						field.SetFloat(value)
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						value, err := strconv.ParseUint(envValue, 10, field.Type().Bits())
						if err != nil {
							panic(fmt.Sprintf("failed to parse %s as uint: %v", envTag, err))
						}
						field.SetUint(value)
					default:
						panic(fmt.Sprintf("unsupported field type: %s", field.Kind()))
					}
				}
			}
		}
	}
}

func hydrateStructFromTOML(ctx context.Context, v reflect.Value, t reflect.Type, prefix string) {
	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)
		tomlTag := fieldType.Tag.Get("toml")

		if tomlTag == "" {
			continue
		}

		if prefix != "" {
			tomlTag = prefix + "_" + tomlTag
		}

		switch field.Kind() {
		case reflect.Ptr:
			panic("pointers are not supported")
		case reflect.Struct:
			hydrateStructFromEnv(ctx, field, fieldType.Type, tomlTag)
		case reflect.Slice:
			panic("slices are not supported")
		default:
			envValue := os.Getenv(tomlTag)

			if envValue != "" {
				if field.CanSet() {
					switch field.Kind() {
					case reflect.String:
						field.SetString(envValue)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						value, err := strconv.ParseInt(envValue, 10, field.Type().Bits())
						if err != nil {
							panic(fmt.Sprintf("failed to parse %s as int: %v", tomlTag, err))
						}
						field.SetInt(value)
					case reflect.Bool:
						value, err := strconv.ParseBool(envValue)
						if err != nil {
							panic(fmt.Sprintf("failed to parse %s as bool: %v", tomlTag, err))
						}
						field.SetBool(value)
					case reflect.Float32, reflect.Float64:
						value, err := strconv.ParseFloat(envValue, field.Type().Bits())
						if err != nil {
							panic(fmt.Sprintf("failed to parse %s as float: %v", tomlTag, err))
						}
						field.SetFloat(value)
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						value, err := strconv.ParseUint(envValue, 10, field.Type().Bits())
						if err != nil {
							panic(fmt.Sprintf("failed to parse %s as uint: %v", tomlTag, err))
						}
						field.SetUint(value)
					default:
						panic(fmt.Sprintf("unsupported field type: %s", field.Kind()))
					}
				}
			}
		}
	}
}
