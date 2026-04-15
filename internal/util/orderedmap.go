package util

// OrderedMap 定义泛型结构体
// K 必须是可比较类型 (int, string 等)
// V 可以是任何类型
type OrderedMap[K comparable, V any] struct {
	data map[K]V
	keys []K
}

// NewOrderedMap 初始化泛型实例
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		data: make(map[K]V),
		keys: make([]K, 0),
	}
}

// Set 插入或更新
func (om *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := om.data[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.data[key] = value
}

// Get 获取
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	val, ok := om.data[key]
	return val, ok
}

// Range 顺序遍历
func (om *OrderedMap[K, V]) Range(f func(key K, value V)) {
	for _, k := range om.keys {
		f(k, om.data[k])
	}
}
