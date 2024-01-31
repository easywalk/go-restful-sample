package model

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null;unique;"`
	Size      uint64    `json:"size" gorm:"type:bigint;"`
	Type      string    `json:"type" gorm:"type:varchar(255);"`
	CreatedAt time.Time `json:"createdAt" gorm:"type:timestamp;"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"type:timestamp;"`
}

func (f *File) SetCreatedAt(t time.Time) {
	f.CreatedAt = t
}

func (f *File) SetUpdatedAt(t time.Time) {
	f.UpdatedAt = t
}

func (f *File) GetID() uuid.UUID {
	return f.ID
}

func (f *File) SetID(id uuid.UUID) {
	f.ID = id
}
