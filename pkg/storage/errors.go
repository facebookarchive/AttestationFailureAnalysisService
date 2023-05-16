package storage

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// ErrInitMySQL implements "error", for the description see Error.
type ErrInitMySQL struct {
	Err error
	DSN string
}

func (err ErrInitMySQL) Error() string {
	return fmt.Sprintf("unable to initialize a MySQL client (DSN: '%s'): %v", err.DSN, err.Err)
}

func (err ErrInitMySQL) Unwrap() error {
	return err.Err
}

// ErrMySQLPing implements "error", for the description see Error.
type ErrMySQLPing struct {
	Err error
}

func (err ErrMySQLPing) Error() string {
	return fmt.Sprintf("unable to ping the MySQL server: %v", err.Err)
}

func (err ErrMySQLPing) Unwrap() error {
	return err.Err
}

// ErrUnableToUpload implements "error", for the description see Error.
type ErrUnableToUpload struct {
	Key []byte
	Err error
}

func (err ErrUnableToUpload) Error() string {
	return fmt.Sprintf("unable to upload file '%X': %v", err.Key, err.Err)
}

func (err ErrUnableToUpload) Unwrap() error {
	return err.Err
}

// ErrUnableToUpdate implements "error", for the description see Error.
type ErrUnableToUpdate struct {
	insertedValue string
	Err           error
}

func (err ErrUnableToUpdate) Error() string {
	return fmt.Sprintf("unable to insert '%s' to the metadata table: %v",
		err.insertedValue, err.Err)
}

func (err ErrUnableToUpdate) Unwrap() error {
	return err.Err
}

// ErrUnableToInsert implements "error", for the description see Error.
type ErrUnableToInsert struct {
	insertedValue string
	Err           error
}

func (err ErrUnableToInsert) Error() string {
	return fmt.Sprintf("unable to insert '%s' to the metadata table: %v",
		err.insertedValue, err.Err)
}

func (err ErrUnableToInsert) Unwrap() error {
	return err.Err
}

// ErrAlreadyExists implements "error", for the description see Error.
type ErrAlreadyExists struct {
	insertedValue string
	Err           *mysql.MySQLError
}

func (err ErrAlreadyExists) Error() string {
	return fmt.Sprintf("image '%s' is already inserted to the metadata table: %v",
		err.insertedValue, err.Err)
}

func (err ErrAlreadyExists) Unwrap() error {
	return err.Err
}

// ErrUnableToUpdateMetadata implements "error", for the description see Error.
type ErrUnableToUpdateMetadata struct {
	Err error
}

func (err ErrUnableToUpdateMetadata) Error() string {
	return fmt.Sprintf("unable to update metadata record: %v", err.Err)
}

func (err ErrUnableToUpdateMetadata) Unwrap() error {
	return err.Err
}

// ErrGetMeta implements "error", for the description see Error.
type ErrGetMeta struct {
	Err error
}

func (err ErrGetMeta) Error() string {
	return fmt.Sprintf("unable to get the metadata record: %v", err.Err)
}

func (err ErrGetMeta) Unwrap() error {
	return err.Err
}

// ErrGetData implements "error", for the description see Error.
type ErrGetData struct {
	Err error
}

func (err ErrGetData) Error() string {
	return fmt.Sprintf("unable to get the data: %v", err.Err)
}

func (err ErrGetData) Unwrap() error {
	return err.Err
}

// ErrSelect implements "error", for the description see Error.
type ErrSelect struct {
	Err error
}

func (err ErrSelect) Error() string {
	return fmt.Sprintf("unable to select rows from MySQL: %v", err.Err)
}

func (err ErrSelect) Unwrap() error {
	return err.Err
}

// ErrNotFound implements "error", for the description see Error.
type ErrNotFound struct {
	Query string
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("not found (query: %s)", err.Query)
}

// ErrTooManyEntries implements "error", for the description see Error.
type ErrTooManyEntries struct {
	Count uint
}

func (err ErrTooManyEntries) Error() string {
	return fmt.Sprintf("too many entries: %d", err.Count)
}

// ErrDownload implements "error", for the description see Error.
type ErrDownload struct {
	Err error
}

func (err ErrDownload) Error() string {
	return fmt.Sprintf("unable to download: %v", err.Err)
}

func (err ErrDownload) Unwrap() error {
	return err.Err
}

// ErrEmptyFilters signals that search filters are empty and effectively
// the request requires to select all the data, which is forbidden.
type ErrEmptyFilters struct{}

func (err ErrEmptyFilters) Error() string {
	return "empty filters"
}
