package exsort

//Sorter :nodoc:
type Sorter interface {
	Sort([]string) (interface{}, error)
}

type defaultSorter struct{}

//Sort :nodoc:
func (d *defaultSorter) Sort(_ []string) (interface{}, error) {
	return nil, nil
}
