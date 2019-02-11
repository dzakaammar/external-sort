package exsort

import (
	"path"
	"runtime"
	"testing"
)

func TestSorting(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)

	s := New(dir+"/file.txt", dir+"/out,txt", dir+"/out", 10000)
	err := s.Sort()
	if err != nil {
		t.Fatalf("Error : %+v", err)
	}
}
