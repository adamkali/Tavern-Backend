package models

type AuthRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserEmail string `json:"user_email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RelationRequest struct {
	Self  string       `json:"self"`  // 32 char user uuid
	Other string       `json:"other"` // 32 char user uuid
	Type  Relationship `json:"type"`  // 32 char relationship uuid for enum
}

type VerifyRequest struct {
	UserEmail string `json:"user_email"`
	UserPhone string `json:"user_phone"`
	Pin       string `json:"pin"`
}

type TagRequest struct {
	UserID string `json:"user_id"`
	Tag    Tag    `json:"tag"`
}

type TagsRequest struct {
	UserID string `json:"user_id"`
	Tags   []Tag  `json:"tags"`
}

type PrefrenceRequest struct {
	UserID string          `json:"user_id"`
	Pref   PlayerPrefrence `json:"pref"`
}

type PrefrencesRequest struct {
	UserID string            `json:"user_id"`
	Prefs  []PlayerPrefrence `json:"prefs"`
}

type RoleRequest struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
}

type RolesRequest struct {
	UserID string `json:"user_id"`
	Roles  []Role `json:"roles"`
}

type RoleChangeRequest struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

type ReportRequest struct {
	ReporterID string `json:"reporter_id"`
	ReportedID string `json:"reported_id"`
	Reason     string `json:"reason"`
	Date       string `json:"date"`
}
