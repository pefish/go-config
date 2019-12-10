package go_config

import (
	"encoding/json"
	"errors"
	"flag"
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

var Config = ConfigClass{}

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

func (this *ConfigClass) parseYaml(arr []string, length int) (map[interface{}]interface{}, error) {
	temp, ok := this.configs[arr[1]].(map[interface{}]interface{})
	if !ok {
		return nil, errors.New(`parse error`)
	}
	for _, v := range arr[2 : length-1] {
		temp, ok = temp[v].(map[interface{}]interface{})
		if !ok {
			return nil, errors.New(`parse error`)
		}
	}
	return temp, nil
}

func (this *ConfigClass) MergeFlag() {
	flag.VisitAll(func(i *flag.Flag) {
		this.configs[i.Name] = i.Value.String()
	})
}

func (this *ConfigClass) parseJson(arr []string, length int) (map[string]interface{}, error) {
	temp, ok := this.configs[arr[1]].(map[string]interface{})
	if !ok {
		return nil, errors.New(`parse error`)
	}
	for _, v := range arr[2 : length-1] {
		temp, ok = temp[v].(map[string]interface{})
		if !ok {
			return nil, errors.New(`parse error`)
		}
	}
	return temp, nil
}

func (this *ConfigClass) GetString(str string) string {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return ``
		}
		if length == 2 {
			result, _ := go_reflect.Reflect.ToString(this.configs[arr[1]])
			return result
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				return ``
			}
			result, _ := go_reflect.Reflect.ToString(temp[arr[length-1]])
			return result
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				return ``
			}
			result, _ := go_reflect.Reflect.ToString(temp[arr[length-1]])
			return result
		}
	}
	result, _ := go_reflect.Reflect.ToString(this.configs[str])
	return result
}

func (this *ConfigClass) GetInt(str string) int {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return 0
		}
		if length == 2 {
			result, _ := go_reflect.Reflect.ToInt(this.configs[arr[1]])
			return result
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				return 0
			}
			result, _ := go_reflect.Reflect.ToInt(temp[arr[length-1]])
			return result
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				return 0
			}
			result, _ := go_reflect.Reflect.ToInt(temp[arr[length-1]])
			return result
		}
	}
	result, _ := go_reflect.Reflect.ToInt(this.configs[str])
	return result
}

func (this *ConfigClass) GetInt64(str string) int64 {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return 0
		}
		if length == 2 {
			result, _ := go_reflect.Reflect.ToInt64(this.configs[arr[1]])
			return result
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				return 0
			}
			result, _ := go_reflect.Reflect.ToInt64(temp[arr[length-1]])
			return result
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				return 0
			}
			result, _ := go_reflect.Reflect.ToInt64(temp[arr[length-1]])
			return result
		}
	}
	result, _ := go_reflect.Reflect.ToInt64(this.configs[str])
	return result
}

func (this *ConfigClass) GetUint64(str string) uint64 {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return 0
		}
		if length == 2 {
			result, _ := go_reflect.Reflect.ToUint64(this.configs[arr[1]])
			return result
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				return 0
			}
			result, _ := go_reflect.Reflect.ToUint64(temp[arr[length-1]])
			return result
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				return 0
			}
			result, _ := go_reflect.Reflect.ToUint64(temp[arr[length-1]])
			return result
		}
	}
	result, _ := go_reflect.Reflect.ToUint64(this.configs[str])
	return result
}

func (this *ConfigClass) GetBool(str string) bool {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return false
		}
		if length == 2 {
			result, _ := go_reflect.Reflect.ToBool(this.configs[arr[1]])
			return result
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				return false
			}
			result, _ := go_reflect.Reflect.ToBool(temp[arr[length-1]])
			return result
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				return false
			}
			result, _ := go_reflect.Reflect.ToBool(temp[arr[length-1]])
			return result
		}
	}
	result, _ := go_reflect.Reflect.ToBool(this.configs[str])
	return result
}

func (this *ConfigClass) GetFloat64(str string) float64 {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return 0
		}
		if length == 2 {
			result, _ := go_reflect.Reflect.ToFloat64(this.configs[arr[1]])
			return result
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				return 0
			}
			result, _ := go_reflect.Reflect.ToFloat64(temp[arr[length-1]])
			return result
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				return 0
			}
			result, _ := go_reflect.Reflect.ToFloat64(temp[arr[length-1]])
			return result
		}
	}
	result, _ := go_reflect.Reflect.ToFloat64(this.configs[str])
	return result
}

func (this *ConfigClass) Get(str string) interface{} {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return nil
		}
		if length == 2 {
			return this.configs[arr[1]]
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				return nil
			}
			return temp[arr[length-1]]
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				return nil
			}
			return temp[arr[length-1]]
		}
	}
	return this.configs[str]
}

func (this *ConfigClass) GetAll() interface{} {
	return this.configs
}

func (this *ConfigClass) MustGetMap(str string) map[string]interface{} {
	result := map[string]interface{}{}
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.MustGetMap(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				panic(err)
			}
			temp1 := temp[arr[length-1]].(map[interface{}]interface{})
			for k, v := range temp1 {
				result[go_reflect.Reflect.MustToString(k)] = v
			}
			return result
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				panic(err)
			}
			return temp[arr[length-1]].(map[string]interface{})
		}
	}

	if this.loadType == YAML_TYPE {
		temp := this.configs[str].(map[interface{}]interface{})
		for k, v := range temp {
			result[go_reflect.Reflect.MustToString(k)] = v
		}
	} else {
		result = this.configs[str].(map[string]interface{})
	}
	return result
}

func (this *ConfigClass) MustGetStruct(str string, s interface{}) {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(this.MustGetMap(str))
	if err != nil {
		panic(err)
	}
}

func (this *ConfigClass) MustGetSlice(str string) []interface{} {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.configs[arr[1]].([]interface{})
		}
		if this.loadType == YAML_TYPE {
			temp, err := this.parseYaml(arr, length)
			if err != nil {
				panic(err)
			}
			return temp[arr[length-1]].([]interface{})
		} else {
			temp, err := this.parseJson(arr, length)
			if err != nil {
				panic(err)
			}
			return temp[arr[length-1]].([]interface{})
		}
	}
	return this.configs[str].([]interface{})
}
