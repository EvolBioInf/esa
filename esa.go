package esa

/*
#cgo LDFLAGS: -ldivsufsort64 -L/opt/homebrew/lib
#cgo CFLAGS: -I/opt/homebrew/include
#include <divsufsort64.h>
#include <stdlib.h>
*/
import "C"
import (
	"log"
	"reflect"
	"unsafe"
)

type Stack []*Interval
type Interval struct {
	Idx int
	Lcp int
}
type Esa struct {
	T   []byte
	Sa  []int
	Lcp []int
	Cld []int
}
type Minterval struct {
	I, J int
	L    int
}

func (s *Stack) Top() *Interval {
	return (*s)[len(*s)-1]
}
func (s *Stack) Pop() *Interval {
	i := (*s)[len(*s)-1]
	(*s) = (*s)[0 : len(*s)-1]
	return i
}
func (s *Stack) Push(i *Interval) {
	(*s) = append(*s, i)
}
func (e *Esa) MatchPref(p []byte) *Minterval {
	k := 0
	m := len(p)
	var parent, child *Minterval
	parent = new(Minterval)
	parent.I = 0
	parent.J = len(e.T) - 1
	for k < m {
		child = e.GetInterval(parent, p[k])
		if child == nil {
			parent.L = k
			return parent
		}
		l := m
		i := child.I
		j := child.J
		if i < j {
			r := 0
			if e.Lcp[i] <= e.Lcp[j+1] {
				r = e.Cld[j]
			} else {
				r = e.Cld[i]
			}
			l = min(l, e.Lcp[r])
		}
		for w := k + 1; w < l; w++ {
			if e.T[e.Sa[i]+w] != p[w] {
				child.L = w
				return child
			}
		}
		k = l
	}
	child.L = k
	return child
}
func Sa(t []byte) []int {
	var sa []int
	header := (*reflect.SliceHeader)(unsafe.Pointer(&t))
	ct := (*C.sauchar_t)(unsafe.Pointer(header.Data))
	n := len(t)
	csa := (*C.saidx64_t)(C.malloc(C.size_t(n * C.sizeof_saidx64_t)))
	cn := C.saidx64_t(n)
	err := int(C.divsufsort64(ct, csa, cn))
	if err != 0 {
		log.Fatalf("divsufsort failed with code %d\n", err)
	}
	header = (*reflect.SliceHeader)((unsafe.Pointer(&sa)))
	header.Cap = n
	header.Len = n
	header.Data = uintptr(unsafe.Pointer(csa))
	return sa
}
func Lcp(t []byte, sa []int) []int {
	n := len(t)
	lcp := make([]int, n)
	isa := make([]int, n)
	for i := 0; i < n; i++ {
		isa[sa[i]] = i
	}
	lcp[0] = -1
	l := 0
	for i := 0; i < n; i++ {
		j := isa[i]
		if j == 0 {
			continue
		}
		k := sa[j-1]
		for k+l < n && i+l < n && t[k+l] == t[i+l] {
			l++
		}
		lcp[j] = l
		l -= 1
		if l < 0 {
			l = 0
		}
	}
	return lcp
}
func Cld(lcp []int) []int {
	var cld []int
	lcp = append(lcp, -1)
	n := len(lcp) - 1
	cld = make([]int, n+1)
	cld[0] = n
	stack := new(Stack)
	iv := newInterval(0, -1)
	stack.Push(iv)
	for i := 1; i <= n; i++ {
		top := stack.Top()
		for lcp[i] < top.Lcp {
			last := stack.Pop()
			top = stack.Top()
			for top.Lcp == last.Lcp {
				cld[top.Idx] = last.Idx
				last = stack.Pop()
				top = stack.Top()
			}
			top = stack.Top()
			if lcp[i] < top.Lcp {
				cld[top.Idx] = last.Idx
			} else {
				cld[i-1] = last.Idx
			}
		}
		iv = newInterval(i, lcp[i])
		stack.Push(iv)
	}
	lcp = lcp[:len(lcp)-1]
	return cld
}
func newInterval(i, l int) *Interval {
	iv := new(Interval)
	iv.Idx = i
	iv.Lcp = l
	return iv
}
func MakeEsa(t []byte) *Esa {
	esa := new(Esa)
	esa.T = t
	esa.T = append(esa.T, 0)
	esa.Sa = Sa(esa.T)
	esa.Lcp = Lcp(esa.T, esa.Sa)
	esa.Lcp = append(esa.Lcp, -1)
	esa.Cld = Cld(esa.Lcp)
	return esa
}
func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
func (e *Esa) GetInterval(iv *Minterval, c byte) *Minterval {
	i := iv.I
	j := iv.J
	if i == j {
		if e.T[e.Sa[i]] == c {
			return iv
		}
	}
	m := 0
	if e.Lcp[i] <= e.Lcp[j+1] {
		m = e.Cld[j]
	} else {
		m = e.Cld[i]
	}
	l := e.Lcp[m]
	k := i
	for e.Lcp[m] == l {
		if e.T[e.Sa[k]+l] == c {
			iv.I = k
			iv.J = m - 1
			return iv
		}
		k = m
		if k == j {
			break
		}
		m = e.Cld[m]
	}
	if e.T[e.Sa[k]+l] == c {
		iv.I = k
		iv.J = j
		return iv
	}
	return nil
}
