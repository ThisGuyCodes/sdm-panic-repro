package biglock

import "github.com/awnumar/memcall"

func init() {
	mem, err := memcall.Alloc(4096 * 16377)
	if err != nil {
		panic(err)
	}

	err = memcall.Lock(mem)
	if err != nil {
		panic(err)
	}
}
