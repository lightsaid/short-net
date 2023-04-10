package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	var nums = []uint{1, 2, 9, 10, 99, 10000, 8765432, 11232390987, 1212128999, 66554232, 90832, 22, 332323, 44, 55998833434, 88}

	for _, num := range nums {
		h := EncodeBase62(num)
		fmt.Println(num, " = ", h)
		n := DecodeBase62(h)

		require.Equal(t, num, n)
	}

	for i := 0; i < 62; i++ {
		num := uint(i)
		h := EncodeBase62(num)
		fmt.Println(num, " = ", h)
		n := DecodeBase62(h)

		require.Equal(t, num, n)
	}

	for i := 3000; i < 99999; i++ {
		num := uint(i)
		h := EncodeBase62(num)
		// fmt.Println(num, " = ", h)
		n := DecodeBase62(h)

		require.Equal(t, num, n)
	}

	for i := 1000000; i < 9999999; i++ {
		num := uint(i)
		h := EncodeBase62(num)
		// fmt.Println(num, " = ", h)
		n := DecodeBase62(h)

		require.Equal(t, num, n)
	}
}
