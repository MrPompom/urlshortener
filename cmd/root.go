package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/axellelanca/urlshortener/internal/config"
	"github.com/spf13/cobra"
)

var Cfg *config.Config

// RootCmd représente la commande de base lorsque l'on appelle l'application sans sous-commande
var RootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
	Long: `url-shortener est une application complète pour gérer des URLs courtes.
Elle inclut un serveur API pour le raccourcissement et la redirection,
ainsi qu'une interface en ligne de commande pour l'administration.

Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.`,
}

// Execute est le point d'entrée principal pour l'application Cobra.
// Il est appelé depuis 'main.go'.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution de la commande: %v\n", err)
		os.Exit(1)
	}
}

// GetConfig retourne l'instance de configuration globale
// Cette fonction est utilisée par les autres commandes pour accéder à la config
func GetConfig() *config.Config {
	if Cfg == nil {
		log.Fatal("Configuration non initialisée. Appelez initConfig() d'abord.")
	}
	return Cfg
}

// init() est une fonction spéciale de Go qui s'exécute automatiquement
// avant la fonction main(). Elle est utilisée ici pour initialiser Cobra
// et ajouter toutes les sous-commandes.
func init() {
	cobra.OnInitialize(initConfig)
	// IMPORTANT : Ici, nous n'appelons PAS RootCmd.AddCommand() directement
	// pour les commandes 'server', 'create', 'stats', 'migrate'.
	// Ces commandes s'enregistreront elles-mêmes via leur propre fonction init().
	//
	// Assurez-vous que tous les fichiers de commande comme
	// 'cmd/server/server.go' et 'cmd/cli/*.go' aient bien
	// un `import "github.com/axellelanca/urlshortener/cmd"`
	// et un `func init() { cmd.RootCmd.AddCommand(MaCommandeCmd) }`
	// C'est ce qui va faire le lien !
}

// initConfig charge la configuration de l'application.
// Cette fonction est appelée au début de l'exécution de chaque commande Cobra
// grâce à `cobra.OnInitialize(initConfig)`.
func initConfig() {
	var err error
	Cfg, err = config.Load() // Correction: utiliser config.Load() au lieu de config.LoadConfig()
	if err != nil {
		// Loggue l'erreur mais ne fait pas un os.Exit(1) ici si Load()
		// gère déjà l'absence de fichier avec des valeurs par défaut.
		// Si Load() termine le programme en cas d'erreur fatale,
		// cette vérification est surtout pour les avertissements.
		log.Printf("Attention: Problème lors du chargement de la configuration: %v. Utilisation des valeurs par défaut.", err)
	}
	// La configuration est maintenant disponible via la variable globale 'cmd.Cfg'
	// et accessible via cmd.GetConfig() depuis les autres fichiers.
}
