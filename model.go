package main

import (
	"fmt"
	"time"
)

type Member struct {
	MemberID       int       `db:"member_id" json:"member_id" pk:"member_id"`
	FirstName      string    `db:"first_name" json:"first_name"`
	LastName       string    `db:"last_name" json:"last_name"`
	Email          string    `db:"email" json:"email"`
	PasswordHash   string    `db:"password_hash" json:"password_hash"`
	DateOfBirth    date      `db:"date_of_birth" json:"date_of_birth"`
	JoinDate       timestamp `db:"join_date" json:"join_date"`
	MembershipType string    `db:"membership_type" json:"membership_type"`
	Status         string    `db:"status" json:"status"`
	CreatedAt      timestamp `db:"created_at" json:"created_at"`
	UpdatedAt      timestamp `db:"updated_at" json:"updated_at"`
}

func (m *Member) Fields() []any {
	return []any{&m.MemberID, &m.FirstName, &m.LastName, &m.Email, &m.PasswordHash, &m.DateOfBirth, &m.JoinDate, &m.MembershipType, &m.Status, &m.CreatedAt, &m.UpdatedAt}
}

type timestamp struct {
	time.Time
}

type date struct {
	time.Time
}

func (d *date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{} // Set to zero time if the value is nil
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *date", value)
	}
	d.Time = t
	return nil
}

func (d *timestamp) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{} // Set to zero time if the value is nil
		return nil
	}
	t, ok := value.(time.Time)
	fmt.Printf("\n%v\n", t)
	if !ok {
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *date", value)
	}
	d.Time = t
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface and parses the
// JSON encoding of the timestamp. The format is in RFC3339, the same as used by
// the time.Time type.
func (ct *timestamp) UnmarshalJSON(b []byte) error {
	// Remove the quotes surrounding the JSON string
	s := string(b)[1 : len(b)-1]

	// Parse the timestamp in RFC3339 format
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	// Set the timestamp to the parsed value
	ct.Time = t
	return nil
}

// MarshalJSON implements the json.Marshaler interface and returns the JSON
// encoding of the timestamp. The format is in RFC3339, the same as used by the
// time.Time type.
func (ct timestamp) MarshalJSON() ([]byte, error) {
	// Return the JSON encoding of the timestamp in RFC3339 format.
	// The quotes are necessary because the timestamp is a string
	return []byte(`"` + ct.Format(time.RFC3339) + `"`), nil
}

func (ct *date) UnmarshalJSON(b []byte) error {
	// Remove the quotes surrounding the JSON string
	s := string(b)[1 : len(b)-1]

	// Parse the DateOnly format
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}

	// Set the Date to the parsed value
	ct.Time = t
	return nil
}
func (ct date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Format(time.DateOnly) + `"`), nil
}
