package cli

import (
	"fmt"
	"log"

	"github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Charger la configuration chargée globalement via cmd.GetConfig()
		cfg := cmd.GetConfig()

		// Initialiser la connexion à la base de données SQLite avec GORM.
		var DB *gorm.DB
		var err error

		log.Printf("Tentative de connexion à la base de données : %s", cfg.Database.Path) // Correction: Path au lieu de Name

		DB, err = gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{}) // Correction: Path au lieu de Name
		if err != nil {
			log.Fatalf("Échec de la connexion à la base de données '%s': %v", cfg.Database.Path, err) // Correction: Path au lieu de Name
		}

		log.Println("Connexion à la base de données SQLite réussie !")

		sqlDB, err := DB.DB() // Correction: DB au lieu de db
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}
		// Assurez-vous que la connexion est fermée après la migration.
		defer sqlDB.Close()

		// Exécuter les migrations automatiques de GORM.
		// Utilisez DB.AutoMigrate() et passez-lui les pointeurs vers tous vos modèles.
		err = DB.AutoMigrate(&models.Link{}, &models.Click{})
		if err != nil {
			log.Fatalf("Échec des migrations: %v", err)
		}

		// Pas touche au log
		fmt.Println("Migrations de la base de données exécutées avec succès.")
	},
}

func init() {
	// Ajouter la commande à RootCmd
	cmd.RootCmd.AddCommand(MigrateCmd)
}
