package repository

import (
	"ponderada-3/domain"
	"gorm.io/gorm"
)


type FigureRepository interface {
	Create(figurinha *domain.Figurinha) (*domain.Figurinha, error)
	FindAll(tipo domain.FigurinhaType, posicao domain.FigurinhaPosition) ([]domain.Figurinha, error)
	FindByID(id uint) (*domain.Figurinha, error)
	Update(figurinha *domain.Figurinha) (*domain.Figurinha, error)
	Delete(id uint) error
}

type figureRepository struct {
	db *gorm.DB
}

func NewFigureRepository(db *gorm.DB) FigureRepository {
	return &figureRepository{db: db}
}


func (r *figureRepository) Create(figurinha *domain.Figurinha) (*domain.Figurinha, error) {
	if err := r.db.Create(figurinha).Error; err != nil {
		return nil, err
	}

	return figurinha, nil
}



func (r *figureRepository) FindAll(tipo domain.FigurinhaType, posicao domain.FigurinhaPosition) ([]domain.Figurinha, error) {
	var figurinhas []domain.Figurinha

	query := r.db.Order("created_at DESC")

	if tipo != "" {
		query = query.Where("tipo = ?", tipo)
	}

	if posicao != "" {
		query = query.Where("posicao = ?", posicao)
	}

	if err := query.Find(&figurinhas).Error; err != nil {
		return nil, err
	}

	return figurinhas, nil
}

func (r *figureRepository) FindByID(id uint) (*domain.Figurinha, error) {
	var figurinha domain.Figurinha

	if err := r.db.First(&figurinha, id).Error; err != nil {
		return nil, err
	}

	return &figurinha, nil
}

func (r *figureRepository) Update(figurinha *domain.Figurinha) (*domain.Figurinha, error) {
	if err := r.db.Save(figurinha).Error; err != nil {
		return nil, err
	}

	return figurinha, nil
}

func (r *figureRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Figurinha{}, id).Error
}