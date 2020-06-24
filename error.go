package hdbscan

import "errors"

var (
	// ErrMCS ...
	ErrMCS = errors.New("minimum cluster size is too small")
	// ErrDataLen ...
	ErrDataLen = errors.New("length of data is less than minimum cluster size")
	// ErrRowLength ...
	ErrRowLength = errors.New("row is incorrect length")
)
