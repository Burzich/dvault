package kv

import "errors"

var ErrPathNotFound = errors.New("path not found")
var ErrVersionNotFound = errors.New("version not found")
var ErrCas = errors.New("cas not mapped")
