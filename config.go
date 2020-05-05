package go_config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pefish/go-reflect"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const (
	YAML_TYPE = `yaml`
	JSON_TYPE = `json`
)

type ConfigClass struct {
	configs  map[string]interface{}
	loadType string
}

var Config = ConfigClass{
	configs: map[string]interface{}{},
	loadType: YAML_TYPE,
}

type Configuration struct {
	ConfigFilepath string
	SecretFilepath string
	ConfigEnvName  string
	SecretEnvName  string
}

func (configInstance *ConfigClass) MustLoadYamlConfig(config Configuration) {
	err := configInstance.LoadYamlConfig(config)
	if err != nil {
		panic(err)
	}
}

func (configInstance *ConfigClass) LoadYamlConfig(config Configuration) error {
	configInstance.loadType = YAML_TYPE

	configFile := ``
	configMap := make(map[string]interface{})
	if config.ConfigEnvName != `` || config.ConfigFilepath != `` {
		if config.ConfigEnvName != `` {
			configFile = os.Getenv(config.ConfigEnvName)
		} else if config.ConfigFilepath != `` {
			configFile = config.ConfigFilepath
		}
		if configFile != `` {
			bytes, err := ioutil.ReadFile(configFile)
			if err == nil {
				err = yaml.Unmarshal(bytes, &configMap)
				if err != nil {
					return err
				}
			}
		}
	}

	secretFile := ``
	secretMap := make(map[string]interface{})
	if config.SecretEnvName != `` || config.SecretFilepath != `` {
		if config.SecretEnvName != `` {
			secretFile = os.Getenv(config.SecretEnvName)
		} else if config.SecretFilepath != `` {
			secretFile = config.SecretFilepath
		}
		if secretFile != `` {
			bytes, err := ioutil.ReadFile(secretFile)
			if err == nil {
				err = yaml.Unmarshal(bytes, &secretMap)
				if err != nil {
					return err
				}
			}
		}
	}

	configInstance.configs = configMap
	for key, val := range secretMap {
		configInstance.configs[key] = val
	}
	return nil
}

func (configInstance *ConfigClass) MustLoadJsonConfig(config Configuration) {
	err := configInstance.LoadJsonConfig(config)
	if err != nil {
		panic(err)
	}
}

func (configInstance *ConfigClass) LoadJsonConfig(config Configuration) error {
	configInstance.loadType = JSON_TYPE
	configFile := ``
	configMap := make(map[string]interface{})
	if config.ConfigEnvName != `` || config.ConfigFilepath != `` {
		if config.ConfigEnvName != `` {
			configFile = os.Getenv(config.ConfigEnvName)
		} else if config.ConfigFilepath != `` {
			configFile = config.ConfigFilepath
		}
		if configFile != `` {
			bytes, err := ioutil.ReadFile(configFile)
			if err == nil {
				var result interface{}
				if err := json.Unmarshal(bytes, &result); err != nil {
					return err
				}
				configMap = result.(map[string]interface{})
			}
		}
	}

	secretFile := ``
	secretMap := make(map[string]interface{})
	if config.SecretEnvName != `` || config.SecretFilepath != `` {
		if config.SecretEnvName != `` {
			secretFile = os.Getenv(config.SecretEnvName)
		} else if config.SecretFilepath != `` {
			secretFile = config.SecretFilepath
		}
		if secretFile != `` {
			bytes, err := ioutil.ReadFile(secretFile)
			if err == nil {
				var result interface{}
				if err := json.Unmarshal(bytes, &result); err != nil {
					return err
				}
				secretMap = result.(map[string]interface{})
			}
		}
	}

	configInstance.configs = configMap
	for key, val := range secretMap {
		configInstance.configs[key] = val
	}
	return nil
}

type NotExistError struct {
	path string
}

func (errorInstance *NotExistError) Error() string {
	return fmt.Sprintf(`config of path [%s] not exist`, errorInstance.path)
}

func NewNotExistError(path string) *NotExistError {
	return &NotExistError{path:path}
}

func (configInstance *ConfigClass) parseYaml(arr []string, length int, path string) (map[interface{}]interface{}, error) {
	temp, ok := configInstance.configs[arr[1]].(map[interface{}]interface{})
	if !ok {
		return nil, NewNotExistError(path)
	}
	for _, v := range arr[2 : length-1] {
		temp, ok = temp[v].(map[interface{}]interface{})
		if !ok {
			return nil, NewNotExistError(path)
		}
	}
	return temp, nil
}

// merge flag config
// priority: flag > config file > flag default value
func (configInstance *ConfigClass) MergeFlagSet(flagSet *flag.FlagSet) {
	flagSet.Visit(func(f *flag.Flag) {
		configInstance.configs[f.Name] = f.Value.String()
	})

	flagSet.VisitAll(func(f *flag.Flag) {
		if _, ok := configInstance.configs[f.Name]; !ok {
			configInstance.configs[f.Name] = f.DefValue
		}
	})
}

func (configInstance *ConfigClass) parseJson(arr []string, length int, path string) (map[string]interface{}, error) {
	temp, ok := configInstance.configs[arr[1]].(map[string]interface{})
	if !ok {
		return nil, NewNotExistError(path)
	}
	for _, v := range arr[2 : length-1] {
		temp, ok = temp[v].(map[string]interface{})
		if !ok {
			return nil, NewNotExistError(path)
		}
	}
	return temp, nil
}

func (configInstance *ConfigClass) MustGetStringDefault(str string, default_ string) string {
	result, err := configInstance.GetStringDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetStringDefault(str string, default_ string) (string, error) {
	result, err := configInstance.GetString(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return ``, err
		}
	}
	return result, nil
}

func (configInstance *ConfigClass) findTarget(str string) (interface{}, error) {
	target := configInstance.configs[str]
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return ``, errors.New(`path error`)
		}
		if length == 2 {
			target = configInstance.configs[arr[1]]
		} else if configInstance.loadType == YAML_TYPE {
			temp, err := configInstance.parseYaml(arr, length, str)
			if err != nil {
				return ``, err
			}
			target = temp[arr[length-1]]
		} else {
			temp, err := configInstance.parseJson(arr, length, str)
			if err != nil {
				return ``, err
			}
			target = temp[arr[length-1]]
		}
	}
	if target == nil {
		return nil, NewNotExistError(str)
	}
	return target, nil
}

func (configInstance *ConfigClass) MustGetString(str string) string {
	result, err := configInstance.GetString(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetString(str string) (string, error) {
	target, err := configInstance.findTarget(str)
	if err != nil {
		return ``, err
	}
	result, err := go_reflect.Reflect.ToString(target)
	if err != nil {
		return ``, err
	}
	return result, nil
}

func (configInstance *ConfigClass) MustGetIntDefault(str string, default_ int) int {
	result, err := configInstance.GetIntDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetIntDefault(str string, default_ int) (int, error) {
	result, err := configInstance.GetInt(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (configInstance *ConfigClass) MustGetInt(str string) int {
	result, err := configInstance.GetInt(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetInt(str string) (int, error) {
	target, err := configInstance.findTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToInt(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (configInstance *ConfigClass) MustGetInt64Default(str string, default_ int64) int64 {
	result, err := configInstance.GetInt64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetInt64Default(str string, default_ int64) (int64, error) {
	result, err := configInstance.GetInt64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (configInstance *ConfigClass) MustGetInt64(str string) int64 {
	result, err := configInstance.GetInt64(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetInt64(str string) (int64, error) {
	target, err := configInstance.findTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToInt64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (configInstance *ConfigClass) MustGetUint64Default(str string, default_ uint64) uint64 {
	result, err := configInstance.GetUint64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetUint64Default(str string, default_ uint64) (uint64, error) {
	result, err := configInstance.GetUint64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (configInstance *ConfigClass) GetUint64(str string) (uint64, error) {
	target, err := configInstance.findTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToUint64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (configInstance *ConfigClass) MustGetBoolDefault(str string, default_ bool) bool {
	result, err := configInstance.GetBoolDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetBoolDefault(str string, default_ bool) (bool, error) {
	result, err := configInstance.GetBool(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return false, err
		}
	}
	return result, nil
}

func (configInstance *ConfigClass) GetBool(str string) (bool, error) {
	target, err := configInstance.findTarget(str)
	if err != nil {
		return false, err
	}
	result, err := go_reflect.Reflect.ToBool(target)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (configInstance *ConfigClass) MustGetFloat64Default(str string, default_ float64) float64 {
	result, err := configInstance.GetFloat64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetFloat64Default(str string, default_ float64) (float64, error) {
	result, err := configInstance.GetFloat64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (configInstance *ConfigClass) GetFloat64(str string) (float64, error) {
	target, err := configInstance.findTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToFloat64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (configInstance *ConfigClass) Get(str string) (interface{}, error) {
	return configInstance.findTarget(str)
}

func (configInstance *ConfigClass) GetAll() interface{} {
	return configInstance.configs
}

func (configInstance *ConfigClass) MustGetMapDefault(str string, default_ map[string]interface{}) map[string]interface{} {
	map_, err := configInstance.GetMapDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return map_
}

func (configInstance *ConfigClass) GetMapDefault(str string, default_ map[string]interface{}) (map[string]interface{}, error) {
	result, err := configInstance.GetMap(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (configInstance *ConfigClass) GetMap(str string) (map[string]interface{}, error) {
	target, err := configInstance.findTarget(str)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	if configInstance.loadType == YAML_TYPE {
		temp, ok := target.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New(`cast error`)
		}
		for k, v := range temp {
			key, err := go_reflect.Reflect.ToString(k)
			if err != nil {
				return nil, err
			}
			result[key] = v
		}
	} else {
		result_, ok := target.(map[string]interface{})
		if !ok {
			return nil, errors.New(`cast error`)
		}
		result = result_
	}
	return result, nil
}

func (configInstance *ConfigClass) MustGetStruct(str string, s interface{}) {
	err := configInstance.GetStruct(str, s)
	if err != nil {
		panic(err)
	}
}

func (configInstance *ConfigClass) GetStruct(str string, s interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	map_, err := configInstance.GetMap(str)
	if err != nil {
		return err
	}
	err = decoder.Decode(map_)
	if err != nil {
		return err
	}
	return nil
}

func (configInstance *ConfigClass) MustGetSliceDefault(str string, default_ []interface{}) []interface{} {
	result, err := configInstance.GetSliceDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigClass) GetSliceDefault(str string, default_ []interface{}) ([]interface{}, error) {
	result, err := configInstance.GetSlice(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (configInstance *ConfigClass) GetSlice(str string) ([]interface{}, error) {
	target, err := configInstance.findTarget(str)
	if err != nil {
		return nil, err
	}
	result, ok := target.([]interface{})
	if !ok {
		return nil, errors.New(`cast error`)
	}
	return result, nil
}
