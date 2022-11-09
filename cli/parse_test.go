package cli

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		a    *App
		args []string
	}
	tests := []struct {
		name           string
		args           args
		wantAppName    string
		wantArgsResult []string
		wantFlagResult map[string]string
		wantErr        bool
	}{
		{
			name: "test1",
			args: args{
				a:    nil,
				args: []string{"app"},
			},
			wantAppName:    "app",
			wantArgsResult: []string{},
			wantErr:        false,
			wantFlagResult: map[string]string{},
		},
		{
			name: "test2",
			args: args{
				a:    nil,
				args: []string{"app", "-h"},
			},
			wantAppName:    "app",
			wantArgsResult: []string{},
			wantErr:        false,
			wantFlagResult: map[string]string{
				"h": "",
			},
		},
		{
			name: "test3",
			args: args{
				a:    nil,
				args: []string{"app", "-h", "-http=10", "-web", "20"},
			},
			wantAppName:    "app",
			wantArgsResult: []string{},
			wantErr:        false,
			wantFlagResult: map[string]string{
				"h":    "",
				"http": "10",
				"web":  "20",
			},
		},
		{
			name: "test4",
			args: args{
				a:    nil,
				args: []string{"app", "-h", "-http=10", "-web", "20", "command1"},
			},
			wantAppName: "app",
			wantArgsResult: []string{
				"command1",
			},
			wantErr: false,
			wantFlagResult: map[string]string{
				"h":    "",
				"http": "10",
				"web":  "20",
			},
		},
		{
			name: "test5",
			args: args{
				a:    nil,
				args: []string{"app", "-h", "-http=10", "-web", "20", "command1", "command2"},
			},
			wantAppName: "app",
			wantArgsResult: []string{
				"command1",
				"command2",
			},
			wantErr: false,
			wantFlagResult: map[string]string{
				"h":    "",
				"http": "10",
				"web":  "20",
			},
		},
		{
			name: "test6",
			args: args{
				a:    nil,
				args: []string{"app", "command1", "sub", "1", "-o"},
			},
			wantAppName: "app",
			wantArgsResult: []string{
				"command1",
				"sub",
				"1",
			},
			wantErr: false,
			wantFlagResult: map[string]string{
				"o": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAppName, gotArgsResult, gotFlagResult, err := Parse(tt.args.a, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAppName != tt.wantAppName {
				t.Errorf("Parse() gotAppName = %v, want %v", gotAppName, tt.wantAppName)
			}
			if !reflect.DeepEqual(gotArgsResult, tt.wantArgsResult) {
				t.Errorf("Parse() gotArgsResult = %v, want %v", gotArgsResult, tt.wantArgsResult)
			}
			if !reflect.DeepEqual(gotFlagResult, tt.wantFlagResult) {
				t.Errorf("Parse() gotFlagResult = %v, want %v", gotFlagResult, tt.wantFlagResult)
			}
		})
	}
}
