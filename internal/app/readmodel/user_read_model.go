package readmodel

import "time"

type UserListItem struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type UserProfile struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ProfilePicture *string    `json:"profile_picture,omitempty"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	PostCount      int64      `json:"post_count"`
	FollowerCount  int64      `json:"follower_count"`
}

type UserStats struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
	NewUsersToday int64 `json:"new_users_today"`
	NewUsersWeek  int64 `json:"new_users_week"`
	NewUsersMonth int64 `json:"new_users_month"`
}

type UserSearch struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`

	Score float64 `json:"score,omitempty"`
}

type UserActivity struct {
	UserID      int64     `json:"user_id"`
	Username    string    `json:"username"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IPAddress   string    `json:"ip_address,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
}
