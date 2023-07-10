package go_config

import (
	"flag"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pefish/go-reflect"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type ConfigManager struct {
	flagSetConfigs        map[string]interface{}
	flagSetDefaultConfigs map[string]interface{}
	envConfigs            map[string]interface{}
	fileConfigs           map[string]interface{}
	configs               map[string]interface{}
}

var ConfigManagerInstance = NewConfigManager()

type Configuration struct {
	ConfigFilepath string
	SecretFilepath string
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs:               make(map[string]interface{}, 5),
		fileConfigs:           make(map[string]interface{}, 5),
		flagSetConfigs:        make(map[string]interface{}, 2),
		flagSetDefaultConfigs: make(map[string]interface{}, 2),
		envConfigs:            make(map[string]interface{}, 2),
	}
}

func (configInstance *ConfigManager) MustLoadConfig(config Configuration) {
	err := configInstance.LoadConfig(config)
	if err != nil {
		panic(err)
	}
}

func (configInstance *ConfigManager) Unmarshal(out interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &out,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(configInstance.configs)
	if err != nil {
		return err
	}
	return nil
}

func (configInstance *ConfigManager) LoadConfig(config Configuration) error {
	configMap := make(map[string]interface{})
	if config.ConfigFilepath != `` {
		if !strings.HasSuffix(config.ConfigFilepath, "yaml") {
			return errors.New("only support yaml file")
		}
		bytes, err := os.ReadFile(config.ConfigFilepath)
		if err == nil {
			err = yaml.Unmarshal(bytes, &configMap)
			if err != nil {
				return errors.WithMessage(err, "file format error")
			}
		}
	}

	secretMap := make(map[string]interface{})
	if config.SecretFilepath != `` {
		if !strings.HasSuffix(config.SecretFilepath, "yaml") {
			return errors.New("only support yaml file")
		}
		bytes, err := os.ReadFile(config.SecretFilepath)
		if err == nil {
			err = yaml.Unmarshal(bytes, &secretMap)
			if err != nil {
				return errors.WithMessage(err, "file format error")
			}
		}
	}

	configInstance.fileConfigs = configMap
	for key, val := range secretMap {
		configInstance.fileConfigs[key] = val
	}
	configInstance.combineConfigs()
	return nil
}

type NotExistError struct {
	path string
}

func (errorInstance *NotExistError) Error() string {
	return fmt.Sprintf(`config of path [%s] not exist`, errorInstance.path)
}

func NewNotExistError(path string) *NotExistError {
	return &NotExistError{path: path}
}

func (configInstance *ConfigManager) parseYaml(arr []string, length int, path string) (map[interface{}]interface{}, error) {
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
// just cover
func (configInstance *ConfigManager) MergeFlagSet(flagSet *flag.FlagSet) {
	flagSet.Visit(func(f *flag.Flag) {
		configInstance.flagSetConfigs[f.Name] = f.Value.String()
	})

	flagSet.VisitAll(func(f *flag.Flag) {
		if _, ok := configInstance.flagSetConfigs[f.Name]; !ok {
			configInstance.flagSetDefaultConfigs[f.Name] = f.DefValue
		}
	})
	//fmt.Println(configInstance.flagSetDefaultConfigs)
	configInstance.combineConfigs()
}

// merge envs
// priority: flag > envs > config file > flag default value
func (configInstance *ConfigManager) MergeEnvs(envKeyPair map[string]string) {
	for envName, keyName := range envKeyPair {
		envValue := os.Getenv(envName)
		if envValue != "" {
			configInstance.envConfigs[keyName] = envValue
		}
	}

	configInstance.combineConfigs()

}

func (configInstance *ConfigManager) combineConfigs() {
	for key, value := range configInstance.flagSetDefaultConfigs {
		configInstance.configs[key] = value
	}
	//fmt.Println(configInstance.configs)

	for key, value := range configInstance.fileConfigs {
		configInstance.configs[key] = value
	}
	//fmt.Println(configInstance.flagSetDefaultConfigs)

	for key, value := range configInstance.envConfigs {
		configInstance.configs[key] = value
	}
	//fmt.Println(configInstance.envConfigs)

	for key, value := range configInstance.flagSetConfigs {
		configInstance.configs[key] = value
	}
	//fmt.Println(configInstance.flagSetConfigs)
}

func (configInstance *ConfigManager) MustGetStringDefault(str string, default_ string) string {
	result, err := configInstance.GetStringDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetStringDefault(str string, default_ string) (string, error) {
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

func (configInstance *ConfigManager) FindTarget(str string) (interface{}, error) {
	target := configInstance.configs[str]
	//fmt.Println(target)
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return ``, errors.New(`path error`)
		}
		if length == 2 {
			target = configInstance.configs[arr[1]]
		} else {
			temp, err := configInstance.parseYaml(arr, length, str)
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

func (configInstance *ConfigManager) MustGetString(str string) string {
	result, err := configInstance.GetString(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetString(str string) (string, error) {
	target, err := configInstance.FindTarget(str)
	if err != nil {
		return ``, err
	}
	result := go_reflect.Reflect.ToString(target)
	return result, nil
}

func (configInstance *ConfigManager) MustGetIntDefault(str string, default_ int) int {
	result, err := configInstance.GetIntDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetIntDefault(str string, default_ int) (int, error) {
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

func (configInstance *ConfigManager) MustGetInt(str string) int {
	result, err := configInstance.GetInt(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetInt(str string) (int, error) {
	target, err := configInstance.FindTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToInt(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (configInstance *ConfigManager) MustGetInt64Default(str string, default_ int64) int64 {
	result, err := configInstance.GetInt64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetInt64Default(str string, default_ int64) (int64, error) {
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

func (configInstance *ConfigManager) MustGetInt64(str string) int64 {
	result, err := configInstance.GetInt64(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetInt64(str string) (int64, error) {
	target, err := configInstance.FindTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToInt64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (configInstance *ConfigManager) MustGetUint64Default(str string, default_ uint64) uint64 {
	result, err := configInstance.GetUint64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetUint64Default(str string, default_ uint64) (uint64, error) {
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

func (configInstance *ConfigManager) GetUint64(str string) (uint64, error) {
	target, err := configInstance.FindTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToUint64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (configInstance *ConfigManager) MustGetBoolDefault(str string, default_ bool) bool {
	result, err := configInstance.GetBoolDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) MustGetBool(str string) bool {
	result, err := configInstance.GetBool(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetBoolDefault(str string, default_ bool) (bool, error) {
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

func (configInstance *ConfigManager) GetBool(str string) (bool, error) {
	target, err := configInstance.FindTarget(str)
	if err != nil {
		return false, err
	}
	result, err := go_reflect.Reflect.ToBool(target)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (configInstance *ConfigManager) MustGetFloat64Default(str string, default_ float64) float64 {
	result, err := configInstance.GetFloat64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (configInstance *ConfigManager) GetFloat64Default(str string, default_ float64) (float64, error) {
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

func (configInstance *ConfigManager) GetFloat64(str string) (float64, error) {
	target, err := configInstance.FindTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_reflect.Reflect.ToFloat64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (configInstance *ConfigManager) Configs() map[string]interface{} {
	return configInstance.configs
}

func (configInstance *ConfigManager) FlagSetDefaultConfigs() map[string]interface{} {
	return configInstance.flagSetDefaultConfigs
}

func (configInstance *ConfigManager) FileConfigs() map[string]interface{} {
	return configInstance.fileConfigs
}

func (configInstance *ConfigManager) EnvConfigs() map[string]interface{} {
	return configInstance.envConfigs
}

func (configInstance *ConfigManager) FlagSetConfigs() map[string]interface{} {
	return configInstance.flagSetConfigs
}

func (configInstance *ConfigManager) MustGetMapDefault(str string, default_ map[string]interface{}) map[string]interface{} {
	map_, err := configInstance.GetMapDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return map_
}

func (configInstance *ConfigManager) MustGetMap(str string) map[string]interface{} {
	map_, err := configInstance.GetMap(str)
	if err != nil {
		panic(err)
	}
	return map_
}

func (configInstance *ConfigManager) GetMapDefault(str string, default_ map[string]interface{}) (map[string]interface{}, error) {
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

func (configInstance *ConfigManager) GetMap(str string) (map[string]interface{}, error) {
	target, err := configInstance.FindTarget(str)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	temp, ok := target.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New(`cast error`)
	}
	for k, v := range temp {
		key := go_reflect.Reflect.ToString(k)
		result[key] = v
	}
	return result, nil
}

func (configInstance *ConfigManager) MustGet(str string, s interface{}) {
	err := configInstance.Get(str, s)
	if err != nil {
		panic(err)
	}
}

func (configInstance *ConfigManager) Get(str string, s interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	target, err := configInstance.FindTarget(str)
	if err != nil {
		return err
	}
	//fmt.Println(target)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}
	return nil
}
