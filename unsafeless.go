package unsafeless

// The runtime representation for a interface{} value, based on iface from src/reflect/value.go
type iface struct {
	typ  uintptr
	word uintptr
}

type transmuter func(interface{}, interface{}) interface{}

type racer interface {
	Race() transmuter
}

type s1 struct {
	f func(iface, iface) iface
}

func (s *s1) Race() transmuter {
	return nil
}

type s2 struct {
	f transmuter
}

func (s *s2) Race() transmuter {
	return s.f
}

// Being able to change the type of a value is the ideal primitive. You may not like it, but this is what peak unsafety
// looks like.
var transmute transmuter

func init() {
	a := &s1{f: func(t iface, u iface) iface {
		return iface{typ: u.typ, word: t.word}
	}}
	b := &s2{}
	r := racer(a)

	go func() {
		// Assigning to r requires storing two words (typ and word). Storing multiple words isn't atomic so
		// r.Race can be called after typ is written, but before word is written, and the method is called with
		// the wrong receiver. If b.Race is called with a as the receiver, then it'll return a.f, but casted
		// to transmuter. This gives us a function that takes interface{} and can modify its typ and word
		// (without using "unsafe" or "reflect" at all).
		for transmute == nil {
			r = b
			if r == nil {
				// This should never be reachable but it prevents the compiler optimizing out the first
				// assignment (which is critical for the data race).
				panic(r)
			}
			r = a
		}
	}()

	for transmute == nil {
		// When the data race causes b.Race to be called with a as the receiver, this will return a.f (but
		// casted to transmuter). Otherwise it returns nil and we try again.
		transmute = r.Race()
	}
}

func Transmute[T, U any](t T) U {
	return *transmute(&t, new(U)).(*U)
}
