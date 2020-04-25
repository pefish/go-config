package go_config

import (
	"flag"
	"fmt"
	"testing"
)

func TestConfigClass_LoadYamlConfig(t *testing.T) {
	Config.MustLoadYamlConfig(Configuration{
		ConfigFilepath: `./test/test.json`,
		SecretFilepath: `./test/test.yaml`,
	})
	a := struct {
		Haha string `json:"haha"`
	}{}
	Config.MustGetStruct(`test`, &a)
	if a.Haha != `a` {
		t.Error()
	}
}

func TestConfigClass_LoadYamlConfig1(t *testing.T) {
	err := Config.LoadYamlConfig(Configuration{
		ConfigFilepath: "",
	})
	if err != nil {
		t.Error(err)
	}
	a := struct {
		Haha string `json:"haha"`
	}{}
	err = Config.GetStruct(`test`, &a)
	if err == nil || err.Error() != "not exist"{
		t.Error()
	}
}

func TestConfigClass_LoadYamlConfig3(t *testing.T) {
	flagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagSet.String("config", "", "path to config file")
	flagSet.String("name", "pefish", "listener name")
	flagSet.Bool("abc", true, "abc")

	Config.MustLoadYamlConfig(Configuration{
		ConfigFilepath: `./test/test.json`,
		SecretFilepath: `./test/test.yaml`,
	})
	Config.MergeFlagSet(flagSet)
	if result, err := Config.GetBool("abc"); err != nil || result != true {
		t.Error()
	}
}

func TestConfigClass_GetString2(t *testing.T) {
	Config.MustLoadYamlConfig(Configuration{
		ConfigFilepath: `./test/test.yaml`,
	})
	str := Config.MustGetString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = Config.MustGetString(`/test/haha`)
	str1 := Config.MustGetString(`/test/haha`)
	fmt.Println(str1, `  cache`)
	if str != `a` {
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
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `./test/test.json`,
	})
	str := Config.MustGetString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = Config.MustGetString(`/test/haha`)
	if str != `a` {
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
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `./test/test.json`,
	})
	int_ := Config.MustGetInt(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetInt642(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `./test/test.json`,
	})
	int_ := Config.MustGetInt64(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetMap(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `./test/test.json`,
	})
	map_ := Config.MustGetMap(`/test3/test2`)
	if map_[`test3`].(float64) != 45 {
		t.Error()
	}
}

func TestConfigClass_GetSlice(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `./test/test.json`,
	})
	slice_ := Config.MustGetSlice(`/test3/test2/test8`)
	if slice_[0].(float64) != 1 {
		t.Error()
	}
}

func TestConfigClass_GetBool(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `./test/test.json`,
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
			name: `test GetString`,
			fields: fields{
				map[string]interface{}{
					`test`: `haha`,
				},
			},
			args: args{
				`test`,
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
	Config.MustLoadYamlConfig(Configuration{
		ConfigFilepath: `./test/test.yaml`,
	})
	str := Config.MustGetStringDefault(`/test1/test2/test4577`, `123`)
	if str != `123` {
		t.Error()
	}
}