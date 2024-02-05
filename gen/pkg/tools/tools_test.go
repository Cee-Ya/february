package tools

import "testing"

func TestFormatStructName(t *testing.T) {
	type args struct {
		prefix    string
		tableName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				prefix:    "t_",
				tableName: "t_sys_role",
			},
			want: "SysRole",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatStructName(tt.args.prefix, tt.args.tableName); got != tt.want {
				t.Errorf("FormatStructName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatJsonColumn(t *testing.T) {
	type args struct {
		tableName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				tableName: "t_sys_role",
			},
			want: "sysrole",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatJsonColumn("t_", tt.args.tableName); got != tt.want {
				t.Errorf("FormatJsonColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateCacheName(t *testing.T) {
	type args struct {
		tableName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				tableName: "User",
			},
			want: "cache::sys::role::",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateCacheName(tt.args.tableName); got != tt.want {
				t.Errorf("CreateCacheName() = %v, want %v", got, tt.want)
			}
		})
	}
}
