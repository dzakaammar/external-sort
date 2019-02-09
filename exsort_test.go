package exsort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSorting(t *testing.T) {
	i, err := OpenFile(source)
	assert.NoError(t, err)

	o, err := OpenFile(output)
	assert.NoError(t, err)

	err = externalSort(i, o, size)
	assert.NoError(t, err)
}
