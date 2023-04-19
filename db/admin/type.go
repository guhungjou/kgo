package admin

import "time"

const (
	AdminUserStatusOK     = "OK"
	AdminUserStatusBanned = "Banned"
)

type AdminUser struct {
	tableName struct{} `pg:"admin_user" json:"-"`

	ID          int64     `pg:"id,notnull,pk" json:"id"`
	Username    string    `pg:"username,notnull" json:"username"`
	Name        string    `pg:"name,notnull,use_zero" json:"name"`
	Salt        string    `pg:"salt,notnull,use_zero" json:"-"`
	PType       string    `pg:"ptype,notnull,use_zero" json:"-"`
	Password    string    `pg:"password,notnull,use_zero" json:"-"`
	Status      string    `pg:"status,notnull,use_zero" json:"status"`
	IsSuperuser bool      `pg:"is_superuser,notnull,use_zero" json:"is_superuser"`
	Phone       string    `pg:"phone,notnull,use_zero" json:"phone"`
	CreatedBy   int64     `pg:"created_by,notnull,use_zero" json:"created_by"`
	CreatedAt   time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt   time.Time `pg:"updated_at" json:"updated_at"`
}

const (
	AdminTokenStatusOK      = "OK"
	AdminTokenStatusInvalid = "Invalid"
)

type AdminToken struct {
	tableName struct{} `pg:"admin_token" json:"-"`

	ID     string `pg:"id,pk" json:"id"`
	UserID int64  `pg:"user_id,notnull" json:"user_id"`
	Device string `pg:"device,notnull,use_zero" json:"device"`
	IP     string `pg:"ip,notnull,use_zero" json:"ip"`
	Status string `pg:"status,notnull,use_zero" json:"status"`

	ExpiresAt time.Time `pg:"expires_at" json:"expires_at"`
	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	User *AdminUser `pg:"rel:has-one" json:"user,omitempty"`
}
