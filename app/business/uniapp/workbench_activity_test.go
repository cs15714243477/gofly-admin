package uniapp

import "testing"

func TestWxIsWorkbenchActivityType(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{in: "showing", want: true},
		{in: "view", want: true},
		{in: "share", want: true},
		{in: "call", want: true},
		{in: "unlock", want: false}, // 开锁单独走 TTLock 相关接口
		{in: "follow", want: false}, // 关注记录走 business_favorites
		{in: "", want: false},
		{in: "  view  ", want: true},
		{in: "VIEW", want: false},
	}
	for _, c := range cases {
		if got := wxIsWorkbenchActivityType(c.in); got != c.want {
			t.Fatalf("wxIsWorkbenchActivityType(%q)=%v, want=%v", c.in, got, c.want)
		}
	}
}

