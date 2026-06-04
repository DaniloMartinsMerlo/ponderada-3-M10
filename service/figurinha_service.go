package service

import (
	"ponderada-3/domain"
	"ponderada-3/repository"

	"errors"
	"time"
	"gorm.io/gorm"
	"regexp"
)

var numeroRegex = regexp.MustCompile(`^[A-Za-z]{3} \d{2}$`)

var (
	ErrFigureNotFound      = errors.New("figurinha não encontrada")
	ErrInvalidNumber       = errors.New("o número informado é invalido, precisa ter o padrão: XXX 00, considere xxx é equivalente a federação e 00 é o número do jogador")
	ErrInvalidType         = errors.New("o tipo informado é invalido")
	ErrInvalidPosition     = errors.New("a posição informada é invalida")
)

type FigureService interface {
	Create(req domain.CreateFigureRequest) (*domain.Figurinha, error)
	ListAll(req domain.FindAllFigureRequest) ([]domain.Figurinha, error)
	GetByID(id uint) (*domain.Figurinha, error)
	Update(id uint, req domain.UpdateFigureRequest) (*domain.Figurinha, error)
	Delete(id uint) error
}

type figureService struct {
	repo repository.FigureRepository
}

func NewFigureService(repo repository.FigureRepository) FigureService {
	return &figureService{repo: repo}
}

func (s *figureService) Create(req domain.CreateFigureRequest) (*domain.Figurinha, error) {
	
	if err := validateNumero(req.Numero); err != nil {
		return nil, err
	}

	if !domain.ValidType[req.Tipo] {
		return nil, ErrInvalidType
	}

	if !domain.ValidPosition[req.Posicao] {
		return nil, ErrInvalidPosition
	}
	
	figurinha := &domain.Figurinha{
		Numero:    req.Numero,
		Tipo:      req.Tipo,
		Posicao:   req.Posicao,
		UpdateAt:  time.Now(),
		CreatedAt: time.Now(),
	}

	created, err := s.repo.Create(figurinha)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *figureService) ListAll(req domain.FindAllFigureRequest) ([]domain.Figurinha, error) {

	if req.Tipo != "" {
		if !domain.ValidType[req.Tipo] {
			return nil, ErrInvalidType
		}
	}

	if req.Posicao != "" {
		if !domain.ValidPosition[req.Posicao] {
			return nil, ErrInvalidPosition
		}
	}

	return s.repo.FindAll(req.Tipo, req.Posicao)
}

func (s *figureService) GetByID(id uint) (*domain.Figurinha, error) {
	figurinha, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFigureNotFound
		}
		return nil, err
	}
	return figurinha, nil
}

func (s *figureService) Update(id uint, req domain.UpdateFigureRequest) (*domain.Figurinha, error) {
	figurinha, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Numero != "" {

		if err := validateNumero(req.Numero); err != nil {
    		return nil, err
		}

		figurinha.Numero = req.Numero
	}

	if req.Tipo != "" {
		if !domain.ValidType[req.Tipo] {
			return nil, ErrInvalidType
		}
		figurinha.Tipo = req.Tipo
	}

	if req.Posicao != "" {
		if !domain.ValidPosition[req.Posicao] {
			return nil, ErrInvalidPosition
		}
		figurinha.Posicao = req.Posicao
	}

	figurinha.UpdateAt = time.Now()

	updated, err := s.repo.Update(figurinha)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *figureService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func validateNumero(numero string) error {
    if !numeroRegex.MatchString(numero) {
        return ErrInvalidNumber
    }
    return nil
}