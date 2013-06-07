package goalloc

import (
	"errors"
	"reflect"
	"unsafe"
)

var AlreadyFreed = errors.New("Already freed.")
var SmallSize = errors.New("Size is too small.")
var LoadError = errors.New("Cannot load value.")
var GCManaged = errors.New("This memory is gc protected.")
var NotAPointer = errors.New("Not a pointer.")
var UnsupportedValue = errors.New("Unsupported value.")

var intSize = int(reflect.TypeOf(int(0)).Size())

type MemBlock struct {
	size      int
	freed     bool
	gcmanaged bool
	ptr       unsafe.Pointer
}

func Alloc(size int) (*MemBlock, error) {
	ptr, err := alloc(size)
	if err != nil {
		return nil, err
	}
	return &MemBlock{size, false, false, ptr}, nil
}

func AllocArray(typ interface{}, size int) (*MemBlock, error) {
	finalSize := size * int(reflect.TypeOf(typ).Size())
	ptr, err := alloc(finalSize)
	if err != nil {
		return nil, err
	}
	return &MemBlock{finalSize, false, false, ptr}, nil
}

func Load(data interface{}) (*MemBlock, error) {

	val := reflect.ValueOf(data)
	prev := reflect.ValueOf(nil)

Switch:
	switch val.Kind() {
	case reflect.Slice:
		return &MemBlock{val.Cap(), false, true, unsafe.Pointer(val.Pointer())}, nil
	case reflect.Array:
		if prev.Kind() != reflect.Ptr {
			return nil, NotAPointer
		}
		return &MemBlock{val.Cap(), false, true, unsafe.Pointer(prev.Pointer())}, nil
	case reflect.Ptr:
		prev = val
		val = val.Elem()
		goto Switch
	case reflect.String:
		var header *reflect.StringHeader
		if prev.Kind() != reflect.Ptr {
			s := val.String()
			header = (*reflect.StringHeader)(unsafe.Pointer(&s))
		} else {
			header = (*reflect.StringHeader)(unsafe.Pointer(prev.Pointer()))
		}
		return &MemBlock{header.Len, false, true, unsafe.Pointer(header.Data)}, nil
	case reflect.Chan, reflect.Map, reflect.Func, reflect.Interface:
		return nil, UnsupportedValue
	default:
		if prev.Kind() != reflect.Ptr {
			return nil, NotAPointer
		}
		return &MemBlock{int(val.Type().Size()), false, true, unsafe.Pointer(prev.Pointer())}, nil
	}
}

func (this *MemBlock) Size() int {
	return this.size
}

func (this *MemBlock) IsGCManaged() bool {
	return this.gcmanaged
}

func (this *MemBlock) WasFreed() bool {
	return this.freed
}

func (this *MemBlock) Resize(size int) error {
	if this.freed {
		return AlreadyFreed
	}
	if this.gcmanaged {
		return GCManaged
	}
	ptr, err := realloc(this.ptr, size)
	if err != nil {
		return err
	}
	this.ptr = ptr
	this.size = size
	return nil
}

func (this *MemBlock) Free() {
	if this.freed || this.gcmanaged {
		return
	}
	free(this.ptr)
	this.freed = true
}

func (this *MemBlock) ToInterface(typ interface{}) (interface{}, error) {
	if this.freed {
		return nil, AlreadyFreed
	}

	v := reflect.TypeOf(typ)
Switch:
	switch v.Kind() {
	case reflect.Slice:
		sSize := this.size / int(v.Elem().Size())
		if sSize == 0 {
			return nil, SmallSize
		}
		nv := reflect.New(v)
		header := (*reflect.SliceHeader)(unsafe.Pointer(nv.Pointer()))
		header.Cap = sSize
		header.Len = sSize
		header.Data = uintptr(this.ptr)
		return nv.Elem().Interface(), nil
	case reflect.Ptr:
		v = v.Elem()
		goto Switch
	case reflect.Chan, reflect.Map, reflect.Func, reflect.Interface:
		return nil, UnsupportedValue
	default:
		if this.size < int(v.Size()) {
			return nil, SmallSize
		}
		return reflect.NewAt(v, this.ptr).Interface(), nil
	}
}

func (this *MemBlock) ToInt() (*int, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < intSize {
		return nil, SmallSize
	}
	return (*int)(this.ptr), nil
}

func (this *MemBlock) ToInt8() (*int8, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 1 {
		return nil, SmallSize
	}
	return (*int8)(this.ptr), nil
}

func (this *MemBlock) ToInt16() (*int16, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 2 {
		return nil, SmallSize
	}
	return (*int16)(this.ptr), nil
}

func (this *MemBlock) ToInt32() (*int32, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 4 {
		return nil, SmallSize
	}
	return (*int32)(this.ptr), nil
}

func (this *MemBlock) ToInt64() (*int64, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 8 {
		return nil, SmallSize
	}
	return (*int64)(this.ptr), nil
}

func (this *MemBlock) ToUInt() (*uint, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < intSize {
		return nil, SmallSize
	}
	return (*uint)(this.ptr), nil
}

func (this *MemBlock) ToUInt8() (*uint8, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 1 {
		return nil, SmallSize
	}
	return (*uint8)(this.ptr), nil
}

func (this *MemBlock) ToByte() (*byte, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 1 {
		return nil, SmallSize
	}
	return (*byte)(this.ptr), nil
}

func (this *MemBlock) ToUInt16() (*uint16, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 2 {
		return nil, SmallSize
	}
	return (*uint16)(this.ptr), nil
}

func (this *MemBlock) ToUInt32() (*uint32, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 4 {
		return nil, SmallSize
	}
	return (*uint32)(this.ptr), nil
}

func (this *MemBlock) ToUInt64() (*uint64, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 8 {
		return nil, SmallSize
	}
	return (*uint64)(this.ptr), nil
}

func (this *MemBlock) ToIntSlice() ([]int, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < intSize {
		return nil, SmallSize
	}
	sSize := this.size / intSize

	var slice []int
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToByteSlice() ([]byte, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 1 {
		return nil, SmallSize
	}
	sSize := this.size

	var slice []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToInt8Slice() ([]int8, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 1 {
		return nil, SmallSize
	}
	sSize := this.size

	var slice []int8
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToInt16Slice() ([]int16, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 2 {
		return nil, SmallSize
	}
	sSize := this.size / 2

	var slice []int16
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToInt32Slice() ([]int32, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 4 {
		return nil, SmallSize
	}
	sSize := this.size / 4

	var slice []int32
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToInt64Slice() ([]int64, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 8 {
		return nil, SmallSize
	}
	sSize := this.size / 8

	var slice []int64
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToUIntSlice() ([]uint, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < intSize {
		return nil, SmallSize
	}
	sSize := this.size / intSize

	var slice []uint
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToUInt8Slice() ([]uint8, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 1 {
		return nil, SmallSize
	}
	sSize := this.size

	var slice []uint8
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToUInt16Slice() ([]uint16, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 2 {
		return nil, SmallSize
	}
	sSize := this.size / 2

	var slice []uint16
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToUInt32Slice() ([]uint32, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 4 {
		return nil, SmallSize
	}
	sSize := this.size / 4

	var slice []uint32
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}

func (this *MemBlock) ToUInt64Slice() ([]uint64, error) {
	if this.freed {
		return nil, AlreadyFreed
	}
	if this.size < 8 {
		return nil, SmallSize
	}
	sSize := this.size / 8

	var slice []uint64
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = sSize
	header.Len = sSize
	header.Data = uintptr(this.ptr)

	return slice, nil
}
