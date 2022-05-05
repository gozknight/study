package main

import (
	"bufio"
	. "fmt"
	"io"
	"os"
)

func B(r io.Reader, w io.Writer) {
	in := bufio.NewReader(r)
	out := bufio.NewWriter(w)
	defer out.Flush()
	var query int
	Fscan(in, &query)
	for ; query > 0; query-- {
		var t, s int
		var sum float64
		Fscan(in, &t, &s)
		sum = sum + 20*1.8
		if s > 20 {
			left := s - 20
			sum = sum + float64(left)*2.5
		}
		sum = sum + float64(t)*0.3
		if sum < 9 {
			Fprintf(out, "%.2f\n", 9.00)
		} else {
			Fprintf(out, "%.2f\n", sum)
		}
	}
}
func A(r io.Reader, w io.Writer) {
	in := bufio.NewReader(r)
	out := bufio.NewWriter(w)
	defer out.Flush()
	var ch string
	var v, t int
	var a, b int
	Fscan(in, &ch, &v, &t)
	if ch[0] == 'a' {
		t = t - 3 - v

		a = t / 2
		if t%3 == 0 {
			b = t / 3
		} else {
			b = t/3 + 1
		}

	} else if ch[0] == 'b' {
		t = t - 3 - 2*v
		a = t
		b = t/3 + t%3
	} else if ch[0] == 'c' {
		t = t - 3 - 3*v
		a = t
		b = t/2 + t%2

	}
	Printf("%d %d\n", v+a, v+b)
}

func main() {
	//A(os.Stdin, os.Stdout)
	//B(os.Stdin, os.Stdout)
	Solve(os.Stdin, os.Stdout)
}

func Solve(r io.Reader, w io.Writer) {
	in := bufio.NewReader(r)
	out := bufio.NewWriter(w)
	defer out.Flush()
	var n, m int
	Fscan(in, &n, &m)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		Fscan(in, &arr[i])
	}
	Println(arr)
	diff := make([]int, n+2)
	insert := func(l, r, cnt int) {
		diff[l] += cnt
		diff[r+1] -= cnt
	}
	for i := 1; i <= n; i++ {
		insert(i, i, arr[i])
	}
	for i := 0; i < m; i++ {
		var left, right, cnt int
		Fscan(in, &left, &right, &cnt)
		insert(left, right, cnt)
	}
	for i := 1; i <= n; i++ {
		diff[i] += diff[i-1]
		Fprintf(out, "%v\n", diff[i])
	}
}
