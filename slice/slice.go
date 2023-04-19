package slice

type DiffIntersect[T1, T2 comparable] interface {
	// Diff
	//  @Description: 差集
	//  @param s1
	//  @param s2
	//  @return []T1 返回s1中存在，s2中不存在的元素
	Diff(s1 []T1, s2 []T2) []T1

	// Intersect
	//  @Description: 交集
	//  @param s1
	//  @param s2
	//  @return []T1 返回s1和s2中均存在的元素
	Intersect(s1 []T1, s2 []T2) []T1
}

type twoSlice[T1, T2 comparable] struct {
	fn CompareFunc[T1, T2]
}

func (t *twoSlice[T1, T2]) Diff(s1 []T1, s2 []T2) []T1 {
	var s3 []T1
loop:
	for _, item1 := range s1 {
		for _, item2 := range s2 {
			if t.fn(item1, item2) {
				continue loop
			}
		}
		s3 = append(s3, item1)
	}
	return s3
}

func (t *twoSlice[T1, T2]) Intersect(s1 []T1, s2 []T2) []T1 {
	var s3 []T1
loop:
	for _, item1 := range s1 {
		for _, item2 := range s2 {
			if t.fn(item1, item2) {
				s3 = append(s3, item1)
				continue loop
			}
		}
	}
	return s3
}

type CompareFunc[T1, T2 comparable] func(T1, T2) bool

// NewCompareHelper
//
//	@Description: 新建一个
//	@param fn
//	@param T2]
//	@return DiffIntersect[T1
//	@return T2]
func NewCompareHelper[T1, T2 comparable](fn CompareFunc[T1, T2]) DiffIntersect[T1, T2] {
	return &twoSlice[T1, T2]{fn: fn}
}

// Map
//
//	@Description: 创建一个新切片，这个新切片由原切片中的每个元素都调用一次提供的函数后的返回值组成。
//	@param s 原始切片
//	@param fn
//	@return []T2
func Map[T1, T2 comparable](s []T1, fn func(element T1) T2) []T2 {
	var s2 = make([]T2, 0, len(s))
	for _, item := range s {
		s2 = append(s2, fn(item))
	}
	return s2
}

// Filter
//
//	@Description: 创建给定切片一部分的浅拷贝，其包含通过所提供函数实现的测试的所有元素。
//	@param s 原始切片
//	@param fn
//	@return []T1
func Filter[T comparable](s []T, fn func(element T) bool) []T {
	var s2 = make([]T, 0, len(s))
	for _, item := range s {
		if fn(item) {
			s2 = append(s2, item)
		}
	}
	return s2
}

// MapFilter
//
//	@Description: 创建一个新切片，这个新切片由原切片中的每个元素都调用一次提供的函数后的符合条件的返回值组成。
//	@param s 原始切片
//	@param fn 处理函数，将T1转变为T2，只有fn返回true的T2类型的值会返回
//	@return []T2
func MapFilter[T1, T2 comparable](s []T1, fn func(element T1) (T2, bool)) []T2 {
	var s2 = make([]T2, 0, len(s))
	for _, item := range s {
		if item2, ok := fn(item); ok {
			s2 = append(s2, item2)
		}
	}
	return s2
}

// ForEach
//
//	@Description: 对切片的每个元素执行一次给定的函数。
//	@param s
//	@param fn func(element T1, index int)
//		element 数组中正在处理的当前元素。
//		index 	数组中正在处理的当前元素的索引。
func ForEach[T comparable](s []T, fn func(element T, index int)) {
	for i, item := range s {
		fn(item, i)
	}
}

// Some
//
//	@Description: 测试切片中是不是至少有 1 个元素通过了被提供的函数测试
//	@param s 原始切片
//	@param fn 条件函数
//	@return bool
func Some[T comparable](s []T, fn func(element T) bool) bool {
	for _, item := range s {
		if fn(item) {
			return true
		}
	}
	return false
}

// Every
//
//	@Description: 测试一个切片内的所有元素是否都能通过指定函数的测试。
//	@param s 原始切片
//	@param fn 条件函数
//	@return bool
func Every[T comparable](s []T, fn func(element T) bool) bool {
	for _, item := range s {
		if !fn(item) {
			return false
		}
	}
	return true
}
