package pkg

type SliceMapFunc[Input any, Result any] = func(item Input) Result

func SliceMap[Input any, Result any](slice []Input, callback SliceMapFunc[Input, Result]) []Result {
	result := make([]Result, len(slice))

	for i, item := range slice {
		result[i] = callback(item)
	}

	return result
}
