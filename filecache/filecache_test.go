package filecache

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileCache(t *testing.T) {
	filename := getTempFileName(t)
	const exp = 100
	var res int
	var nCalled int
	for i := 0; i < 10; i++ {
		err := Load(&res, filename, func() {
			nCalled++
			res = exp
		})
		if err != nil {
			t.Fatalf("Error occurred while loading cache: %v", err)
		}
		if res != exp {
			t.Errorf("Expected %d, got %d", exp, res)
		}
	}
	if nCalled != 1 {
		t.Errorf("Epected a function to be called once, was called %d times", nCalled)
	}
}

func getTempFileName(t *testing.T) string {
	t.Helper()
	fd, err := ioutil.TempFile("", "filecache")
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}
	if err := fd.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}
	if err := os.Remove(fd.Name()); err != nil {
		t.Fatalf("Failed to delete temporary file: %v", err)
	}
	return fd.Name()
}
