package controllers

import (
	"Tavern-Backend/models"

	"gorm.io/gorm"
)

type RelationshipController struct {
	H BaseHandler[models.UserRelationship]
}

func NewRelationshipController(DB *gorm.DB) *RelationshipController {
	return &RelationshipController{
		H: *NewHandler(DB, models.UserRelationship{}, "relationship"),
	}
}

// FIXME: #1 add in `/api/auth/Relationships/report` endpoint
// :ENDFIXME
