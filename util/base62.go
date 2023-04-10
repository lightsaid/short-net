package util

import (
	"math"
	"strings"
)

// base62 编码和解码实现原理

/**

十进制转二进制，如 100 转二进制, 62进制同理

100 / 2 = 50 余 0
50  / 2 = 25 余 0
25  / 2 = 12 余 1
12  / 2 = 6  余 0
6   / 2 = 3  余 0
3   / 2 = 1  余 1
1   / 2 = 0  余 1

最终得二进制:  1100100

======================

二进制转十进制，1100100 转十进制，62进制同理

1 * 2^6 = 64
1 * 2^5 = 32
0 * 2^4 = 0
0 * 2^3 = 0
1 * 2^2 = 4
0 * 2^1 = 0
0 * 2^0 = 0

64 + 32 + 4 = 100

多少进制底数就是多少，指数是长度-1

*/

var base62Characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// EncodeBase62 编码
func EncodeBase62(shortID uint) string {
	var size = uint(len(base62Characters))
	var hash = []byte{}

	if shortID == 0 {
		return "0"
	}

	for shortID > 0 {
		rem := shortID % size // 取余数
		// NOTE: 这里有个坑，这样组装的 hash 是反转的，解码会对不上, 如：10000 = iB2，应该是 2Bi 才到
		// hash = append(hash, base62Characters[rem])
		hash = append([]byte{base62Characters[rem]}, hash...)
		shortID = shortID / size // 取商
	}

	return string(hash)
}

// DecodeBase62 解码
func DecodeBase62(hash string) uint {
	var shortID uint
	l := len(hash)

	for i := 0; i < l; i++ {
		index := strings.Index(base62Characters, string(hash[i]))
		shortID += uint(index) * uint(math.Pow(62, float64(l-i-1)))
	}

	return shortID
}
