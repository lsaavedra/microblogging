package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tweet struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	CreatedBy       uuid.UUID //TODO se relaciona con el user que lo creo
	Message         string    // TODO: hacer un check para validar que el largo no sea mayor a 280 char
	MediaContentUrl string    //TODO: modelarlo como un archivo o url que nos pasan donde queda el contenido
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}
