package models

import (
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type ErrorWithPath struct {
	error `json:"error"`
	Path  string `json:"path"`
	sep   string
}

func (e ErrorWithPath) Error() string {
	return fmt.Errorf("%s: %v", e.Path, e.error).Error()
}

func (e ErrorWithPath) Unwrap() error {
	return e.error
}

func (e ErrorWithPath) PrependPath(path string) ErrorWithPath {
	return ErrorWithPath{
		error: e.error,
		Path:  fmt.Sprintf("%s%s%s", path, e.sep, e.Path),
		sep:   ".",
	}
}

func (e ErrorWithPath) PrependIndex(index int) ErrorWithPath {
	return ErrorWithPath{
		error: e.error,
		Path:  fmt.Sprintf("[%d]%s%s", index, e.sep, e.Path),
		sep:   "",
	}
}

func (e ErrorWithPath) AsErrorDetail() huma.ErrorDetail {
	return huma.ErrorDetail{
		Message:  e.Error(),
		Location: e.Path,
	}
}

func WrapErrorPath(err error, path string) ErrorWithPath {
	if errors.As(err, &ErrorWithPath{}) {
		return err.(ErrorWithPath).PrependPath(path)
	}

	return ErrorWithPath{
		error: err,
		Path:  path,
		sep:   ".",
	}
}

func WrapErrorIndex(err error, index int) ErrorWithPath {
	if errors.As(err, &ErrorWithPath{}) {
		return err.(ErrorWithPath).PrependIndex(index)
	}

	return ErrorWithPath{
		error: err,
		Path:  fmt.Sprintf("[%d]", index),
		sep:   "",
	}
}
