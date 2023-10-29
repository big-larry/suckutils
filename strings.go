package suckutils

import (
	"math"
	"strings"
	"unsafe"
)

func WhereString(str []string, condition func(s string) bool) []string {
	oldlen := len(str)
	mask := make([]byte, oldlen/8+1)
	newlen := 0
	i := 0
	for i = 0; i < oldlen; i++ {
		if condition(str[i]) {
			mask[i/8] |= 1 << (i % 8)
			newlen++
		}
	}
	result := make([]string, newlen)
	c := 0
	n := 0
	for n = 0; n < oldlen/8+1; n++ {
		for i = 0; i < 8; i++ {
			// fmt.Print(i, n, ":", 1<<i, ",", b&(1<<i), " ")
			if mask[n]&(1<<i) != 0 {
				result[c] = str[i+n*8]
				c++
			}
		}
	}
	return result
}

func ReduceString(str []string, reduce func(s string) string) []string {
	for i := 0; i < len(str); i++ {
		str[i] = reduce(str[i])
	}
	return str
}

func Contains(str []string, sub string) bool {
	for _, s := range str {
		if sub == s {
			return true
		}
	}
	return false
}

func ContainsAny(str []string, sub []string) bool {
	for _, s := range str {
		for _, ss := range sub {
			if strings.Contains(s, ss) {
				return true
			}
		}
	}
	return false
}
func ContainsIndex(str []string, sub []string) int {
	for i, s := range str {
		for _, ss := range sub {
			if strings.Contains(s, ss) {
				return i
			}
		}
	}
	return -1
}

func ConcatNonString(elems ...interface{}) string {
	n := 0
	for i := 0; i < len(elems); i++ {
		n += len(elems[i].(string))
	}
	result := make([]byte, 0, n)
	for i := 0; i < len(elems); i++ {
		result = append(result, elems[i].(string)...)
	}
	return *(*string)(unsafe.Pointer(&result))
}

func Concat(elems ...string) string {
	// switch k {
	// case 0:
	// 	return ""
	// case 1:
	// 	return elems[0]
	// }
	n := 0
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}
	// result := make([]byte, n, n)
	result := make([]byte, 0, n)
	//var j int
	for i := 0; i < len(elems); i++ {
		result = append(result, elems[i]...)
		//j += copy(result[j:], elems[i])
	}
	return *(*string)(unsafe.Pointer(&result))
}
func ConcatTwo(s1, s2 string) string {
	result := make([]byte, len(s1)+len(s2))
	copy(result[copy(result, s1):], s2)
	return *(*string)(unsafe.Pointer(&result))
}
func ConcatThree(s1, s2, s3 string) string {
	result := make([]byte, len(s1)+len(s2)+len(s3))
	copy(result[len(s1)+copy(result[copy(result, s1):], s2):], s3)
	return *(*string)(unsafe.Pointer(&result))
}
func ConcatFour(s1, s2, s3, s4 string) string {
	result := make([]byte, len(s1)+len(s2)+len(s3)+len(s4))
	// result = append(result, elems[i]...)
	copy(result[len(s1)+len(s2)+copy(result[len(s1)+copy(result[copy(result, s1):], s2):], s3):], s4)
	return *(*string)(unsafe.Pointer(&result))
}

const digits = "0123456789"

func Itoa(num uint32) string {
	result := make([]byte, int(math.Log10(float64(num))+1))
	for i := len(result) - 1; i >= 0; i-- {
		result[i] = digits[num%10]
		num /= 10
	}

	return string(result)
}
