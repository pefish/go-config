package go_config

import (
	"fmt"
	"testing"
)

func TestConfigClass_LoadYamlConfig(t *testing.T) {
	Config.MustLoadYamlConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-core/_example/config/local.yaml`,
		SecretFilepath: `/Users/joy/Work/backend/go-core/_example/secret/local1.yaml`,
	})
	a := struct {
		Host string `json:"host"`
	}{}
	Config.MustGetStruct(`mysql`, &a)
	fmt.Println(a)
}

func TestConfigClass_GetString2(t *testing.T) {
	Config.MustLoadYamlConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-config/test/test.yaml`,
	})
	str := Config.GetString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = Config.GetString(`/test/haha`)
	str1 := Config.GetString(`/test/haha`)
	fmt.Println(str1, `  cache`)
	if str != `a` {
		t.Error()
	}

	str = Config.GetString(`/test2`)
	if str != `b` {
		t.Error()
	}

	str = Config.GetString(`/test3/test2/test4/test5/test6`)
	if str != `122` {
		t.Error()
	}

	str = Config.GetString(`/test3/test2/test4/test5/test7`)
	if str != `11` {
		t.Error()
	}
}

func TestConfigClass_GetString3(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-config/test/test.json`,
	})
	str := Config.GetString(`/test1/test2/test3`)
	if str != `45` {
		t.Error()
	}

	str = Config.GetString(`/test/haha`)
	if str != `a` {
		t.Error()
	}

	str = Config.GetString(`/test2`)
	if str != `b` {
		t.Error()
	}

	str = Config.GetString(`/test3/test2/test4/test5/test6`)
	if str != `122` {
		t.Error()
	}

	str = Config.GetString(`/test3/test2/test4/test5/test7`)
	if str != `11` {
		t.Error()
	}
}

func TestConfigClass_GetInt2(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-config/test/test.json`,
	})
	int_ := Config.GetInt(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetInt642(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-config/test/test.json`,
	})
	int_ := Config.GetInt64(`/test1/test2/test3`)
	if int_ != 45 {
		t.Error()
	}
}

func TestConfigClass_GetMap(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-config/test/test.json`,
	})
	map_ := Config.MustGetMap(`/test3/test2`)
	if map_[`test3`].(float64) != 45 {
		t.Error()
	}
}

func TestConfigClass_GetSlice(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-config/test/test.json`,
	})
	slice_ := Config.MustGetSlice(`/test3/test2/test8`)
	if slice_[0].(float64) != 1 {
		t.Error()
	}
}

func TestConfigClass_GetBool(t *testing.T) {
	Config.MustLoadJsonConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-config/test/test.json`,
	})
	fmt.Println(Config.GetBool(`xixi`))
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
			if got := this.GetString(tt.args.str); got != tt.want {
				t.Errorf("ConfigClass.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}
