package main

import (
	"bufio"
	. "fmt"
	"io"
	"sort"
	"strconv"
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
	Println(maxSmaller(24333, []int{4, 3}))
}
func maxSmaller(N int, arr []int) string {
	sort.Ints(arr)
	str := strconv.Itoa(N)
	n := len(str)
	ans := make([]byte, n)
	for i := 0; i < n; i++ {
		cur := search(arr, int(str[i]&15))
		ans[i] = byte(cur + '0')
	}
	for i := 0; i < n; i++ {
		if ans[i] < str[i] {
			for j := i + 1; j < n; j++ {
				ans[j] = byte(arr[len(arr)-1] + '0')
			}
			return string(ans)
		}
	}
	for i := 0; i < n; i++ {
		if ans[i] > str[i] {
			for j := i - 1; j >= 0; j-- {
				if int(ans[j]&15) > arr[0] {
					cur := searchLess(arr, int(ans[j]&15))
					ans[j] = byte(cur + '0')
					for k := j + 1; k < n; k++ {
						ans[k] = byte(arr[len(arr)-1] + '0')
					}
					return string(ans)
				}
			}
			var tmp []byte
			for j := 0; j < n-1; j++ {
				tmp = append(tmp, byte(arr[len(arr)-1]+'0'))
			}
			return string(tmp)
		}
	}
	for i := n - 1; i >= 0; i-- {
		if int(ans[i]&15) > arr[0] {
			cur := searchLess(arr, int(ans[i]&15))
			ans[i] = byte(cur + '0')
			for j := i + 1; j < n; j++ {
				ans[j] = byte(arr[len(arr)-1] + '0')
			}
			return string(ans)
		}
	}
	var tmp []byte
	for i := 0; i < n-1; i++ {
		tmp = append(tmp, byte(arr[len(arr)-1]+'0'))
	}
	return string(tmp)
}
func search(arr []int, target int) int {
	ans := arr[0]
	for _, num := range arr {
		if num > target {
			break
		}
		ans = num
	}
	return ans
}
func searchLess(arr []int, target int) int {
	ans := arr[0]
	for _, num := range arr {
		if num >= target {
			break
		}
		ans = num
	}
	return ans
}
