package configo

import (
	"github.com/mitchellh/mapstructure"
	"path"
	"reflect"
	"strconv"
	"strings"
)
import "github.com/spf13/viper"

type Format string

const FormatYaml = "yaml"
const FormatJson = "json"
const FormatEnv = "env"

// Option is struct with parameters for parse func
type Option struct {
	// Path config file
	// example:
	// /usr/share/something/config.json
	Path string

	// Key is parent name for config
	// for example if you need parse config into struct:
	// type MyConf struct {
	//   Address string `conf: "address"`
	//	 Port 	 int `conf: "port"`
	// }
	// from file:
	// {
	//	  "http": {
	//      "address": "localhost",
	//	    "port": 23
	//	  }
	// }
	// you need write Key as `http`
	Key string

	// EnvPrefix need for prefix env
	// if you need parse config into struct:
	// type MyConf struct {
	//   Address string `conf: "address"`
	//	 Port 	 int `conf: "port"`
	// }
	// from env with names:
	// MY_ENV_ADDRESS and MY_ENV_PORT
	// you need set EnvPrefix = "MY_ENV" or EnvPrefix = "my_env"
	EnvPrefix string

	// EnvInclude if need include env vars from system
	EnvInclude bool
}

func MustParse[T any](option Option) (t T) {
	t, err := Parse[T](option)
	if err != nil {
		panic(err)
	}

	return
}

// Parse func for parse config
// example:
//
//	type MyConf struct {
//	  Address string `conf: "address" default:"localhost"`
//		 Port 	 int 	`conf: "port" default:"80"`
//	}
//
// parsedConfig := config.Parse[MyConf](config.Option{})
func Parse[T any](option Option) (t T, err error) {
	raw := make(map[string]any)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:              "configo",
		IgnoreUntaggedFields: true,
		Result:               &raw,
	})

	if err != nil {
		return
	}

	provideDefaultTag(&t)

	// parse default variable from struct
	err = decoder.Decode(&t)

	if len(option.Key) != 0 {
		holder := make(map[string]any)
		holder[option.Key] = raw
		raw = holder
	}

	if err != nil {
		return
	}

	if len(option.Path) != 0 {
		viper.SetConfigFile(option.Path)

		ext := path.Ext(option.Path)

		if len(ext) > 1 {
			ext = ext[1:]
		}

		viper.SetConfigType(ext)
	}

	if option.EnvInclude && len(option.EnvPrefix) != 0 {
		viper.SetEnvPrefix(option.EnvPrefix)
	}

	if option.EnvInclude {
		viper.AutomaticEnv()
		viper.AllowEmptyEnv(true)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	err = viper.MergeConfigMap(raw)

	if err != nil {
		return
	}

	if len(option.Key) != 0 {
		confHolderOut := make(map[string]T)
		err = viper.Unmarshal(&confHolderOut)

		if err != nil {
			return
		}

		t = confHolderOut[option.Key]

		return
	} else {
		err = viper.Unmarshal(&t)

		if err != nil {
			return
		}

		return
	}
}

func provideDefaultTag(entity any) {
	if reflect.TypeOf(entity).Kind() != reflect.Ptr {
		return
	}

	provideDefaultTagInternal(reflect.ValueOf(entity))
}

func provideDefaultTagInternal(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		if configTag := t.Field(i).Tag.Get("configo"); configTag != "-" && len(configTag) > 0 {
			if v.Field(i).Kind() == reflect.Struct {
				provideDefaultTagInternal(v.Field(i))
			} else if v.Field(i).Kind() == reflect.Pointer {
				if v.Field(i).IsNil() {
					v.Field(i).Set(reflect.New(v.Field(i).Type().Elem()))
					if v.Field(i).Elem().Kind() == reflect.Struct {
						provideDefaultTagInternal(v.Field(i))
					} else {
						if defaultVal := t.Field(i).Tag.Get("default"); defaultVal != "-" && len(defaultVal) > 0 {
							setValue(v.Field(i).Elem(), defaultVal)
						}
					}
				}
			} else {
				if defaultVal := t.Field(i).Tag.Get("default"); defaultVal != "-" && len(defaultVal) > 0 {
					setValue(v.Field(i), defaultVal)
				}
			}
		}
	}
}

func setValue(entity reflect.Value, value string) {
	if !entity.CanSet() {
		return
	}

	switch entity.Kind() {
	case
		reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseInt(value, 10, 64); err == nil {
			entity.Set(reflect.ValueOf(val).Convert(entity.Type()))
		}
	case reflect.String:
		entity.Set(reflect.ValueOf(value).Convert(entity.Type()))
	case reflect.Bool:
		if value == "true" {
			entity.Set(reflect.ValueOf(true).Convert(entity.Type()))
		} else {
			entity.Set(reflect.ValueOf(false).Convert(entity.Type()))
		}
	case reflect.Struct:
		setValue(entity, value)
	default:
		return
	}
}
