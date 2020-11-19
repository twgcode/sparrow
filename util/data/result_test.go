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

func TestRequestParaErrJson(t *testing.T) {
	type test struct { // 定义test结构体
		input string
		want  string
	}
	tests := map[string]test{
		"empty":   {"", "参数错误"},
		"english": {"english", "参数错误, english"},
		"chinese": {"中文", "参数错误, 中文"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := RequestParaErrJson(tc.input)
			if result.Message != tc.want {
				t.Errorf("name:%s excepted:%s, got:%s", name, tc.want, result.Message)
			}
		})
	}

}
