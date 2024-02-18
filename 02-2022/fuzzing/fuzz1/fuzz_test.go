package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")
}
func FuzzHex(f *testing.F) {
	for _, seed := range [][]byte{{}, {0}, {9}, {0xa}, {0xf}, {1, 2, 3, 4}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		enc := hex.EncodeToString(in)
		out, err := hex.DecodeString(enc)
		if err != nil {
			t.Fatalf("%v: decode: %v", in, err)
		}
		if !bytes.Equal(in, out) {
			t.Fatalf("%v: not equal after round trip: %v", in, out)
		}
	})
}

//A fuzz test maintains a seed corpus, or a set of inputs which are run by default, and can seed input generation. Seed inputs may be registered by calling (*F).Add or by storing files in the directory testdata/fuzz/<Name> (where <Name> is the name of the fuzz test) within the package containing the fuzz test. Seed inputs are optional, but the fuzzing engine may find bugs more efficiently when provided with a set of small seed inputs with good code coverage. These seed inputs can also serve as regression tests for bugs identified through fuzzing.

//The function passed to (*F).Fuzz within the fuzz test is considered the fuzz target. A fuzz target must accept a *T parameter, followed by one or more parameters for random inputs. The types of arguments passed to (*F).Add must be identical to the types of these parameters. The fuzz target may signal that it's found a problem the same way tests do: by calling T.Fail (or any method that calls it like T.Error or T.Fatal) or by panicking.
