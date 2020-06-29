/**
@Author: wei-g
@Date:   2020/6/29 4:02 下午
@Description:
*/

package sliceutil

// ContainsElement 查看某个元素是否在切片中
func ContainsElement(strList []string, element string) bool {
	for _, v := range strList {
		if v == element {
			return true
		}
	}
	return false
}
