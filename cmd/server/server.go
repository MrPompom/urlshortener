package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/monitor"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/spf13/cobra"
	// Driver SQLite pour GORM
)

// RunServerCmd représente la commande 'run-server' de Cobra.
// C'est le point d'entrée pour lancer le serveur de l'application.
var DB *gorm.DB

var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
	démarre les workers asynchrones pour les clics et le moniteur d'URLs,
	puis lance le serveur HTTP.`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		if cmd.Cfg == nil {
			log.Fatalf("Configuration non chargée.")
		}

		var err error

		log.Printf("Tentative de connexion à la base de données : %s", cmd.Cfg.Database.Name)

		DB, err = gorm.Open(sqlite.Open(cmd.Cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("Échec de la connexion à la base de données '%s': %v", cmd.Cfg.Database.Name, err)
		}

		log.Println("Connexion à la base de données SQLite réussie !")

		// TODO : Initialiser les repositories.
		// Créez des instances de GormLinkRepository et GormClickRepository.

		// Laissez le log
		log.Println("Repositories initialisés.")

		// TODO : Initialiser les services métiers.
		// Créez des instances de LinkService et ClickService, en leur passant les repositories nécessaires.

		// Laissez le log
		log.Println("Services métiers initialisés.")

		// TODO : Initialiser le channel ClickEventsChannel (api/handlers) des événements de clic et lancer les workers (StartClickWorkers).
		// Le channel est bufferisé avec la taille configurée.
		// Passez le channel et le clickRepo aux workers.

		// TODO : Remplacer les XXX par les bonnes variables
		log.Printf("Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
			XXX, XXX)

		// TODO : Initialiser et lancer le moniteur d'URLs.
		// Utilisez l'intervalle configuré (cfg.Monitor.IntervalMinutes).
		// Lancez le moniteur dans sa propre goroutine.
		monitorInterval := time.Duration(XXX) * time.Minute
		urlMonitor := monitor.NewUrlMonitor() // Le moniteur a besoin du linkRepo et de l'interval
		go urlMonitor.Start()
		log.Printf("Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)

		// TODO : Configurer le routeur Gin et les handlers API.
		// Passez les services nécessaires aux fonctions de configuration des routes.

		// Pas toucher au log
		log.Println("Routes API configurées.")

		// Créer le serveur HTTP Gin
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// TODO : Démarrer le serveur Gin dans une goroutine anonyme pour ne pas bloquer.
		// Pensez à logger des ptites informations...

		// Gére l'arrêt propre du serveur (graceful shutdown).
		// Créez un channel pour les signaux OS (SIGINT, SIGTERM).
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Attendre Ctrl+C ou signal d'arrêt

		// Bloquer jusqu'à ce qu'un signal d'arrêt soit reçu.
		<-quit
		log.Println("Signal d'arrêt reçu. Arrêt du serveur...")

		// Arrêt propre du serveur HTTP avec un timeout.
		log.Println("Arrêt en cours... Donnez un peu de temps aux workers pour finir.")
		time.Sleep(5 * time.Second)

		log.Println("Serveur arrêté proprement.")
	},
}

func init() {
	cmd.RootCmd.AddCommand(RunServerCmd)
}
