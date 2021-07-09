package logger

import (
	"reflect"
	"testing"
)

func TestKV_Key(t *testing.T) {
	type fields struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should pass",
			fields: fields{
				key:   "some key",
				value: "some value",
			},
			want: "some key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KV{
				key:   tt.fields.key,
				value: tt.fields.value,
			}
			if got := k.Key(); got != tt.want {
				t.Errorf("KV.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKV_Value(t *testing.T) {
	type fields struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "should pass",
			fields: fields{
				key:   "some key",
				value: "some value",
			},
			want: "some value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KV{
				key:   tt.fields.key,
				value: tt.fields.value,
			}
			if got := k.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KV.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_ToFields(t *testing.T) {
	tests := []struct {
		name string
		m    Map
		want []Field
	}{
		{
			name: "should pass",
			m: Map{
				"test": 999,
			},
			want: []Field{
				KV{"test", 999},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ToFields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.ToFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
