package go_config

import (
	"flag"
	go_format "github.com/pefish/go-format"
	go_test_ "github.com/pefish/go-test"
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
	go_test_.Equal(t, `a2`, a.Haha)

	b := make([]uint64, 0)
	ConfigManagerInstance.MustGet(`/test3/test2/test8`, &b)
	go_test_.Equal(t, 3, len(b))

	c := make([]string, 0)
	ConfigManagerInstance.MustGet(`test4`, &c)
	go_test_.Equal(t, 2, len(c))
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
	go_test_.Equal(t, false, err == nil)
	go_test_.Equal(t, true, strings.Contains(err.Error(), "not exist"))
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
	result, err := configManagerInstance.Bool("abc")
	go_test_.Equal(t, nil, err)
	go_test_.Equal(t, true, result)

	result1, err1 := configManagerInstance.String("name")
	go_test_.Equal(t, nil, err1)
	go_test_.Equal(t, "_example", result1)

	result2, ok := configManagerInstance.FlagSetDefaultConfigs()["abc"]
	go_test_.Equal(t, true, ok)
	go_test_.Equal(t, true, go_format.FormatInstance.MustToBool(result2))
}

func TestConfigClass_GetString2(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := ConfigManagerInstance.MustString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = ConfigManagerInstance.MustString(`/test/haha`)
	//fmt.Println(str1, `  cache`)
	if str != `a2` {
		t.Error()
	}

	str = ConfigManagerInstance.MustString(`/test2`)
	if str != `b` {
		t.Error()
	}

	str = ConfigManagerInstance.MustString(`/test3/test2/test4/test5/test6`)
	if str != `122` {
		t.Error()
	}

	str = ConfigManagerInstance.MustString(`/test3/test2/test4/test5/test7`)
	if str != `11` {
		t.Error()
	}
}

func TestConfigClass_GetString3(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := ConfigManagerInstance.MustString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = ConfigManagerInstance.MustString(`/test/haha`)
	if str != `a2` {
		t.Error()
	}

	str = ConfigManagerInstance.MustString(`/test2`)
	if str != `b` {
		t.Error()
	}

	str = ConfigManagerInstance.MustString(`/test3/test2/test4/test5/test6`)
	if str != `122` {
		t.Error()
	}

	str = ConfigManagerInstance.MustString(`/test3/test2/test4/test5/test7`)
	if str != `11` {
		t.Error()
	}
}

func TestConfigClass_GetInt2(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	int_ := ConfigManagerInstance.MustInt(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetInt642(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	int_ := ConfigManagerInstance.MustInt64(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetMap(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	map_ := ConfigManagerInstance.MustMap(`/test3/test2`)
	go_test_.Equal(t, 45, go_format.FormatInstance.MustToInt(map_[`test3`]))

	map1_ := ConfigManagerInstance.MustMap(`/test3/test225235`)
	go_test_.Equal(t, map[string]interface{}(nil), map1_)
}

func TestConfigClass_GetBool(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	r := ConfigManagerInstance.MustBool(`xixi`)
	go_test_.Equal(t, false, r)
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
			if got := this.MustString(tt.args.str); got != tt.want {
				t.Errorf("ConfigManager.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigClass_MustGetStringDefault(t *testing.T) {
	ConfigManagerInstance.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := ConfigManagerInstance.MustString(`/test1/test2/test4577`)
	go_test_.Equal(t, "", str)
}
