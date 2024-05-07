package core

import "errors"

var (
	ErrOutputIsNil   = errors.New("output is nil")
	ErrAlreadyExists = errors.New("already exists")
	ErrLinkIsNotTidy = errors.New("link is not tidy")
)
