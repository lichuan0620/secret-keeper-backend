/*
Copyright 2018 Caicloud Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"reflect"
	"testing"
	"time"
)

func TestConverterFor(t *testing.T) {
	wantTime, _ := time.Parse(time.RFC3339, "2020-08-25T05:12:18Z")
	tests := []*struct {
		tpy         reflect.Type
		invalidType bool
		data        []string
		invalidData bool
		want        interface{}
		pointer     bool
	}{
		{
			tpy:         reflect.TypeOf(struct{}{}),
			invalidType: true,
		},
		{
			tpy:  reflect.TypeOf(false),
			data: []string{"true", "false"},
			want: true,
		},
		{
			tpy:  reflect.TypeOf(false),
			data: []string{"0"},
			want: false,
		},
		{
			tpy:         reflect.TypeOf(false),
			data:        []string{"00"},
			invalidData: true,
		},
		{
			tpy:  reflect.TypeOf(0),
			data: []string{"1", "2"},
			want: 1,
		},
		{
			tpy:  reflect.TypeOf(int8(0)),
			data: []string{"1", "2"},
			want: int8(1),
		},
		{
			tpy:  reflect.TypeOf(int32(0)),
			data: []string{"1", "2"},
			want: int32(1),
		},
		{
			tpy:  reflect.TypeOf(int64(0)),
			data: []string{"1", "2"},
			want: int64(1),
		},
		{
			tpy:  reflect.TypeOf(uint(0)),
			data: []string{"1", "2"},
			want: uint(1),
		},
		{
			tpy:  reflect.TypeOf(uint8(0)),
			data: []string{"1", "2"},
			want: uint8(1),
		},
		{
			tpy:  reflect.TypeOf(uint16(0)),
			data: []string{"1", "2"},
			want: uint16(1),
		},
		{
			tpy:  reflect.TypeOf(uint32(0)),
			data: []string{"1", "2"},
			want: uint32(1),
		},
		{
			tpy:  reflect.TypeOf(uint64(0)),
			data: []string{"1", "2"},
			want: uint64(1),
		},
		{
			tpy:  reflect.TypeOf(float32(0)),
			data: []string{"1.2", "2"},
			want: float32(1.2),
		},
		{
			tpy:  reflect.TypeOf(float64(0)),
			data: []string{"1.2", "2"},
			want: float64(1.2),
		},
		{
			tpy:  reflect.TypeOf(""),
			data: []string{"1", "2"},
			want: "1",
		},
		{
			tpy:  reflect.TypeOf(time.Time{}),
			data: []string{"2020-08-25T05:12:18Z"},
			want: wantTime,
		},
		{
			tpy:  reflect.TypeOf(time.Duration(0)),
			data: []string{"5m3s", "24h"},
			want: 5*time.Minute + 3*time.Second,
		},
		{
			tpy:         reflect.TypeOf(time.Duration(0)),
			data:        []string{"53"},
			invalidData: true,
		},
		{
			tpy:     reflect.TypeOf(new(time.Duration)),
			data:    []string{"5m3s", "24h"},
			want:    5*time.Minute + 3*time.Second,
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(bool)),
			data:    []string{"true", "2"},
			want:    true,
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(int)),
			data:    []string{"1", "2"},
			want:    1,
			pointer: true,
		},
		{
			tpy:         reflect.TypeOf(new(int)),
			data:        []string{"1.2"},
			invalidData: true,
		},
		{
			tpy:     reflect.TypeOf(new(int8)),
			data:    []string{"1", "2"},
			want:    int8(1),
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(int16)),
			data:    []string{"1", "2"},
			want:    int16(1),
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(int32)),
			data:    []string{"1", "2"},
			want:    int32(1),
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(int64)),
			data:    []string{"1", "2"},
			want:    int64(1),
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(uint)),
			data:    []string{"1", "2"},
			want:    uint(1),
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(uint8)),
			data:    []string{"1", "2"},
			want:    uint8(1),
			pointer: true,
		},

		{
			tpy:     reflect.TypeOf(new(uint16)),
			data:    []string{"1", "2"},
			want:    uint16(1),
			pointer: true,
		},

		{
			tpy:     reflect.TypeOf(new(uint32)),
			data:    []string{"1", "2"},
			want:    uint32(1),
			pointer: true,
		},

		{
			tpy:     reflect.TypeOf(new(uint64)),
			data:    []string{"1", "2"},
			want:    uint64(1),
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(float32)),
			data:    []string{"1.2", "2"},
			want:    float32(1.2),
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(float64)),
			data:    []string{"1.2", "2"},
			want:    float64(1.2),
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(string)),
			data:    []string{"1.2", "2"},
			want:    "1.2",
			pointer: true,
		},
		{
			tpy:     reflect.TypeOf(new(time.Time)),
			data:    []string{"2020-08-25T05:12:18Z", "2020-08-25T05:12:18Z"},
			want:    wantTime,
			pointer: true,
		},
		{
			tpy:  reflect.TypeOf([]bool{}),
			data: []string{"true", "false"},
			want: []bool{true, false},
		},
		{
			tpy:  reflect.TypeOf([]int{}),
			data: []string{"1", "2"},
			want: []int{1, 2},
		},
		{
			tpy:  reflect.TypeOf([]float64{}),
			data: []string{"1.2", "2.2"},
			want: []float64{1.2, 2.2},
		},
		{
			tpy:  reflect.TypeOf([]string{}),
			data: []string{"1.2", "2.2"},
			want: []string{"1.2", "2.2"},
		},
	}
	for _, tc := range tests {
		t.Run("", func(tt *testing.T) {
			converter, err := ConverterFor(tc.tpy)
			if tc.invalidType == (err == nil) {
				tt.Fatalf("expecting invalid type: %v; got converter error: %v", tc.invalidType, err)
			}
			if tc.invalidType {
				return
			}

			got, err := converter(tc.data)
			if tc.invalidData == (err == nil) {
				tt.Fatalf("expecting invalid data: %v; got convertor error: %v", tc.invalidData, err)
			}
			if tc.invalidData {
				return
			}

			value := got
			if tc.pointer {
				value = reflect.ValueOf(got).Elem().Interface()
			}
			if !reflect.DeepEqual(value, tc.want) {
				tt.Fatalf("convertor for type %v, expecting %v, got %v", tc.tpy, tc.want, value)
			}
		})
	}

}
