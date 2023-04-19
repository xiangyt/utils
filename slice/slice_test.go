package slice

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestEvery(t *testing.T) {
	type args[T comparable] struct {
		s  []T
		fn func(element T) bool
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "全大于0",
			args: args[int]{
				s: []int{1, 2, 3},
				fn: func(element int) bool {
					return element > 0
				},
			},
			want: true,
		},
		{
			name: "部分大于0",
			args: args[int]{
				s: []int{1, 2, -3},
				fn: func(element int) bool {
					return element > 0
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Every(tt.args.s, tt.args.fn); got != tt.want {
				t.Errorf("Every() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type args[T comparable] struct {
		s  []T
		fn func(element T) bool
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "过滤大于0的元素",
			args: args[int]{
				s: []int{1, 2, -3},
				fn: func(element int) bool {
					return element > 0
				},
			},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.s, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForEach(t *testing.T) {
	type args[T comparable] struct {
		s  []T
		fn func(element T, index int)
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "对每个元素都执行一次函数",
			args: args[int]{
				s: []int{1, 2, -3},
				fn: func(element int, index int) {
					fmt.Println(element, index)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ForEach(tt.args.s, tt.args.fn)
		})
	}
}

func TestMap(t *testing.T) {
	type args[T1 comparable, T2 comparable] struct {
		s  []T1
		fn func(element T1) T2
	}
	type testCase[T1 comparable, T2 comparable] struct {
		name string
		args args[T1, T2]
		want []T2
	}
	tests := []testCase[int, string]{
		{
			name: "对每个元素执行一次函数，返回由函数返回值组成的新切片",
			args: args[int, string]{
				s: []int{1, 2, -3},
				fn: func(element int) string {
					return strconv.Itoa(element)
				},
			},
			want: []string{"1", "2", "-3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.s, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapFilter(t *testing.T) {
	type args[T1 comparable, T2 comparable] struct {
		s  []T1
		fn func(element T1) (T2, bool)
	}
	type testCase[T1 comparable, T2 comparable] struct {
		name string
		args args[T1, T2]
		want []T2
	}
	tests := []testCase[int, string]{
		{
			name: "同时执行map和filter",
			args: args[int, string]{
				s: []int{0, 1, 2, -3},
				fn: func(element int) (string, bool) {
					if element <= 0 {
						return "", false
					}
					return strconv.Itoa(element), true
				},
			},
			want: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapFilter(tt.args.s, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSome(t *testing.T) {
	type args[T comparable] struct {
		s  []T
		fn func(element T) bool
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "至少有一个大于0",
			args: args[int]{
				s: []int{1, 2, -3},
				fn: func(element int) bool {
					return element > 0
				},
			},
			want: true,
		},
		{
			name: "全部小于或等于0",
			args: args[int]{
				s: []int{0, -2, -3},
				fn: func(element int) bool {
					return element > 0
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Some(tt.args.s, tt.args.fn); got != tt.want {
				t.Errorf("Some() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_twoSlice_Diff(t1 *testing.T) {
	type args[T1 comparable, T2 comparable] struct {
		s1 []T1
		s2 []T2
	}
	type testCase[T1 comparable, T2 comparable] struct {
		name string
		t    twoSlice[T1, T2]
		args args[T1, T2]
		want []T1
	}
	tests := []testCase[int, int]{
		{
			name: "差集",
			t: twoSlice[int, int]{
				fn: func(i int, j int) bool {
					return i == j
				},
			},
			args: args[int, int]{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 4},
			},
			want: []int{3},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.t.Diff(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_twoSlice_Intersect(t1 *testing.T) {
	type args[T1 comparable, T2 comparable] struct {
		s1 []T1
		s2 []T2
	}
	type testCase[T1 comparable, T2 comparable] struct {
		name string
		t    twoSlice[T1, T2]
		args args[T1, T2]
		want []T1
	}
	tests := []testCase[int, int]{
		{
			name: "交集",
			t: twoSlice[int, int]{
				fn: func(i int, j int) bool {
					return i == j
				},
			},
			args: args[int, int]{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 4},
			},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.t.Intersect(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}
