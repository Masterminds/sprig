package sprig

// maxAllocSize is the maximum number of elements or bytes that template
// functions are allowed to allocate in a single call.  This prevents
// denial-of-service when user-controlled input is passed as a size parameter.
const maxAllocSize = 1 << 20 // 1 MiB / ~1 million elements

func clampInt(n, max int) int {
	if n < 0 {
		n = 0
	}
	if n > max {
		return max
	}
	return n
}
