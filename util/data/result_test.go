/**
@Author: wei-g
@Date:   2020/3/26 3:13 下午
@Description: 简介
*/

package data

import "testing"

func TestCode(t *testing.T) {
	t.Log(SucCode)
	t.Log(FailedCode)
}
func TestResponseInternalErrJson(t *testing.T) {
	ret := ResponseInternalErrJson("test")
	t.Log(ret)
	t.Log(ret.Data)
	ret1 := ResponseInternalErrJson("test", map[string]int{"a": 1})
	t.Log(ret1)
	t.Log(ret1.Data)
}
