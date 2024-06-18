package go_config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	go_format "github.com/pefish/go-format"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ConfigManager struct {
	flagSetConfigs        map[string]interface{}
	flagSetDefaultConfigs map[string]interface{}
	fileConfigs           map[string]interface{}
	configs               map[string]interface{}
}

var ConfigManagerInstance = NewConfigManager()

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs:               make(map[string]interface{}, 5),
		fileConfigs:           make(map[string]interface{}, 5),
		flagSetConfigs:        make(map[string]interface{}, 2),
		flagSetDefaultConfigs: make(map[string]interface{}, 2),
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

// merge config file
func (c *ConfigManager) MergeConfigFile(configFilepath string) error {
	configMap := make(map[string]interface{})
	if !strings.HasSuffix(configFilepath, "yaml") {
		return errors.New("Only support yaml file.")
	}
	bytes, err := os.ReadFile(configFilepath)
	if err == nil {
		err = yaml.Unmarshal(bytes, &configMap)
		if err != nil {
			return errors.WithMessage(err, "File format error.")
		}
	}

	c.fileConfigs = configMap
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

// priority: env > flag > config file > flag default value
func (c *ConfigManager) combineConfigs() {
	for key, value := range c.flagSetDefaultConfigs {
		c.configs[key] = value
	}
	//fmt.Println(c.configs)

	for key, value := range c.fileConfigs {
		c.configs[key] = value
	}
	//fmt.Println(c.flagSetDefaultConfigs)

	for key, value := range c.flagSetConfigs {
		c.configs[key] = value
	}
	//fmt.Println(c.flagSetConfigs)

	// 查找环境变量中有没有匹配的配置项
	for k := range c.Configs() {
		envValue := os.Getenv(strings.ReplaceAll(strings.ToUpper(k), "-", "_"))
		if envValue != "" {
			c.configs[k] = envValue
		}
	}
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

func (c *ConfigManager) MustString(str string) string {
	result, err := c.String(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) String(str string) (string, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		var notExistError *NotExistError
		if errors.As(err, &notExistError) {
			return "", nil
		}
		return "", err
	}
	result := go_format.FormatInstance.ToString(target)
	return result, nil
}

func (c *ConfigManager) MustInt(str string) int {
	result, err := c.Int(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) Int(str string) (int, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		var notExistError *NotExistError
		if errors.As(err, &notExistError) {
			return 0, nil
		}
		return 0, err
	}
	result, err := go_format.FormatInstance.ToInt(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *ConfigManager) MustInt64(str string) int64 {
	result, err := c.Int64(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) Int64(str string) (int64, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		var notExistError *NotExistError
		if errors.As(err, &notExistError) {
			return 0, nil
		}
		return 0, err
	}
	result, err := go_format.FormatInstance.ToInt64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *ConfigManager) MustUint64(str string) uint64 {
	result, err := c.Uint64(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) Uint64(str string) (uint64, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		var notExistError *NotExistError
		if errors.As(err, &notExistError) {
			return 0, nil
		}
		return 0, err
	}
	result, err := go_format.FormatInstance.ToUint64(target)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *ConfigManager) MustBool(str string) bool {
	result, err := c.Bool(str)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ConfigManager) Bool(str string) (bool, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		var notExistError *NotExistError
		if errors.As(err, &notExistError) {
			return false, nil
		}
		return false, err
	}
	result, err := go_format.FormatInstance.ToBool(target)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (c *ConfigManager) Float64(str string) (float64, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		var notExistError *NotExistError
		if errors.As(err, &notExistError) {
			return 0.0, nil
		}
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

func (c *ConfigManager) FlagSetConfigs() map[string]interface{} {
	return c.flagSetConfigs
}

func (c *ConfigManager) MustMap(str string) map[string]interface{} {
	map_, err := c.Map(str)
	if err != nil {
		panic(err)
	}
	return map_
}

func (c *ConfigManager) Map(str string) (map[string]interface{}, error) {
	target, err := c.FindTarget(str)
	if err != nil {
		var notExistError *NotExistError
		if errors.As(err, &notExistError) {
			return nil, nil
		}
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
