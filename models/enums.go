package models

type Tag struct {
	ID      string `json:"id" gorm:"primaryKey; not null; type:varchar(32);"`
	TagID   int    `json:"tag_id" gorm:"column:tag_id;type:smallint(255);not null"`
	TagName string `json:"tag_name" gorm:"column:tag_name;varchar(32) not null"`
}

type PlayerPrefrence struct {
	ID           string `json:"id" gorm:"primaryKey; not null; type:varchar(32);"`
	PreferenceID int    `json:"pref_id" gorm:"column:pref_id;type:smallint(255);not null"`
	Preference   string `json:"pref_name" gorm:"column:pref_name;varchar(32) not null"`
}

type Relationship struct {
	ID               string `json:"id" gorm:"primaryKey; not null; type:varchar(32);"`
	RelationshipName string `json:"relationship_name" gorm:"column:relationship_name;varchar(32) not null"`
	Negative         bool   `json:"negative" gorm:"column:negative;type:tinyint(1);not null"`
}

type Tags struct {
	TempID string `json:"temp_id"`
	Tags   []Tag  `json:"tags"`
}

type PlayerPrefrences struct {
	TempID           string            `json:"temp_id"`
	PlayerPrefrences []PlayerPrefrence `json:"player_prefrences"`
}

type Relationships struct {
	TempID        string         `json:"temp_id"`
	Relationships []Relationship `json:"relationships"`
}

// Types of roles a user can have
// - Admin: Can do anything
// - LifetimePremium: Can access premium features for life
// - Premium: Can access premium features for a limited time
// - User: Can access basic features
// - Banned: Can only access Home Page
type Role struct {
	ID       string `json:"id" gorm:"primaryKey; not null; type:varchar(32);"`
	RoleName string `json:"role_name" gorm:"column:role_name;varchar(32) not null"`
}

// #region COMMON FUNCTIONS
func (t Tag) SetID(id string)             { t.ID = id }
func (t PlayerPrefrence) SetID(id string) { t.ID = id }
func (t Relationship) SetID(id string)    { t.ID = id }
func (t Role) SetID(id string)            { t.ID = id }
func (t Tag) GetID() string               { return t.ID }
func (t PlayerPrefrence) GetID() string   { return t.ID }
func (t Relationship) GetID() string      { return t.ID }
func (t Role) GetID() string              { return t.ID }

// #endregion
