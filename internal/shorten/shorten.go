package shorten

import "strings"

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = uint32(len(alphabet))

func Shorten(id uint32) string {
	var (
		nums    []uint32
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		nums = append(nums, num%alphabetLen)
		num = num / alphabetLen
	}

	reverse(nums)

	for _, num := range nums {
		builder.WriteString(string(alphabet[num]))
	}

	return builder.String()
}

func reverse(s []uint32) []uint32 {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
