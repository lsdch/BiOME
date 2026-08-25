package app_errors

import (
	"database/sql"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type ErrorCode string

//generate:enum
const (
	ErrorCodeInconsistentTaxa ErrorCode = "inconsistent_taxa"
	ErrorCodeInvalidCSVRow    ErrorCode = "invalid_csv_format"
)

type ErrorCategory string

//generate:enum
const (
	ErrorCategoryImport ErrorCategory = "import"
)

var ErrRegistry = map[ErrorCategory]map[ErrorCode]ErrorCode{
	ErrorCategoryImport: {
		ErrorCodeInconsistentTaxa: ErrorCodeInconsistentTaxa,
	},
}

type AppError struct {
	huma.ErrorModel
	Code     ErrorCode     `json:"code,omitempty"`
	Category ErrorCategory `json:"category,omitempty"`
}

type AppErrorDetail = huma.ErrorDetail

type AppErrorProvider interface {
	AppError() *AppError
}

func InternalError(err error) *AppError {
	return &AppError{
		ErrorModel: huma.ErrorModel{
			Status: 500,
			Title:  "Internal server error",
			Detail: err.Error(),
		},
	}
}

func NotFoundError(err error) *AppError {
	return &AppError{
		ErrorModel: huma.ErrorModel{
			Status: 404,
			Title:  "Not found",
			Detail: err.Error(),
		},
	}
}

func ForbiddenError(err error) *AppError {
	return &AppError{
		ErrorModel: huma.ErrorModel{
			Status: 403,
			Title:  "Forbidden",
			Detail: err.Error(),
		},
	}
}

func AsAppError(err error) *AppError {
	logrus.Infof("AsAppError: %v", err)
	if err == nil {
		return nil
	}

	var provider AppErrorProvider
	if errors.As(err, &provider) {
		return provider.AppError()
	}

	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		logrus.Infof("AsAppError: returning NotFoundError for err: %v", err)
		return NotFoundError(err)
	}

	return InternalError(err)
}
