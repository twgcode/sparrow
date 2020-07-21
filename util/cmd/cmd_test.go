/**
@Author: wei-g
@Date:   2020/6/19 5:05 下午
@Description:
*/

package cmd

import (
	"fmt"
	"testing"
)

func TestInitCmd(t *testing.T) {
	err := InitCmd("sparrow", "", "", nil, "v0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(config)
}
