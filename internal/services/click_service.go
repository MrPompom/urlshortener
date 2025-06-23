package services

import (
	"errors"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// TODO : créer la struct
// ClickService est une structure qui fournit des méthodes pour la logique métier des clics.
// Elle est juste composer de clickRepo qui est de type ClickRepository
type ClickService struct {
	clickRepo repository.ClickRepository // Le ClickRepository pour interagir avec la base de données des clics.
}

// NewClickService crée et retourne une nouvelle instance de ClickService.
// C'est la fonction recommandée pour obtenir un service, assurant que toutes ses dépendances sont injectées.
func NewClickService(clickRepo repository.ClickRepository) *ClickService {
	return &ClickService{
		clickRepo: clickRepo,
	}
}

// RecordClick enregistre un nouvel événement de clic dans la base de données.
// Cette méthode est appelée par le worker asynchrone.
func (s *ClickService) RecordClick(click *models.Click) error {
	// TODO 1: Appeler le ClickRepository (CreateClick) pour créer l'enregistrement de clic.
	// Gérer toute erreur provenant du repository.
	if click == nil {
		return errors.New("click is nil") // Retourne une erreur si le clic est nil
	}
	return s.clickRepo.CreateClick(click) // Appelle le repository pour enregistrer le clic

}

// GetClicksCountByLinkID récupère le nombre total de clics pour un LinkID donné.
// Cette méthode pourrait être utilisée par le LinkService pour les statistiques, ou directement par l'API stats.
func (s *ClickService) GetClicksCountByLinkID(linkID uint) (int, error) {
	// TODO 2: Appeler le ClickRepository (CountclicksByLinkID) pour compter les clics par LinkID.
	if linkID == 0 {
		return 0, errors.New("linkID cannot be zero") // Retourne une erreur si linkID est zéro
	}
	count, err := s.clickRepo.CountClicksByLinkID(linkID) // Appelle le repository pour compter les clics
	if err != nil {
		return 0, err // Retourne l'erreur si la récupération échoue
	}
	return count, nil // Retourne le nombre de clics
}
