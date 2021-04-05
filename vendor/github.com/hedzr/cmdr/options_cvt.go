// Copyright © 2019 Hedzr Yeh.

package cmdr

import "strconv"

func stringSliceToIntSlice(in []string) (out []int) {
	for _, ii := range in {
		if i, err := strconv.Atoi(ii); err == nil {
			out = append(out, i)
		}
	}
	return
}

func stringSliceToInt64Slice(in []string) (out []int64) {
	for _, ii := range in {
		if i, err := strconv.ParseInt(ii, 0, 64); err == nil {
			out = append(out, i)
		}
	}
	return
}

func stringSliceToUint64Slice(in []string) (out []uint64) {
	for _, ii := range in {
		if i, err := strconv.ParseUint(ii, 0, 64); err == nil {
			out = append(out, i)
		}
	}
	return
}

func intSliceToInt64Slice(in []int) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}

func intSliceToUint64Slice(in []int) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

func int64SliceToIntSlice(in []int64) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func int64SliceToUint64Slice(in []int64) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

func uint64SliceToIntSlice(in []uint64) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func uint64SliceToInt64Slice(in []uint64) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}
