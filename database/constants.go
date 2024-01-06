package database

const (
	TableUsers    string = "users"
	TableSessions string = "sessions"
)

const (
	ColumnId           string = "id"
	ColumnUserId       string = "user_id"
	ColumnEmail        string = "email"
	ColumnPasswordHash string = "password_hash"
	ColumnTokenHash    string = "token_hash"
)
