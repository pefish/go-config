package go_config

import (
	"testing"

	go_format "github.com/pefish/go-format"
	go_test_ "github.com/pefish/go-test"
)

func TestConfigClass_MergeConfigFile(t *testing.T) {
	err := ConfigManagerInstance.MergeConfigFile(`./test.yaml`)
	go_test_.Equal(t, nil, err)
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

func TestConfigClass_GetMap(t *testing.T) {
	err := ConfigManagerInstance.MergeConfigFile(`./test.yaml`)
	go_test_.Equal(t, nil, err)
	map_ := ConfigManagerInstance.MustMap(`/test3/test2`)
	go_test_.Equal(t, 45, go_format.FormatInstance.MustToInt(map_[`test3`]))

	map1_ := ConfigManagerInstance.MustMap(`/test3/test225235`)
	go_test_.Equal(t, map[string]interface{}(nil), map1_)
}

func TestConfigClass_GetBool(t *testing.T) {
	err := ConfigManagerInstance.MergeConfigFile(`./test.yaml`)
	go_test_.Equal(t, nil, err)
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
	err := ConfigManagerInstance.MergeConfigFile(`./test.yaml`)
	go_test_.Equal(t, nil, err)
	str := ConfigManagerInstance.MustString(`/test1/test2/test4577`)
	go_test_.Equal(t, "", str)
}
