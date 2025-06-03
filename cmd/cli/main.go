package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"myapi/config"
	"myapi/internal/db"
	"myapi/internal/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [arguments...]")
		fmt.Println("Available commands:")
		fmt.Println("  update_pwd <email> <new_password>")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Handle commands
	command := os.Args[1]
	switch command {
	case "update_pwd":
		if len(os.Args) != 4 {
			fmt.Println("Usage: cli update_pwd <email> <new_password>")
			os.Exit(1)
		}
		email := os.Args[2]
		newPassword := os.Args[3]
		update_pwd(database, email, newPassword)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// update_pwd updates a user's password by email
func update_pwd(db *gorm.DB, email, new_pwd string) {
	var account models.CoachAccount

	// Find the account with the email
	result := db.Where("provider_type = ? AND provider_id = ?",
		models.AccountProviderTypeEmailWithPwd, email).First(&account)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("No account found with email: %s\n", email)
			os.Exit(1)
		}
		log.Fatalf("Database error: %v", result.Error)
		return
	}

	// Hash the new password
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(new_pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Update the password
	account.ProviderArg1 = string(hashed_pwd)
	result = db.Model(&account).Where("provider_id = ?", email).Updates(&account)
	if result.Error != nil {
		log.Fatalf("Failed to update password: %v", result.Error)
		return
	}

	fmt.Printf("Password updated successfully for email: %s\n", email)
}
