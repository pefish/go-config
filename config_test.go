package go_config

import (
	"flag"
	go_reflect "github.com/pefish/go-reflect"
	"github.com/pefish/go-test-assert"
	"strings"
	"testing"
)

func TestConfigClass_LoadConfig(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		SecretFilepath: `./_example/test.yaml`,
	})
	a := struct {
		Haha string `json:"haha"`
	}{}
	ConfigManagerInstance.MustGet(`test`, &a)
	test.Equal(t, `a2`, a.Haha)

	b := make([]uint64, 0)
	ConfigManagerInstance.MustGet(`/test3/test2/test8`, &b)
	test.Equal(t, 3, len(b))

	c := make([]string, 0)
	ConfigManagerInstance.MustGet(`test4`, &c)
	test.Equal(t, 2, len(c))
}

func TestConfigClass_LoadYamlConfig1(t *testing.T) {
	instance := NewConfigManager()
	err := instance.LoadConfig(Configuration{
		ConfigFilepath: "",
	})
	if err != nil {
		t.Error(err)
	}
	a := struct {
		Haha string `json:"haha"`
	}{}
	err = instance.Get(`test`, &a)
	test.Equal(t, false, err == nil)
	test.Equal(t, true, strings.Contains(err.Error(), "not exist"))
}

func TestConfigClass_LoadYamlConfig3(t *testing.T) {
	configManagerInstance := NewConfigManager()

	flagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagSet.String("config", "", "path to config file")
	flagSet.String("name", "pefish", "listener name")
	flagSet.Bool("abc", true, "abc")

	configManagerInstance.MustLoadConfig(Configuration{
		SecretFilepath: `./_example/test.yaml`,
	})
	configManagerInstance.MergeFlagSet(flagSet)
	result, err := configManagerInstance.GetBool("abc")
	test.Equal(t, nil, err)
	test.Equal(t, true, result)

	result1, err1 := configManagerInstance.GetString("name")
	test.Equal(t, nil, err1)
	test.Equal(t, "_example", result1)

	result2, ok := configManagerInstance.FlagSetDefaultConfigs()["abc"]
	test.Equal(t, true, ok)
	test.Equal(t, true, go_reflect.Reflect.MustToBool(result2))
}

func TestConfigClass_GetString2(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := ConfigManagerInstance.MustGetString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = ConfigManagerInstance.MustGetString(`/test/haha`)
	//fmt.Println(str1, `  cache`)
	if str != `a2` {
		t.Error()
	}

	str = ConfigManagerInstance.MustGetString(`/test2`)
	if str != `b` {
		t.Error()
	}

	str = ConfigManagerInstance.MustGetString(`/test3/test2/test4/test5/test6`)
	if str != `122` {
		t.Error()
	}

	str = ConfigManagerInstance.MustGetString(`/test3/test2/test4/test5/test7`)
	if str != `11` {
		t.Error()
	}
}

func TestConfigClass_GetString3(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := ConfigManagerInstance.MustGetString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = ConfigManagerInstance.MustGetString(`/test/haha`)
	if str != `a2` {
		t.Error()
	}

	str = ConfigManagerInstance.MustGetString(`/test2`)
	if str != `b` {
		t.Error()
	}

	str = ConfigManagerInstance.MustGetString(`/test3/test2/test4/test5/test6`)
	if str != `122` {
		t.Error()
	}

	str = ConfigManagerInstance.MustGetString(`/test3/test2/test4/test5/test7`)
	if str != `11` {
		t.Error()
	}
}

func TestConfigClass_GetInt2(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	int_ := ConfigManagerInstance.MustGetInt(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetInt642(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	int_ := ConfigManagerInstance.MustGetInt64(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetMap(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	map_ := ConfigManagerInstance.MustGetMapDefault(`/test3/test2`, nil)
	test.Equal(t, 45, go_reflect.Reflect.MustToInt(map_[`test3`]))

	map1_ := ConfigManagerInstance.MustGetMapDefault(`/test3/test225235`, map[string]interface{}{
		"haha111": "36573",
	})
	if map1_[`haha111`].(string) != "36573" {
		t.Error()
	}
}

func TestConfigClass_GetBool(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	_, err := ConfigManagerInstance.GetBool(`xixi`)
	if _, ok := err.(*NotExistError); !ok {
		t.Error()
	}
}

func TestConfigClass_GetString(t *testing.T) {
	type fields struct {
		configs map[string]interface{}
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: `_example GetString`,
			fields: fields{
				map[string]interface{}{
					`_example`: `haha`,
				},
			},
			args: args{
				`_example`,
			},
			want: `haha`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ConfigManager{
				configs: tt.fields.configs,
			}
			if got := this.MustGetString(tt.args.str); got != tt.want {
				t.Errorf("ConfigManager.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigClass_MustGetStringDefault(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := ConfigManagerInstance.MustGetStringDefault(`/test1/test2/test4577`, `123`)
	if str != `123` {
		t.Error()
	}
}