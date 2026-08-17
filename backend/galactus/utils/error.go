package utils

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
)

func GormStatusError(error error, entityName string) error {
	if errors.Is(error, gorm.ErrRecordNotFound) {
		return status.Errorf(codes.NotFound, "%s not found", entityName)
	}
	if strings.Contains(error.Error(), "duplicate key") ||
		strings.Contains(error.Error(), "Duplicate entry") ||
		strings.Contains(error.Error(), "UNIQUE constraint failed") {
		return status.Errorf(codes.AlreadyExists, "%s already exists", entityName)
	}
	if strings.Contains(error.Error(), "foreign key constraint") ||
		strings.Contains(error.Error(), "FOREIGN KEY constraint failed") {
		return status.Errorf(codes.FailedPrecondition, "foreign key constraint error")
	}
	if strings.Contains(error.Error(), "check constraint") ||
		strings.Contains(error.Error(), "CHECK constraint failed") {
		return status.Errorf(codes.InvalidArgument, "validation failed: %v", error)
	}

	return status.Errorf(codes.Internal, "database error occurred")
}
