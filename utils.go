package whiteboard

import "reflect"

func IsMapKeyTypeEqual(m reflect.Value, key reflect.Value) bool {
	// 判断 m 是否为 map 类型
	if m.Kind() != reflect.Map {
		return false
	}

	// 获取 map 的 key 类型
	mapKeyType := m.Type().Key()

	// 判断 key 的类型是否与 map 的 key 类型一致
	return key.Type().AssignableTo(mapKeyType)
}

func IsValidMatch(v reflect.Value, key reflect.Value) bool {
	// 判断 m 是否为 map 类型
	if v.Kind() == reflect.Map {
		// 获取 map 的 key 类型
		mapKeyType := v.Type().Key()

		// 判断 key 的类型是否与 map 的 key 类型一致
		return key.Type().AssignableTo(mapKeyType)
	}
	// v 是一个 array
	if v.Kind() == reflect.Array {

		// 判断 key 是否为有效的下标
		keyValue := reflect.ValueOf(key)
		keyInt := int(keyValue.Int())
		if keyValue.IsValid() && keyValue.Type().Kind() == reflect.Int && keyInt >= 0 && keyInt < v.Len() {
			// 获取 keyInt 下标的元素值
			// elemValue := v.Index(keyInt)
			// if elemValue.IsValid() {
			// 	return true
			// }
			return true
		}
	}

	return false
}
