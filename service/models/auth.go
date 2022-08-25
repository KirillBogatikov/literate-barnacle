package models

import (
	"github.com/Viva-Victoria/bear-jwt"
	"github.com/google/uuid"
)

type Authorization struct {
	UserId uuid.UUID `json:"userId"`
	Role   Role      `json:"role"`
}

type TokenClaims struct {
	jwt.BasicClaims
	Authorization
}
