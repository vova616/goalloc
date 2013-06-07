package goalloc

//#include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"
)

var AllocationFailed = errors.New("The allocation has failed.")

func alloc(size int) (unsafe.Pointer, error) {
	p := C.calloc(C.size_t(size), C.size_t(1))
	if p == nil {
		return nil, AllocationFailed
	}
	return p, nil
}

func free(ptr unsafe.Pointer) {
	C.free(ptr)
}

func realloc(ptr unsafe.Pointer, size int) (unsafe.Pointer, error) {
	p := C.realloc(ptr, C.size_t(size))
	if p == nil {
		return nil, AllocationFailed
	}
	return p, nil
}
