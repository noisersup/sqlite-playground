package db

import "fmt"

type iterator struct {
	docBin []byte
	docLen int
	i      int
}

var ErrIteratorDone = fmt.Errorf("Iterator done")

func NewIterator(doc []byte) *iterator {
	return &iterator{
		docBin: doc,
		docLen: len(doc),
	}
}

func (iter *iterator) Next() (byte, error) {
	if iter.i >= iter.docLen {
		return 0, ErrIteratorDone
	}

	res := iter.docBin[iter.i]
	iter.i++

	return res, nil
}
