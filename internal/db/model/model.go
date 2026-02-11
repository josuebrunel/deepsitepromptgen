package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type UserPrompt struct {
	ID        uuid.UUID `db:"id" json:"id" form:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id" form:"user_id"`
	Name      string    `db:"name" json:"name" form:"name"`
	Content   string    `db:"content" json:"content" form:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
