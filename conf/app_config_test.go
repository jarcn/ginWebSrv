package conf

import "testing"

func TestSetDefault(t *testing.T) {
	ret := setDefault("1", "1", "3")
	t.Log(ret)
	s := setDefault("1", "2", "3")
	t.Log(s)
}
