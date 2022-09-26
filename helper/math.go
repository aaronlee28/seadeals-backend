package helper

func SumInt(array []int) int {
	result := int(0)
	for _, v := range array {
		result += v
	}
	return result
}
