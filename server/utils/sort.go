package utils

type SortType interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | string
}

// 快速排序,right一般为数组长度
func QuickSort[T SortType](arr []T, left, right int) []T {
	if left < right {
		pivot := arr[left]
		j := left
		for i := left; i < right; i++ {
			if arr[i] < pivot {
				j++
				arr[j], arr[i] = arr[i], arr[j]
			}
		}
		arr[left], arr[j] = arr[j], arr[left]
		QuickSort(arr, left, j)
		QuickSort(arr, j+1, right)
	}
	return arr
}
