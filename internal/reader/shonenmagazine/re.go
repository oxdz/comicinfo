package shonenmagazine

// 定义坐标点结构体
type Point struct {
	X int
	Y int
}

// 定义转换映射结构体
type Mapping struct {
	Source Point
	Dest   Point
}

func le(seed int) []*Mapping {
	mp := []*Mapping{}
	a := 4

	// 创建初始数组 [0, 1, 2, ..., a²-1]
	size := a * a
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = i
	}

	// 调用re函数进行处理
	processed := re(arr, uint32(seed))

	// 生成映射并发送到通道
	for t, o := range processed {
		mp = append(mp, &Mapping{
			Source: Point{
				X: o % a,
				Y: o / a,
			},
			Dest: Point{
				X: t % a,
				Y: t / a,
			},
		})
	}
	return mp
}

// re 函数实现
func re(a []int, e uint32) []int {
	s := ve(e, len(a))
	pairs := make([][2]uint32, len(a))

	// 生成随机数与原始值的对
	for i, t := range a {
		pairs[i] = [2]uint32{s[i], uint32(t)}
	}

	// 按随机数排序
	sortPairs(pairs)

	// 提取排序后的原始值
	result := make([]int, len(a))
	for i, pair := range pairs {
		result[i] = int(pair[1])
	}

	return result
}

// sortPairs 对二维数组按第一元素排序
func sortPairs(pairs [][2]uint32) {
	// 简单的冒泡排序实现
	n := len(pairs)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if pairs[j][0] > pairs[j+1][0] {
				pairs[j], pairs[j+1] = pairs[j+1], pairs[j]
			}
		}
	}
}

// ve 函数实现（返回一个生成随机数的函数）
func ve(state uint32, l int) (r []uint32) {
	for i := 0; i < l; i++ {
		state ^= state << 13
		state ^= state >> 17
		state ^= state << 5
		r = append(r, state)
	}
	return r
}

// Ue 计算符合条件的宽度和高度
func BlockSize(weight, high, s int) *struct {
	Width  int
	Height int
} {
	const y = 8 // 根据上下文可能需要调整
	// 检查输入是否满足条件
	if weight < s*y || high < s*y {
		return nil
	}

	// 计算中间值
	o := weight / y // 等价于JavaScript的Math.floor(a / y)
	t := high / y   // 等价于JavaScript的Math.floor(e / y)
	u := o / s      // 等价于JavaScript的Math.floor(o / s)
	r := t / s      // 等价于JavaScript的Math.floor(t / s)

	// 返回结果结构体
	return &struct {
		Width  int
		Height int
	}{
		Width:  u * y,
		Height: r * y,
	}
}
