package utilx

func Mapping[T1 any, T2 any](slice []T1, mapperFunction func(T1) T2) []T2 {

	result := make([]T2, len(slice))

	for i, e := range slice {
		result[i] = mapperFunction(e)
	}

	return result
}

func MappingToMap[T1 any, T2 comparable, T3 any](slice []T1, mapperFunction func(T1) (T2, T3)) map[T2]T3 {

	result := make(map[T2]T3, len(slice))

	for _, e := range slice {
		key, value := mapperFunction(e)
		result[key] = value
	}

	return result
}

func MappingToMultiMap[T1 any, T2 comparable, T3 any](slice []T1, mapperFunction func(T1) (T2, T3)) map[T2][]T3 {

	result := make(map[T2][]T3, len(slice))

	for _, e := range slice {
		key, value := mapperFunction(e)
		values, ok := result[key]
		if !ok {
			values = make([]T3, 0, 4)
		}
		values = append(values, value)
		result[key] = values
	}

	return result
}
