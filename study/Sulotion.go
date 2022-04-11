package main

import (
	"bufio"
	. "fmt"
	"io"
	"os"
	"sort"
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
	A(os.Stdin, os.Stdout)
	B(os.Stdin, os.Stdout)
}

func canPartitionKSubsets(nums []int, k int) bool {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	if sum%k != 0 {
		return false
	}
	sum /= k
	sort.Ints(nums)
	n := len(nums)
	if nums[n-1] > sum || sum < nums[0] {
		return false
	}
	bucket := make([]int, k)
	for i := 0; i < k; i++ {
		bucket[i] = sum
	}
	var dfs func(index int) bool
	dfs = func(index int) bool {
		if index < 0 {
			return true
		}
		for i := 0; i < k; i++ {
			if bucket[i] == nums[index] || bucket[i]-nums[index] >= nums[0] {
				bucket[i] -= nums[index]
				if dfs(index - 1) {
					return true
				}
				bucket[i] += nums[index]
			}
		}
		return false
	}
	return dfs(n - 1)
}

func uniqueMorseRepresentations(words []string) int {
	alpha := []string{".-", "-...", "-.-.", "-..", ".", "..-.", "--.", "....", "..", ".---", "-.-", ".-..", "--", "-.", "---", ".--.", "--.-", ".-.", "...", "-", "..-", "...-", ".--", "-..-", "-.--", "--.."}
	encode := func(word string) (ans string) {
		for _, ch := range word {
			ans += alpha[ch-'a']
		}
		return
	}
	set := make(map[string]bool)
	for _, word := range words {
		code := encode(word)
		set[code] = true
	}
	return len(set)
}
