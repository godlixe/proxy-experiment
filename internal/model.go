package internal

import "time"

type User struct {
	// Username is the user's name
	// which is defined by the header
	// x-username
	Username string `json:"username"`

	// Count is the number of hits
	// by Username field.
	Count int `json:"count"`
}

// AccessLog defines the accesses
// to the app.
type AccessLog struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	Username   string    `json:"username"`
	IPAddress  string    `json:"ip_address"`
	AccessTime time.Time `json:"access_time"`
}
