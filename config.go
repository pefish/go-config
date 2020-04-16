package go_config

import (
	"encoding/json"
	"errors"
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

func (this *ConfigClass) MustLoadYamlConfig(config Configuration) {
	this.loadType = YAML_TYPE

	configFile := ``
	configMap := map[string]interface{}{}
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
					panic(err)
				}
			}
		}
	}

	secretFile := ``
	secretMap := map[string]interface{}{}
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
					panic(err)
				}
			}
		}
	}

	if configFile == `` && secretFile == `` {
		panic(errors.New(`unspecified config file and secret file`))
	}
	this.configs = configMap
	for key, val := range secretMap {
		this.configs[key] = val
	}
}

func (this *ConfigClass) MustLoadJsonConfig(config Configuration) {
	this.loadType = JSON_TYPE
	configFile := ``
	configMap := map[string]interface{}{}
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
					panic(err)
				}
				configMap = result.(map[string]interface{})
			}
		}
	}

	secretFile := ``
	secretMap := map[string]interface{}{}
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
					panic(err)
				}
				secretMap = result.(map[string]interface{})
			}
		}
	}

	if configFile == `` && secretFile == `` {
		panic(errors.New(`unspecified config file and secret file`))
	}
	this.configs = configMap
	for key, val := range secretMap {
		this.configs[key] = val
	}
}

type NotExistError struct {

}

func (errorInstance *NotExistError) Error() string {
	return `not exist`
}

func (this *ConfigClass) parseYaml(arr []string, length int) (map[interface{}]interface{}, error) {
	temp, ok := this.configs[arr[1]].(map[interface{}]interface{})
	if !ok {
		return nil, &NotExistError{}
	}
	for _, v := range arr[2 : length-1] {
		temp, ok = temp[v].(map[interface{}]interface{})
		if !ok {
			return nil, &NotExistError{}
		}
	}
	return temp, nil
}

// merge命令行参数
func (this *ConfigClass) MustMergeFlag() {
	result, err := Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	for k, v := range result {
		this.configs[k] = v
	}
}

func (this *ConfigClass) parseJson(arr []string, length int) (map[string]interface{}, error) {
	temp, ok := this.configs[arr[1]].(map[string]interface{})
	if !ok {
		return nil, &NotExistError{}
	}
	for _, v := range arr[2 : length-1] {
		temp, ok = temp[v].(map[string]interface{})
		if !ok {
			return nil, &NotExistError{}
		}
	}
	return temp, nil
}

func (this *ConfigClass) MustGetStringDefault(str string, default_ string) string {
	result, err := this.GetStringDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetStringDefault(str string, default_ string) (string, error) {
	result, err := this.GetString(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return ``, err
		}
	}
	return result, nil
}

func (this *ConfigClass) findTarget(str string) (interface{}, error) {
	target := this.configs[str]
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return ``, errors.New(`path error`)
		}
		if length == 2 {
			target = this.configs[arr[1]]
		} else if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				return ``, err
			}
			target = temp[arr[length-1]]
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				return ``, err
			}
			target = temp[arr[length-1]]
		}
	}
	if target == nil {
		return nil, &NotExistError{}
	}
	return target, nil
}

func (this *ConfigClass) MustGetString(str string) string {
	result, err := this.GetString(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetString(str string) (string, error) {
	target, err := this.findTarget(str)
	if err != nil {
		return ``, err
	}
	result, err := go_reflect.Reflect.ToString(target)
	if err != nil {
		return ``, err
	}
	return result, nil
}

func (this *ConfigClass) MustGetIntDefault(str string, default_ int) int {
	result, err := this.GetIntDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetIntDefault(str string, default_ int) (int, error) {
	result, err := this.GetInt(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (this *ConfigClass) MustGetInt(str string) int {
	result, err := this.GetInt(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetInt(str string) (int, error) {
	target, err := this.findTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToInt(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (this *ConfigClass) MustGetInt64Default(str string, default_ int64) int64 {
	result, err := this.GetInt64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetInt64Default(str string, default_ int64) (int64, error) {
	result, err := this.GetInt64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (this *ConfigClass) MustGetInt64(str string) int64 {
	result, err := this.GetInt64(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetInt64(str string) (int64, error) {
	target, err := this.findTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToInt64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (this *ConfigClass) MustGetUint64Default(str string, default_ uint64) uint64 {
	result, err := this.GetUint64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetUint64Default(str string, default_ uint64) (uint64, error) {
	result, err := this.GetUint64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (this *ConfigClass) GetUint64(str string) (uint64, error) {
	target, err := this.findTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToUint64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (this *ConfigClass) MustGetBoolDefault(str string, default_ bool) bool {
	result, err := this.GetBoolDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetBoolDefault(str string, default_ bool) (bool, error) {
	result, err := this.GetBool(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return false, err
		}
	}
	return result, nil
}

func (this *ConfigClass) GetBool(str string) (bool, error) {
	target, err := this.findTarget(str)
	if err != nil {
		return false, err
	}
	result, err := go_reflect.Reflect.ToBool(target)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (this *ConfigClass) MustGetFloat64Default(str string, default_ float64) float64 {
	result, err := this.GetFloat64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetFloat64Default(str string, default_ float64) (float64, error) {
	result, err := this.GetFloat64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (this *ConfigClass) GetFloat64(str string) (float64, error) {
	target, err := this.findTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToFloat64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (this *ConfigClass) Get(str string) (interface{}, error) {
	return this.findTarget(str)
}

func (this *ConfigClass) GetAll() interface{} {
	return this.configs
}

func (this *ConfigClass) MustGetMap(str string) map[string]interface{} {
	map_, err := this.GetMap(str)
	if err != nil {
		panic(err)
	}
	return map_
}

func (this *ConfigClass) GetMap(str string) (map[string]interface{}, error) {
	target, err := this.findTarget(str)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	if this.loadType == YAML_TYPE {
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

func (this *ConfigClass) MustGetStruct(str string, s interface{}) {
	err := this.GetStruct(str, s)
	if err != nil {
		panic(err)
	}
}

func (this *ConfigClass) GetStruct(str string, s interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	map_, err := this.GetMap(str)
	if err != nil {
		return err
	}
	err = decoder.Decode(map_)
	if err != nil {
		return err
	}
	return nil
}

func (this *ConfigClass) MustGetSlice(str string) []interface{} {
	result, err := this.GetSlice(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ConfigClass) GetSlice(str string) ([]interface{}, error) {
	target, err := this.findTarget(str)
	if err != nil {
		return nil, err
	}
	result, ok := target.([]interface{})
	if !ok {
		return nil, errors.New(`cast error`)
	}
	return result, nil
}

func (this *ConfigClass) MustGetStringSlice(str string) []string {
	var result []string
	results := this.MustGetSlice(str)
	for _, v := range results {
		result = append(result, v.(string))
	}
	return result
}

func (this *ConfigClass) MustGetUint64Slice(str string) []uint64 {
	var result []uint64
	results := this.MustGetSlice(str)
	for _, v := range results {
		result = append(result, go_reflect.Reflect.MustToUint64(v))
	}
	return result
}
