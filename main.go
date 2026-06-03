package main

import (
	domain "ponderada-3/domain"
	"ponderada-3/handler"
	"ponderada-3/repository"
	"ponderada-3/service"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. Banco de dados — erro fatal aqui, pois sem banco nada funciona
	db, err := gorm.Open(sqlite.Open("figurinhas.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("erro ao abrir banco de dados: %v", err)
	}

	if err := db.AutoMigrate(&domain.Figurinha{}); err != nil {
		log.Fatalf("erro na migração: %v", err)
	}

	repo := repository.NewFigureRepository(db)
	svc := service.NewFigureService(repo)
	h := handler.NewFigureHandler(svc)

	r := gin.Default()

	h.RegisterRoutes(r)

	log.Println("Servidor iniciado em http://localhost:8080")
	if err := r.Run(":42069"); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
