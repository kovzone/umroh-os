// helpers.go — shared helper utilities for the ops-svc sqlc package.

package sqlc

import "github.com/jackc/pgx/v5/pgtype"

// nullText converts an empty string to a null pgtype.Text.
func nullText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}
