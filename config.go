package go_config

import (
	"flag"
	"fmt"
	"github.com/mitchellh/mapstructure"
	go_format "github.com/pefish/go-format"
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

func (c *ConfigManager) MustLoadConfig(config Configuration) {
	err := c.LoadConfig(config)
	if err != nil {
		panic(err)
	}
}

func (c *ConfigManager) Unmarshal(out interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &out,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(c.configs)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigManager) LoadConfig(config Configuration) error {
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

	c.fileConfigs = configMap
	for key, val := range secretMap {
		c.fileConfigs[key] = val
	}
	c.combineConfigs()
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

func (c *ConfigManager) parseYaml(arr []string, length int, path string) (map[interface{}]interface{}, error) {
	temp, ok := c.configs[arr[1]].(map[interface{}]interface{})
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
func (c *ConfigManager) MergeFlagSet(flagSet *flag.FlagSet) {
	flagSet.Visit(func(f *flag.Flag) {
		c.flagSetConfigs[f.Name] = f.Value.String()
	})

	flagSet.VisitAll(func(f *flag.Flag) {
		if _, ok := c.flagSetConfigs[f.Name]; !ok {
			c.flagSetDefaultConfigs[f.Name] = f.DefValue
		}
	})
	//fmt.Println(c.flagSetDefaultConfigs)
	c.combineConfigs()
}

// merge envs
// priority: flag > envs > config file > flag default value
func (c *ConfigManager) MergeEnvs(envKeyPair map[string]string) {
	for envName, keyName := range envKeyPair {
		envValue := os.Getenv(envName)
		if envValue != "" {
			c.envConfigs[keyName] = envValue
		}
	}

	c.combineConfigs()

}

func (c *ConfigManager) combineConfigs() {
	for key, value := range c.flagSetDefaultConfigs {
		c.configs[key] = value
	}
	//fmt.Println(c.configs)

	for key, value := range c.fileConfigs {
		c.configs[key] = value
	}
	//fmt.Println(c.flagSetDefaultConfigs)

	for key, value := range c.envConfigs {
		c.configs[key] = value
	}
	//fmt.Println(c.envConfigs)

	for key, value := range c.flagSetConfigs {
		c.configs[key] = value
	}
	//fmt.Println(c.flagSetConfigs)
}

func (c *ConfigManager) MustGetStringDefault(str string, default_ string) string {
	result, err := c.GetStringDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetStringDefault(str string, default_ string) (string, error) {
	result, err := c.GetString(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return ``, err
		}
	}
	return result, nil
}

func (c *ConfigManager) FindTarget(str string) (interface{}, error) {
	target := c.configs[str]
	//fmt.Println(target)
	if strings.HasPrefix(str, `/`) {
		arr := strings.Split(str, `/`)
		length := len(arr)
		if length <= 1 {
			return ``, errors.New(`path error`)
		}
		if length == 2 {
			target = c.configs[arr[1]]
		} else {
			temp, err := c.parseYaml(arr, length, str)
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

func (c *ConfigManager) MustGetString(str string) string {
	result, err := c.GetString(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetString(str string) (string, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		return ``, err
	}
	result := go_format.FormatInstance.ToString(target)
	return result, nil
}

func (c *ConfigManager) MustGetIntDefault(str string, default_ int) int {
	result, err := c.GetIntDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetIntDefault(str string, default_ int) (int, error) {
	result, err := c.GetInt(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (c *ConfigManager) MustGetInt(str string) int {
	result, err := c.GetInt(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetInt(str string) (int, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_format.FormatInstance.ToInt(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *ConfigManager) MustGetInt64Default(str string, default_ int64) int64 {
	result, err := c.GetInt64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetInt64Default(str string, default_ int64) (int64, error) {
	result, err := c.GetInt64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (c *ConfigManager) MustGetInt64(str string) int64 {
	result, err := c.GetInt64(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetInt64(str string) (int64, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_format.FormatInstance.ToInt64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *ConfigManager) MustGetUint64Default(str string, default_ uint64) uint64 {
	result, err := c.GetUint64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetUint64Default(str string, default_ uint64) (uint64, error) {
	result, err := c.GetUint64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (c *ConfigManager) MustGetUint64(str string) uint64 {
	result, err := c.GetUint64(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetUint64(str string) (uint64, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_format.FormatInstance.ToUint64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *ConfigManager) MustGetBoolDefault(str string, default_ bool) bool {
	result, err := c.GetBoolDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) MustGetBool(str string) bool {
	result, err := c.GetBool(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetBoolDefault(str string, default_ bool) (bool, error) {
	result, err := c.GetBool(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return false, err
		}
	}
	return result, nil
}

func (c *ConfigManager) GetBool(str string) (bool, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		return false, err
	}
	result, err := go_format.FormatInstance.ToBool(target)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (c *ConfigManager) MustGetFloat64Default(str string, default_ float64) float64 {
	result, err := c.GetFloat64Default(str, default_)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) GetFloat64Default(str string, default_ float64) (float64, error) {
	result, err := c.GetFloat64(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return 0, err
		}
	}
	return result, nil
}

func (c *ConfigManager) GetFloat64(str string) (float64, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		return 0, err
	}
	result, err := go_format.FormatInstance.ToFloat64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *ConfigManager) Configs() map[string]interface{} {
	return c.configs
}

func (c *ConfigManager) FlagSetDefaultConfigs() map[string]interface{} {
	return c.flagSetDefaultConfigs
}

func (c *ConfigManager) FileConfigs() map[string]interface{} {
	return c.fileConfigs
}

func (c *ConfigManager) EnvConfigs() map[string]interface{} {
	return c.envConfigs
}

func (c *ConfigManager) FlagSetConfigs() map[string]interface{} {
	return c.flagSetConfigs
}

func (c *ConfigManager) MustGetMapDefault(str string, default_ map[string]interface{}) map[string]interface{} {
	map_, err := c.GetMapDefault(str, default_)
	if err != nil {
		panic(err)
	}
	return map_
}

func (c *ConfigManager) MustGetMap(str string) map[string]interface{} {
	map_, err := c.GetMap(str)
	if err != nil {
		panic(err)
	}
	return map_
}

func (c *ConfigManager) GetMapDefault(str string, default_ map[string]interface{}) (map[string]interface{}, error) {
	result, err := c.GetMap(str)
	if err != nil {
		if _, ok := err.(*NotExistError); ok {
			return default_, nil
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (c *ConfigManager) GetMap(str string) (map[string]interface{}, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	temp, ok := target.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New(`cast error`)
	}
	for k, v := range temp {
		key := go_format.FormatInstance.ToString(k)
		result[key] = v
	}
	return result, nil
}

func (c *ConfigManager) MustGet(str string, s interface{}) {
	err := c.Get(str, s)
	if err != nil {
		panic(err)
	}
}

func (c *ConfigManager) Get(str string, s interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	target, err := c.FindTarget(str)
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
