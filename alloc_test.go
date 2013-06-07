package goalloc

import "testing"

func TestAlloc(t *testing.T) {
	b, e := Alloc(8)
	if e != nil {
		t.Fatal(e)
	}
	if e = b.Resize(4); e != nil {
		t.Fatal(e)
	}
	if e = b.Resize(16); e != nil {
		t.Fatal(e)
	}
	if b.Size() != 16 {
		t.Fatal("Need size 16 have", b.Size())
	}
}

func TestAllocArray(t *testing.T) {
	b, e := AllocArray(int16(0), 8)
	if e != nil {
		t.Fatal(e)
	}
	arr, e := b.ToInt16Slice()
	if e != nil {
		t.Fatal(e)
	}
	if len(arr) != 8 {
		t.Fatal("Need size 8 have", len(arr))
	}
}

func TestInt(t *testing.T) {
	b, e := Alloc(8)
	if e != nil {
		t.Fatal(e)
	}
	i, e := b.ToInt()
	if e != nil {
		t.Fatal(e)
	}
	if *i != 0 {
		t.Fatal("int is not zero, is", i)
	}
	*i = 200
	i, e = b.ToInt()
	if e != nil {
		t.Fatal(e)
	}
	if *i != 200 {
		t.Fatal("int is not 200, is ", i)
	}
	b.Resize(3)
	_, e = b.ToInt()
	if e == nil {
		t.Fatal("impossible conversion")
	}
}

func TestByteSlice(t *testing.T) {
	b, e := Alloc(8)
	if e != nil {
		t.Fatal(e)
	}
	i, e := b.ToInt32()
	if e != nil {
		t.Fatal(e)
	}
	if *i != 0 {
		t.Fatal("int is not zero, is", i)
	}
	*i = 0x11223344
	arr, e := b.ToByteSlice()
	if e != nil {
		t.Fatal(e)
	}
	if arr[0] != 0x44 {
		t.Fatal("index 0 is not 0x44, is", arr[0])
	}
	if arr[1] != 0x33 {
		t.Fatal("index 1 is not 0x33, is", arr[1])
	}
	if arr[2] != 0x22 {
		t.Fatal("index 2 is not 0x22, is", arr[2])
	}
	if arr[3] != 0x11 {
		t.Fatal("index 3 is not 0x11, is", arr[3])
	}
}

func TestLoadSlice(t *testing.T) {
	b := []byte{0x44, 0x33, 0x22, 0x11}
	m, e := Load(b)
	if e != nil {
		t.Fatal(e)
	}
	i, e := m.ToInt32()
	if e != nil {
		t.Fatal(e)
	}
	if *i != 0x11223344 {
		t.Fatal("int is not 0x11223344, is", *i)
	}
}

func TestLoadArray(t *testing.T) {
	b := &[4]byte{0x44, 0x33, 0x22, 0x11}
	m, e := Load(b)
	if e != nil {
		t.Fatal(e)
	}
	i, e := m.ToInt32()
	if e != nil {
		t.Fatal(e)
	}
	if *i != 0x11223344 {
		t.Fatal("int is not 0x11223344, is", *i)
	}
	*i = 0x44332211
	if b[3] != 0x44 {
		t.Fatal("index 3 is not 0x44, is", b[3])
	}
	if b[2] != 0x33 {
		t.Fatal("index 2 is not 0x33, is", b[2])
	}
	if b[1] != 0x22 {
		t.Fatal("index 1 is not 0x22, is", b[1])
	}
	if b[0] != 0x11 {
		t.Fatal("index 0 is not 0x11, is", b[0])
	}
}

func TestLoadPtr(t *testing.T) {
	b := []byte{0x44, 0x33, 0x22, 0x11}
	m, e := Load(&b)
	if e != nil {
		t.Fatal(e)
	}
	i, e := m.ToInt32()
	if e != nil {
		t.Fatal(e)
	}
	if *i != 0x11223344 {
		t.Fatal("int is not 0x11223344, is", *i)
	}
}

func TestLoadStruct(t *testing.T) {
	b := struct {
		a int
		b int
	}{1, 2}
	_, e := Load(&b)
	if e != nil {
		t.Fatal(e)
	}
	_, e = Load(b)
	if e == nil {
		t.Fatal("should not give error, got", e)
	}
}

func TestLoadString(t *testing.T) {
	s := "1234"
	m, e := Load(&s)
	if e != nil {
		t.Fatal(e)
	}

	arr, e := m.ToByteSlice()
	if e != nil {
		t.Fatal(e)
	}
	if arr[0] != 0x31 {
		t.Fatal("index 0 is not 0x31, is", arr[0])
	}
	if arr[1] != 0x32 {
		t.Fatal("index 1 is not 0x32, is", arr[1])
	}
	if arr[2] != 0x33 {
		t.Fatal("index 2 is not 0x33, is", arr[2])
	}
	if arr[3] != 0x34 {
		t.Fatal("index 3 is not 0x34, is", arr[3])
	}

	m, e = Load(s)
	if e != nil {
		t.Fatal(e)
	}
	arr, e = m.ToByteSlice()
	if e != nil {
		t.Fatal(e)
	}
	if arr[0] != 0x31 {
		t.Fatal("index 0 is not 0x31, is", arr[0])
	}
	if arr[1] != 0x32 {
		t.Fatal("index 1 is not 0x32, is", arr[1])
	}
	if arr[2] != 0x33 {
		t.Fatal("index 2 is not 0x33, is", arr[2])
	}
	if arr[3] != 0x34 {
		t.Fatal("index 3 is not 0x34, is", arr[3])
	}
}

func TestToInterface(t *testing.T) {
	b := 0x11223344
	m, e := Load(&b)
	if e != nil {
		t.Fatal(e)
	}
	var arr []uint16
	arrI, e := m.ToInterface(arr)
	arr = arrI.([]uint16)
	if e != nil {
		t.Fatal(e)
	}
	if arr[0] != 0x3344 {
		t.Fatal("index 0 is not 0x3344, is", arr[0])
	}
	if arr[1] != 0x1122 {
		t.Fatal("index 1 is not 0x1122, is", arr[1])
	}
}

func TestFromToInterface(t *testing.T) {
	b := 0x11223344
	m, e := Load(interface{}(&b))
	if e != nil {
		t.Fatal(e)
	}
	var arr []uint16
	arrI, e := m.ToInterface(arr)
	arr = arrI.([]uint16)
	if e != nil {
		t.Fatal(e)
	}
	if arr[0] != 0x3344 {
		t.Fatal("index 0 is not 0x3344, is", arr[0])
	}
	if arr[1] != 0x1122 {
		t.Fatal("index 1 is not 0x1122, is", arr[1])
	}
}

func TestFromToInterface2(t *testing.T) {
	b := 0x11223344
	m, e := Load(interface{}(&b))
	if e != nil {
		t.Fatal(e)
	}
	type Foo struct {
		A, B uint16
	}
	arrI, e := m.ToInterface(interface{}(Foo{}))
	arr := arrI.(*Foo)
	if e != nil {
		t.Fatal(e)
	}
	if arr.A != 0x3344 {
		t.Fatal("A is not 0x3344, is", arr.A)
	}
	if arr.B != 0x1122 {
		t.Fatal("B is not 0x1122, is", arr.B)
	}
}

func TestFromToInterface3(t *testing.T) {
	b := 0x11223344
	m, e := Load(interface{}(&b))
	if e != nil {
		t.Fatal(e)
	}
	nI, e := m.ToInterface(int(0))
	n := nI.(*int)
	if e != nil {
		t.Fatal(e)
	}
	if *n != 0x11223344 {
		t.Fatal("n is not 0x11223344, is", *n)
	}
}

func TestFromToInterface4(t *testing.T) {
	b := 0x11223344
	m, e := Load(interface{}(&b))
	if e != nil {
		t.Fatal(e)
	}
	arrI, e := m.ToInterface([2]uint16{})
	arr := arrI.(*[2]uint16)
	if e != nil {
		t.Fatal(e)
	}
	if arr[0] != 0x3344 {
		t.Fatal("index 0 is not 0x3344, is", arr[0])
	}
	if arr[1] != 0x1122 {
		t.Fatal("index 1 is not 0x1122, is", arr[1])
	}
}
