package utils

// PermuteStrings generates all permutations of strings provided
func PermuteStrings(parts ...[]string) (ret [][]string) {
	{
		var n = 1
		for _, ar := range parts {
			n *= len(ar)
		}
		ret = make([][]string, 0, n)
	}
	var at = make([]int, len(parts))
	var buf2 []string
loop:
	for {
		// increment position counters
		for i := len(parts) - 1; i >= 0; i-- {
			if at[i] > 0 && at[i] >= len(parts[i]) {
				if i == 0 || (i == 1 && at[i-1] == len(parts[0])-1) {
					break loop
				}
				at[i] = 0
				at[i-1]++
			}
		}
		// construct permutated string
		buf2 = []string{}
		for i, ar := range parts {
			var p = at[i]
			if p >= 0 && p < len(ar) {
				buf2 = append(buf2, ar[p])
			}
		}

		ret = append(ret, buf2)
		at[len(parts)-1]++
	}
	return ret
}
