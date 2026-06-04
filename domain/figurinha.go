package domain

import "time"

type FigureType string

const (
	TypeCommun       FigureType = "comum"
	TypeShine        FigureType = "brilhante"
	TypeLegendGold   FigureType = "legends_ouro"
	TypeLegendCopper FigureType = "legends_bronze"
)

type FigurePosition string

const (
	PositionGoalkeeper FigurePosition = "goleiro"
	PositionDefender   FigurePosition = "zagueiro"
	PositionMidfielder FigurePosition = "meio-campista"
	PositionForward    FigurePosition = "atacante"
)

var ValidType = map[FigureType]bool{
	TypeCommun:       true,
	TypeShine:        true,
	TypeLegendGold:   true,
	TypeLegendCopper: true,
}

var ValidPosition = map[FigurePosition]bool{
	PositionGoalkeeper: true,
	PositionDefender:   true,
	PositionMidfielder: true,
	PositionForward:    true,
}

type Figurinha struct {
	ID        uint           `json:"id"          gorm:"primaryKey"`
	Numero    string         `json:"numero"      gorm:"not null"`
	Tipo      FigureType     `json:"tipo"        gorm:"not null"`
	Posicao   FigurePosition `json:"posicao"     gorm:"not null"`
	UpdateAt  time.Time      `json:"updated_at"   gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"  gorm:"not null"`
}

type CreateFigureRequest struct {
	Numero  string         `json:"numero"    binding:"required"`
	Tipo    FigureType     `json:"tipo"      binding:"required"`
	Posicao FigurePosition `json:"posicao"   binding:"required"`
}

type UpdateFigureRequest struct {
	Numero  string         `json:"numero"    binding:"omitempty"`
	Tipo    FigureType     `json:"tipo"      binding:"omitempty"`
	Posicao FigurePosition `json:"posicao"   binding:"omitempty"`
}

type FindAllFigureRequest struct {
	Tipo    FigureType     `json:"tipo"      binding:"omitempty"`
	Posicao FigurePosition `json:"posicao"   binding:"omitempty"`
}
