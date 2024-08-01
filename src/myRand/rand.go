package myRand

const (
	a = 214013
	c = 2531011
)

var x int = 0

func Seed(s int64) {
	x = int(s)
}

func Int() int {
	x = a*x + c
	return x
}

func Intn(n int) int {
	tmp := Int() % n
	if tmp < 0 {
		tmp = tmp + n
	}
	return tmp
}
