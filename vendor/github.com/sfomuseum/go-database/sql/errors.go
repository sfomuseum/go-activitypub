package sql

import (
	"fmt"
)

// WrapError returns a new error wrapping 'err' and prepending with the value of 't's Name() method.
func WrapError(t Table, err error) error {
	return fmt.Errorf("[%s] %w", t.Name(), err)
}

// InitializeTableError returns a new error with a default message for database initialization problems wrapping 'err' and prepending with the value of 't's Name() method.
func InitializeTableError(t Table, err error) error {
	return WrapError(t, fmt.Errorf("Failed to initialize database table, %w", err))
}

// DatabaseConnectionError returns a new error with a default message for database connection problems wrapping 'err' and prepending with the value of 't's Name() method.
func DatabaseConnectionError(t Table, err error) error {
	return WrapError(t, fmt.Errorf("Failed to establish database connection, %w", err))
}

// BeginTransactionError returns a new error with a default message for database transaction initialization problems wrapping 'err' and prepending with the value of 't's Name() method.
func BeginTransactionError(t Table, err error) error {
	return WrapError(t, fmt.Errorf("Failed to begin database transaction, %w", err))
}

// CommitTransactionError returns a new error with a default message for problems committing database transactions wrapping 'err' and prepending with the value of 't's Name() method.
func CommitTransactionError(t Table, err error) error {
	return WrapError(t, fmt.Errorf("Failed to commit database transaction, %w", err))
}

// PrepareStatementError returns a new error with a default message for problems preparing database (SQL) statements wrapping 'err' and prepending with the value of 't's Name() method.
func PrepareStatementError(t Table, err error) error {
	return WrapError(t, fmt.Errorf("Failed to prepare SQL statement, %w", err))
}

// ExecuteStatementError returns a new error with a default message for problems executing database (SQL) statements wrapping 'err' and prepending with the value of 't's Name() method.
func ExecuteStatementError(t Table, err error) error {
	return WrapError(t, fmt.Errorf("Failed to execute SQL statement, %w", err))
}
