package unsafeless

import (
	"testing"
)

func TestUint32(t *testing.T) {
	x := uint32(0x12345678)
	addr := Transmute[*uint32, uintptr](&x)
	*Transmute[uintptr, *uint16](addr + 2) = 0xdead
	*Transmute[uintptr, *uint16](addr) = 0xbeef
	if x != 0xdeadbeef {
		t.Errorf("x = 0x%x; want 0xdeadbeef", x)
	}
}

func TestString(t *testing.T) {
	// A sufficiently smart compiler will optimize this and put it in .rodata and you'll get a SIGSEGV in a few
	// lines. Luckily for us, Go is not that compiler.
	s := "Hello World--- "
	s += "TRIM ---"
	h := Transmute[*string, *struct {
		data uintptr
		len  int
		cap int
	}](&s)
	if h.len != 23 {
		t.Errorf("h.len = %d; want 23", h.len)
	}
	// There's actually no Cap field in reflect.StringHeader, but hopefully the memory layout fairies will be on
	// our side and leave us space for this field (which we need for transmuting to []byte).
	h.cap = h.len
	h.len = 11
	if s != "Hello World" {
		t.Errorf("s = %q; want \"Hello World\"", s)
	}
	// We use pointers here because by-value will discard our masterfully prepared Cap field.
	b := *Transmute[*string, *[]byte](&s)
	if len(b) != 11 {
		t.Errorf("len(b) = %d; want 11", len(b))
	}
	copy(b[:5], []byte("ABCDE"))
	if s != "ABCDE World" {
		t.Errorf("s = %q; want \"ABCDE World\"", s)
	}
}
