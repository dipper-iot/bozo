package cli

import (
	"os"
	"reflect"
	"testing"
)

func TestStringFlag_GetName(t *testing.T) {
	type fields struct {
		Name         string
		Aliases      []string
		Usage        string
		Env          []string
		DefaultValue string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test",
			fields: fields{
				Name: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringFlag{
				Name:         tt.fields.Name,
				Aliases:      tt.fields.Aliases,
				Usage:        tt.fields.Usage,
				Env:          tt.fields.Env,
				DefaultValue: tt.fields.DefaultValue,
			}
			if got := s.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringFlag_GetUsage(t *testing.T) {
	type fields struct {
		Name         string
		Aliases      []string
		Usage        string
		Env          []string
		DefaultValue string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test",
			fields: fields{
				Name:  "test",
				Usage: "test",
			},
			want: "-test - test",
		},
		{
			name: "test",
			fields: fields{
				Name:         "test",
				Usage:        "test",
				DefaultValue: "6",
			},
			want: "-test - test (Default: 6)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringFlag{
				Name:         tt.fields.Name,
				Aliases:      tt.fields.Aliases,
				Usage:        tt.fields.Usage,
				Env:          tt.fields.Env,
				DefaultValue: tt.fields.DefaultValue,
			}
			if got := s.GetUsage(); got != tt.want {
				t.Errorf("GetUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringFlag_Setup(t *testing.T) {
	flag := StringFlag{
		Name:         "test",
		DefaultValue: "1",
	}
	flag2 := StringFlag{
		Name:         "test",
		Env:          []string{"TEST_ENV"},
		DefaultValue: "1",
	}
	os.Setenv("TEST_ENV", "1")
	type args struct {
		c *Context
	}
	c := &Context{
		args: []string{},
		dataFlagVal: map[string]string{
			"test": "1",
		},
		flagVal: map[string]FlagGet{},
	}
	c.flagVal[flag.GetName()] = flag.Setup(c)

	tests := []struct {
		name   string
		fields StringFlag
		args   args
		want   string
	}{
		{
			name:   "test",
			fields: flag,
			args: args{
				c: c,
			},
			want: "1",
		},
		{
			name:   "test2",
			fields: flag2,
			args: args{
				c: c,
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringFlag{
				Name:         tt.fields.Name,
				Aliases:      tt.fields.Aliases,
				Usage:        tt.fields.Usage,
				Env:          tt.fields.Env,
				DefaultValue: tt.fields.DefaultValue,
			}
			rs := s.Setup(tt.args.c)()
			if !reflect.DeepEqual(rs, tt.want) {
				t.Errorf("Setup() = %s, want %s", rs, tt.want)
			}
		})
	}
}
