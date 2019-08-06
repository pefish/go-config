package go_config

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/pefish/go-reflect"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ConfigClass struct {
	configs map[string]interface{}
}

var Config = ConfigClass{}

type Configuration struct {
	ConfigFilepath string
	SecretFilepath string
}

func (this *ConfigClass) LoadYamlConfig(config Configuration) {
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

func (this *ConfigClass) GetString(str string) string {
	return go_reflect.Reflect.ToString(this.configs[str])
}

func (this *ConfigClass) GetInt(str string) int {
	return go_reflect.Reflect.ToInt(this.configs[str])
}

func (this *ConfigClass) GetInt64(str string) int64 {
	return go_reflect.Reflect.ToInt64(this.configs[str])
}

func (this *ConfigClass) GetUint64(str string) uint64 {
	return go_reflect.Reflect.ToUint64(this.configs[str])
}

func (this *ConfigClass) GetBool(str string) bool {
	return go_reflect.Reflect.ToBool(this.configs[str])
}

func (this *ConfigClass) GetFloat64(str string) float64 {
	return go_reflect.Reflect.ToFloat64(this.configs[str])
}

func (this *ConfigClass) Get(str string) interface{} {
	return this.configs[str]
}

func (this *ConfigClass) GetMap(str string) map[string]interface{} {
	result := map[string]interface{}{}
	switch this.configs[str].(type) {
	case map[interface{}]interface{}:
		temp := this.configs[str].(map[interface{}]interface{})
		for k, v := range temp {
			result[go_reflect.Reflect.ToString(k)] = v
		}
	default:
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
	return this.configs[str].([]interface{})
}

func (this *ConfigClass) GetSliceString(str string) []string {
	return this.configs[str].([]string)
}

func (this *ConfigClass) GetSliceWithErr(str string) ([]interface{}, error) {
	result, ok := this.configs[str].([]interface{})
	if !ok {
		return nil, errors.New(`type assert error`)
	}
	return result, nil
}
