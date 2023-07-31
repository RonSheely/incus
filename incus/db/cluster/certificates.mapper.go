//go:build linux && cgo && !agent

package cluster

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cyphar/incus/incus/db/query"
	"github.com/cyphar/incus/shared/api"
)

var _ = api.ServerEnvironment{}

var certificateObjects = RegisterStmt(`
SELECT certificates.id, certificates.fingerprint, certificates.type, certificates.name, certificates.certificate, certificates.restricted
  FROM certificates
  ORDER BY certificates.fingerprint
`)

var certificateObjectsByID = RegisterStmt(`
SELECT certificates.id, certificates.fingerprint, certificates.type, certificates.name, certificates.certificate, certificates.restricted
  FROM certificates
  WHERE ( certificates.id = ? )
  ORDER BY certificates.fingerprint
`)

var certificateObjectsByFingerprint = RegisterStmt(`
SELECT certificates.id, certificates.fingerprint, certificates.type, certificates.name, certificates.certificate, certificates.restricted
  FROM certificates
  WHERE ( certificates.fingerprint = ? )
  ORDER BY certificates.fingerprint
`)

var certificateID = RegisterStmt(`
SELECT certificates.id FROM certificates
  WHERE certificates.fingerprint = ?
`)

var certificateCreate = RegisterStmt(`
INSERT INTO certificates (fingerprint, type, name, certificate, restricted)
  VALUES (?, ?, ?, ?, ?)
`)

var certificateDeleteByFingerprint = RegisterStmt(`
DELETE FROM certificates WHERE fingerprint = ?
`)

var certificateDeleteByNameAndType = RegisterStmt(`
DELETE FROM certificates WHERE name = ? AND type = ?
`)

var certificateUpdate = RegisterStmt(`
UPDATE certificates
  SET fingerprint = ?, type = ?, name = ?, certificate = ?, restricted = ?
 WHERE id = ?
`)

// certificateColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the Certificate entity.
func certificateColumns() string {
	return "certificates.id, certificates.fingerprint, certificates.type, certificates.name, certificates.certificate, certificates.restricted"
}

// getCertificates can be used to run handwritten sql.Stmts to return a slice of objects.
func getCertificates(ctx context.Context, stmt *sql.Stmt, args ...any) ([]Certificate, error) {
	objects := make([]Certificate, 0)

	dest := func(scan func(dest ...any) error) error {
		c := Certificate{}
		err := scan(&c.ID, &c.Fingerprint, &c.Type, &c.Name, &c.Certificate, &c.Restricted)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"certificates\" table: %w", err)
	}

	return objects, nil
}

// getCertificatesRaw can be used to run handwritten query strings to return a slice of objects.
func getCertificatesRaw(ctx context.Context, tx *sql.Tx, sql string, args ...any) ([]Certificate, error) {
	objects := make([]Certificate, 0)

	dest := func(scan func(dest ...any) error) error {
		c := Certificate{}
		err := scan(&c.ID, &c.Fingerprint, &c.Type, &c.Name, &c.Certificate, &c.Restricted)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"certificates\" table: %w", err)
	}

	return objects, nil
}

// GetCertificates returns all available certificates.
// generator: certificate GetMany
func GetCertificates(ctx context.Context, tx *sql.Tx, filters ...CertificateFilter) ([]Certificate, error) {
	var err error

	// Result slice.
	objects := make([]Certificate, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = Stmt(tx, certificateObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"certificateObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.ID != nil && filter.Fingerprint == nil && filter.Name == nil && filter.Type == nil {
			args = append(args, []any{filter.ID}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, certificateObjectsByID)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"certificateObjectsByID\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(certificateObjectsByID)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"certificateObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Fingerprint != nil && filter.ID == nil && filter.Name == nil && filter.Type == nil {
			args = append(args, []any{filter.Fingerprint}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, certificateObjectsByFingerprint)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"certificateObjectsByFingerprint\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(certificateObjectsByFingerprint)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"certificateObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID == nil && filter.Fingerprint == nil && filter.Name == nil && filter.Type == nil {
			return nil, fmt.Errorf("Cannot filter on empty CertificateFilter")
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getCertificates(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getCertificatesRaw(ctx, tx, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"certificates\" table: %w", err)
	}

	return objects, nil
}

// GetCertificate returns the certificate with the given key.
// generator: certificate GetOne
func GetCertificate(ctx context.Context, tx *sql.Tx, fingerprint string) (*Certificate, error) {
	filter := CertificateFilter{}
	filter.Fingerprint = &fingerprint

	objects, err := GetCertificates(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"certificates\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "Certificate not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"certificates\" entry matches")
	}
}

// GetCertificateID return the ID of the certificate with the given key.
// generator: certificate ID
func GetCertificateID(ctx context.Context, tx *sql.Tx, fingerprint string) (int64, error) {
	stmt, err := Stmt(tx, certificateID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"certificateID\" prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, fingerprint)
	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, api.StatusErrorf(http.StatusNotFound, "Certificate not found")
	}

	if err != nil {
		return -1, fmt.Errorf("Failed to get \"certificates\" ID: %w", err)
	}

	return id, nil
}

// CertificateExists checks if a certificate with the given key exists.
// generator: certificate Exists
func CertificateExists(ctx context.Context, tx *sql.Tx, fingerprint string) (bool, error) {
	_, err := GetCertificateID(ctx, tx, fingerprint)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateCertificate adds a new certificate to the database.
// generator: certificate Create
func CreateCertificate(ctx context.Context, tx *sql.Tx, object Certificate) (int64, error) {
	// Check if a certificate with the same key exists.
	exists, err := CertificateExists(ctx, tx, object.Fingerprint)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"certificates\" entry already exists")
	}

	args := make([]any, 5)

	// Populate the statement arguments.
	args[0] = object.Fingerprint
	args[1] = object.Type
	args[2] = object.Name
	args[3] = object.Certificate
	args[4] = object.Restricted

	// Prepared statement to use.
	stmt, err := Stmt(tx, certificateCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"certificateCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"certificates\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"certificates\" entry ID: %w", err)
	}

	return id, nil
}

// DeleteCertificate deletes the certificate matching the given key parameters.
// generator: certificate DeleteOne-by-Fingerprint
func DeleteCertificate(ctx context.Context, tx *sql.Tx, fingerprint string) error {
	stmt, err := Stmt(tx, certificateDeleteByFingerprint)
	if err != nil {
		return fmt.Errorf("Failed to get \"certificateDeleteByFingerprint\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(fingerprint)
	if err != nil {
		return fmt.Errorf("Delete \"certificates\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "Certificate not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d Certificate rows instead of 1", n)
	}

	return nil
}

// DeleteCertificates deletes the certificate matching the given key parameters.
// generator: certificate DeleteMany-by-Name-and-Type
func DeleteCertificates(ctx context.Context, tx *sql.Tx, name string, certificateType CertificateType) error {
	stmt, err := Stmt(tx, certificateDeleteByNameAndType)
	if err != nil {
		return fmt.Errorf("Failed to get \"certificateDeleteByNameAndType\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(name, certificateType)
	if err != nil {
		return fmt.Errorf("Delete \"certificates\": %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	return nil
}

// UpdateCertificate updates the certificate matching the given key parameters.
// generator: certificate Update
func UpdateCertificate(ctx context.Context, tx *sql.Tx, fingerprint string, object Certificate) error {
	id, err := GetCertificateID(ctx, tx, fingerprint)
	if err != nil {
		return err
	}

	stmt, err := Stmt(tx, certificateUpdate)
	if err != nil {
		return fmt.Errorf("Failed to get \"certificateUpdate\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(object.Fingerprint, object.Type, object.Name, object.Certificate, object.Restricted, id)
	if err != nil {
		return fmt.Errorf("Update \"certificates\" entry failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query updated %d rows instead of 1", n)
	}

	return nil
}
