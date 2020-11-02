package go_config

import (
	"flag"
	go_reflect "github.com/pefish/go-reflect"
	"strings"
	"testing"
	"github.com/pefish/go-test-assert"
)

func TestConfigClass_LoadConfig(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		SecretFilepath: `./_example/test.yaml`,
	})
	a := struct {
		Haha string `json:"haha"`
	}{}
	Config.MustGetStruct(`test`, &a)
	if a.Haha != `a2` {
		t.Error()
	}
}

func TestConfigClass_LoadYamlConfig1(t *testing.T) {
	err := Config.LoadConfig(Configuration{
		ConfigFilepath: "",
	})
	if err != nil {
		t.Error(err)
	}
	a := struct {
		Haha string `json:"haha"`
	}{}
	err = Config.GetStruct(`test`, &a)
	if err == nil || !strings.Contains(err.Error(), "not exist"){
		t.Error()
	}
}

func TestConfigClass_LoadYamlConfig3(t *testing.T) {
	flagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagSet.String("config", "", "path to config file")
	flagSet.String("name", "pefish", "listener name")
	flagSet.Bool("abc", true, "abc")

	Config.MustLoadConfig(Configuration{
		SecretFilepath: `./_example/test.yaml`,
	})
	Config.MergeFlagSet(flagSet)
	result, err := Config.GetBool("abc")
	test.Equal(t, nil, err)
	test.Equal(t, true, result)

	result1, err1 := Config.GetString("name")
	test.Equal(t, nil, err1)
	test.Equal(t, "_example", result1)

	result2, ok := Config.FlagSetConfigs()["abc"]
	test.Equal(t, true, ok)
	test.Equal(t, true, go_reflect.Reflect.MustToBool(result2))
}

func TestConfigClass_GetString2(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := Config.MustGetString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = Config.MustGetString(`/test/haha`)
	//fmt.Println(str1, `  cache`)
	if str != `a2` {
		t.Error()
	}

	str = Config.MustGetString(`/test2`)
	if str != `b` {
		t.Error()
	}

	str = Config.MustGetString(`/test3/test2/test4/test5/test6`)
	if str != `122` {
		t.Error()
	}

	str = Config.MustGetString(`/test3/test2/test4/test5/test7`)
	if str != `11` {
		t.Error()
	}
}

func TestConfigClass_GetString3(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := Config.MustGetString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = Config.MustGetString(`/test/haha`)
	if str != `a2` {
		t.Error()
	}

	str = Config.MustGetString(`/test2`)
	if str != `b` {
		t.Error()
	}

	str = Config.MustGetString(`/test3/test2/test4/test5/test6`)
	if str != `122` {
		t.Error()
	}

	str = Config.MustGetString(`/test3/test2/test4/test5/test7`)
	if str != `11` {
		t.Error()
	}
}

func TestConfigClass_GetInt2(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	int_ := Config.MustGetInt(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetInt642(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	int_ := Config.MustGetInt64(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetMap(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	map_ := Config.MustGetMapDefault(`/test3/test2`, nil)
	test.Equal(t, 45, go_reflect.Reflect.MustToInt(map_[`test3`]))

	map1_ := Config.MustGetMapDefault(`/test3/test225235`, map[string]interface{}{
		"haha111": "36573",
	})
	if map1_[`haha111`].(string) != "36573" {
		t.Error()
	}
}

func TestConfigClass_GetSlice(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	slice_ := Config.MustGetSliceDefault(`/test3/test2/test8`, nil)
	test.Equal(t, 1, go_reflect.Reflect.MustToInt(slice_[0]))
}

func TestConfigClass_GetBool(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	_, err := Config.GetBool(`xixi`)
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
			this := &ConfigClass{
				configs: tt.fields.configs,
			}
			if got := this.MustGetString(tt.args.str); got != tt.want {
				t.Errorf("ConfigClass.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigClass_MustGetStringDefault(t *testing.T) {
	Config.MustLoadConfig(Configuration{
		ConfigFilepath: `./_example/test.yaml`,
	})
	str := Config.MustGetStringDefault(`/test1/test2/test4577`, `123`)
	if str != `123` {
		t.Error()
	}
}