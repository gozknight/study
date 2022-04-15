package test

import (
	"fmt"
	"gozknight.com/online-judge/pkg"
	"testing"
)

func TestLRUCache(t *testing.T) {
	lru := pkg.Constructor(10)
	for i := '0'; i < '9'; i++ {
		lru.Put(string(i), string(i))
	}
	v3 := lru.Get("1")
	fmt.Println(v3)
	v4 := lru.Get("9")
	fmt.Println(v4)
	v5 := lru.Get("8")
	fmt.Println(v5)
}
