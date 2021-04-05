// Copyright © 2020 Hedzr Yeh.

package tool

import (
	"sort"
	"strings"
)

type byDottedSlice []string

func (s byDottedSlice) Len() int      { return len(s) }
func (s byDottedSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byDottedSlice) Less(i, j int) (ret bool) {
	sa := strings.Split(s[i], ".")
	sb := strings.Split(s[j], ".")
	la, lb := len(sa), len(sb)

	ll := la
	if ll > lb {
		ll = lb
	}

	var lastRet int
	for i := 0; i < ll; i++ {
		switch rb := strings.Compare(sa[i], sb[i]); {
		case rb == -1:
			if lastRet != rb {
				return true
			}
			lastRet = rb
			continue
		case rb == 0:
			if lastRet != rb {
				return lastRet < 0
			}
			lastRet = rb
			continue
		default:
			ret = false
			return
		}
	}

	rr := strings.Compare(sa[ll-1], sb[ll-1])
	if rr < 0 {
		ret = true
		return
	}

	if rr == 0 {
		ret = la <= lb
		return
	}

	return
}

type byDottedSliceRev []string

func (s byDottedSliceRev) Len() int      { return len(s) }
func (s byDottedSliceRev) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byDottedSliceRev) Less(i, j int) (ret bool) {
	sa := strings.Split(s[i], ".")
	sb := strings.Split(s[j], ".")
	la, lb := len(sa), len(sb)

	ll := la
	if ll > lb {
		ll = lb
	}

	var lastRet int
	for i := 0; i < ll; i++ {
		switch rb := strings.Compare(sa[i], sb[i]); {
		case rb == -1:
			if lastRet != rb {
				return false
			}
			lastRet = rb
			continue
		case rb == 0:
			if lastRet != rb {
				return lastRet > 0
			}
			lastRet = rb
			continue
		default:
			ret = true
			return
		}
	}

	rr := strings.Compare(sa[ll-1], sb[ll-1])
	if rr < 0 {
		ret = false
		return
	}

	if rr == 0 {
		ret = la > lb
		return
	}

	return
}

// SortAsDottedSlice sorts a slice by which is treated as a dot-separated path string
func SortAsDottedSlice(ks []string) {
	sort.Sort(byDottedSlice(ks))
}

// SortAsDottedSliceReverse sorts a slice by which is treated as a dot-separated path string
func SortAsDottedSliceReverse(ks []string) {
	sort.Sort(byDottedSliceRev(ks))
}
