package monitor

import (
	"log"
	"net/http"
	"sync" // Pour protéger l'accès concurrentiel à knownStates
	"time"

	_ "github.com/axellelanca/urlshortener/internal/models"   // Importe les modèles de liens
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le repository de liens
)

// UrlMonitor gère la surveillance périodique des URLs longues.
type UrlMonitor struct {
	linkRepo    repository.LinkRepository // Pour récupérer les URLs à surveiller
	interval    time.Duration             // Intervalle entre chaque vérification (ex: 5 minutes)
	knownStates map[uint]bool             // État connu de chaque URL: map[LinkID]estAccessible (true/false)
	mu          sync.Mutex                // Mutex pour protéger l'accès concurrentiel à knownStates
}

// NewUrlMonitor crée et retourne une nouvelle instance de UrlMonitor.
// Attention: retourne un pointeur
func NewUrlMonitor(linkRepo repository.LinkRepository, interval time.Duration) *UrlMonitor {
	return &UrlMonitor{
		linkRepo:    linkRepo,            // Injecte le repository de liens pour récupérer les URLs à surveiller
		interval:    interval,            // Définit l'intervalle de vérification
		knownStates: make(map[uint]bool), // Initialise la map pour stocker les états connus des URLs
		mu:          sync.Mutex{},        // Initialise le mutex pour protéger l'accès concurrentiel
	}
}

// Start lance la boucle de surveillance périodique des URLs.
// Cette fonction est conçue pour être lancée dans une goroutine séparée.
func (m *UrlMonitor) Start() {
	log.Printf("[MONITOR] Démarrage du moniteur d'URLs avec un intervalle de %v...", m.interval)
	ticker := time.NewTicker(m.interval) // Crée un ticker qui envoie un signal à chaque intervalle
	defer ticker.Stop()                  // S'assure que le ticker est arrêté quand Start se termine

	// Exécute une première vérification immédiatement au démarrage
	m.checkUrls()

	// Boucle principale du moniteur, déclenchée par le ticker
	for range ticker.C {
		m.checkUrls()
	}
}

// checkUrls effectue une vérification de l'état de toutes les URLs longues enregistrées.
func (m *UrlMonitor) checkUrls() {
	log.Println("[MONITOR] Lancement de la vérification de l'état des URLs...")

	// Gérer l'erreur si la récupération échoue.
	// Si erreur : log.Printf("[MONITOR] ERREUR lors de la récupération des liens pour la surveillance : %v", err)
	links, err := m.linkRepo.GetAllLinks()
	if err != nil {
		log.Printf("[MONITOR] ERREUR lors de la récupération des liens pour la surveillance : %v", err)
		return // Sort de la fonction si une erreur se produit
	}

	for _, link := range links {
		currentState := m.isUrlAccessible(link.LongURL)
		if currentState {
			log.Printf("[MONITOR] L'URL %s (%s) est ACCESSIBLE",
				link.Shortcode, link.LongURL)
		} else {
			log.Printf("[MONITOR] L'URL %s (%s) est INACCESSIBLE",
				link.Shortcode, link.LongURL)
		}

		// Protéger l'accès à la map 'knownStates' car 'checkUrls' peut être exécuté concurremment
		m.mu.Lock()
		previousState, exists := m.knownStates[link.ID] // Récupère l'état précédent
		m.knownStates[link.ID] = currentState           // Met à jour l'état actuel
		m.mu.Unlock()

		// Si c'est la première vérification pour ce lien, on initialise l'état sans notifier.
		if !exists {
			log.Printf("[MONITOR] État initial pour le lien %s (%s) : %s",
				link.Shortcode, link.LongURL, formatState(currentState))
			continue
		}

		// Si l'état a changé, générer une fausse notification dans les logs.
		// log.Printf("[NOTIFICATION] Le lien %s (%s) est passé de %s à %s !"
		if currentState != previousState {
			log.Printf("[NOTIFICATION] Le lien %s (%s) est passé de %s à %s !",
				link.Shortcode, link.LongURL,
				formatState(previousState), formatState(currentState))
		}
	}
	log.Println("[MONITOR] Vérification de l'état des URLs terminée.")
}

// isUrlAccessible effectue une requête HTTP HEAD pour vérifier l'accessibilité d'une URL.
func (m *UrlMonitor) isUrlAccessible(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second, // Timeout de 5 secondes pour la requête HTTP
	}
	// Un code de statut 2xx ou 3xx indique que l'URL est accessible.
	// Si err : log.Printf("[MONITOR] Erreur d'accès à l'URL '%s': %v", url, err)
	resp, err := client.Head(url)
	if err != nil {
		log.Printf("[MONITOR] Erreur d'accès à l'URL '%s': %v", url, err)
		return false // Si une erreur se produit, on considère l'URL comme inaccessible
	}

	defer resp.Body.Close() // Assurez-vous de fermer le corps de la réponse pour libérer les ressources
	log.Printf("[MONITOR] Requête HEAD pour l'URL '%s' a renvoyé le code de statut %d", url, resp.StatusCode)
	// Déterminer l'accessibilité basée sur le code de statut HTTP.
	return resp.StatusCode >= 200 && resp.StatusCode < 400 // Codes 2xx ou 3xx
}

// formatState est une fonction utilitaire pour rendre l'état plus lisible dans les logs.
func formatState(accessible bool) string {
	if accessible {
		return "ACCESSIBLE"
	}
	return "INACCESSIBLE"
}
