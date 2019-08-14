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

var Config = ConfigClass{}

type Configuration struct {
	ConfigFilepath string
	SecretFilepath string
}

func (this *ConfigClass) LoadYamlConfig(config Configuration) {
	this.loadType = YAML_TYPE
	configFile := ``
	configMap := map[string]interface{}{}
	if config.ConfigFilepath == `` {
		configFile = os.Getenv(`GO_CONFIG`)
	} else {
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

	secretFile := ``
	secretMap := map[string]interface{}{}
	if config.SecretFilepath == `` {
		secretFile = os.Getenv(`GO_SECRET`)
	} else {
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

	if configFile == `` && secretFile == `` {
		panic(errors.New(`unspecified config file and secret file`))
	}
	this.configs = configMap
	for key, val := range secretMap {
		this.configs[key] = val
	}
}

func (this *ConfigClass) LoadJsonConfig(config Configuration) {
	this.loadType = JSON_TYPE
	configFile := ``
	configMap := map[string]interface{}{}
	if config.ConfigFilepath == `` {
		configFile = os.Getenv(`GO_CONFIG`)
	} else {
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

	secretFile := ``
	secretMap := map[string]interface{}{}
	if config.SecretFilepath == `` {
		secretFile = os.Getenv(`GO_SECRET`)
	} else {
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

	if configFile == `` && secretFile == `` {
		panic(errors.New(`unspecified config file and secret file`))
	}
	this.configs = configMap
	for key, val := range secretMap {
		this.configs[key] = val
	}
}

func (this *ConfigClass) parseYaml(arr []string, length int) map[interface{}]interface{} {
	temp := this.configs[arr[1]].(map[interface{}]interface{})
	for _, v := range arr[2 : length-1] {
		temp = temp[v].(map[interface{}]interface{})
	}
	return temp
}

func (this *ConfigClass) parseJson(arr []string, length int) map[string]interface{} {
	temp := this.configs[arr[1]].(map[string]interface{})
	for _, v := range arr[2 : length-1] {
		temp = temp[v].(map[string]interface{})
	}
	return temp
}

func (this *ConfigClass) GetString(str string) string {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetString(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			return go_reflect.Reflect.ToString(temp[arr[length-1]])
		} else {
			temp := this.parseJson(arr, length)
			return go_reflect.Reflect.ToString(temp[arr[length-1]])
		}
	}
	return go_reflect.Reflect.ToString(this.configs[str])
}

func (this *ConfigClass) GetInt(str string) int {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetInt(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			return go_reflect.Reflect.ToInt(temp[arr[length-1]])
		} else {
			temp := this.parseJson(arr, length)
			return go_reflect.Reflect.ToInt(temp[arr[length-1]])
		}
	}
	return go_reflect.Reflect.ToInt(this.configs[str])
}

func (this *ConfigClass) GetInt64(str string) int64 {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetInt64(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			return go_reflect.Reflect.ToInt64(temp[arr[length-1]])
		} else {
			temp := this.parseJson(arr, length)
			return go_reflect.Reflect.ToInt64(temp[arr[length-1]])
		}
	}
	return go_reflect.Reflect.ToInt64(this.configs[str])
}

func (this *ConfigClass) GetUint64(str string) uint64 {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetUint64(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			return go_reflect.Reflect.ToUint64(temp[arr[length-1]])
		} else {
			temp := this.parseJson(arr, length)
			return go_reflect.Reflect.ToUint64(temp[arr[length-1]])
		}
	}
	return go_reflect.Reflect.ToUint64(this.configs[str])
}

func (this *ConfigClass) GetBool(str string) bool {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetBool(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			return go_reflect.Reflect.ToBool(temp[arr[length-1]])
		} else {
			temp := this.parseJson(arr, length)
			return go_reflect.Reflect.ToBool(temp[arr[length-1]])
		}
	}
	return go_reflect.Reflect.ToBool(this.configs[str])
}

func (this *ConfigClass) GetFloat64(str string) float64 {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetFloat64(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			return go_reflect.Reflect.ToFloat64(temp[arr[length-1]])
		} else {
			temp := this.parseJson(arr, length)
			return go_reflect.Reflect.ToFloat64(temp[arr[length-1]])
		}
	}
	return go_reflect.Reflect.ToFloat64(this.configs[str])
}

func (this *ConfigClass) Get(str string) interface{} {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetUint64(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			return temp[arr[length-1]]
		} else {
			temp := this.parseJson(arr, length)
			return temp[arr[length-1]]
		}
	}
	return this.configs[str]
}

func (this *ConfigClass) GetAll() interface{} {
	return this.configs
}

func (this *ConfigClass) GetMap(str string) map[string]interface{} {
	result := map[string]interface{}{}
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetMap(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			temp1 := temp[arr[length-1]].(map[interface{}]interface{})
			for k, v := range temp1 {
				result[go_reflect.Reflect.ToString(k)] = v
			}
			return result
		} else {
			temp := this.parseJson(arr, length)
			return temp[arr[length-1]].(map[string]interface{})
		}
	}

	if this.loadType == YAML_TYPE {
		temp := this.configs[str].(map[interface{}]interface{})
		for k, v := range temp {
			result[go_reflect.Reflect.ToString(k)] = v
		}
	} else {
		result = this.configs[str].(map[string]interface{})
	}
	return result
}

func (this *ConfigClass) GetStruct(str string, s interface{}) {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(this.GetMap(str))
	if err != nil {
		panic(err)
	}
}

func (this *ConfigClass) GetSlice(str string) []interface{} {
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			panic(`path error`)
		}
		if length == 2 {
			return this.GetSlice(arr[1])
		}
		if this.loadType == YAML_TYPE {
			temp := this.parseYaml(arr, length)
			return temp[arr[length-1]].([]interface{})
		} else {
			temp := this.parseJson(arr, length)
			return temp[arr[length-1]].([]interface{})
		}
	}
	return this.configs[str].([]interface{})
}
