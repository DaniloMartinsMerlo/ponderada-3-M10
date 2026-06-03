package domain

import "time"

type FigurinhaType string

const (
	TypeCommun FigurinhaType = "comum"
	TypeShine FigurinhaType = "brilhante"
	TypeLegendGold FigurinhaType = "legends_ouro"
	TypeLegendCopper FigurinhaType = "legends_bronze"
)

type FigurinhaPosition string

const (
	PositionGoalkeeper FigurinhaPosition = "Goleiro" 
	PositionDefender FigurinhaPosition = "Zagueiro"
	PositionMidfielder FigurinhaPosition = "Meio-campista"
	PositionForward FigurinhaPosition = "Atacante"
)

var ValidType = map[FigurinhaType]bool{
	TypeCommun:      true,
	TypeShine: true,
	TypeLegendGold:    true,
	TypeLegendCopper: true,
}

var ValidPosition = map[FigurinhaPosition]bool{
	PositionGoalkeeper:      true,
	PositionDefender: true,
	PositionMidfielder:    true,
	PositionForward: true,
}

type Figurinha struct {
	ID          uint              `json:"id"          gorm:"primaryKey"`
	Numero      string            `json:"numero"      gorm:"not null"`
	Tipo        FigurinhaType     `json:"tipo"        gorm:"not null"`
	Posicao     FigurinhaPosition `json:"posicao"     gorm:"not null"`
	UpdateAt    time.Time      	  `json:"updated_at"   gorm:"not null"`
	CreatedAt   time.Time         `json:"created_at"  gorm:"not null"` 
}


type CreateFigureRequest struct {
	Numero     string            `json:"numero"    binding:"required"`
	Tipo       FigurinhaType     `json:"tipo"      binding:"required"`
	Posicao    FigurinhaPosition `json:"posicao"   binding:"required"`
}

type UpdateFigureRequest struct {
	Numero     string            `json:"numero"    binding:"omitempty"`
	Tipo       FigurinhaType     `json:"tipo"      binding:"omitempty"`
	Posicao    FigurinhaPosition `json:"posicao"   binding:"omitempty"`
}

type FindAllFigureRequest struct {
	Tipo       FigurinhaType     `json:"tipo"      binding:"omitempty"`
	Posicao    FigurinhaPosition `json:"posicao"   binding:"omitempty"`
}