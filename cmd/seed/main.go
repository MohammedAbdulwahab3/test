package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"flutter_go_crud/backend/internal/config"
	"flutter_go_crud/backend/internal/database"
	"flutter_go_crud/backend/internal/models"
)

func main() {
	cfg := config.Load()
	db, closeDB, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer closeDB()

	// Auto migrate first to ensure tables exist
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Seeding database...")

	// Sample data
	items := []models.Item{
		{
			Name:        "MacBook Pro 16",
			Description: "M3 Max chip, 32GB RAM, 1TB SSD. The ultimate pro laptop.",
			Price:       2499.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "iPhone 15 Pro",
			Description: "Titanium design, A17 Pro chip, 48MP Main camera.",
			Price:       999.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "Industry-leading noise canceling headphones with premium sound.",
			Price:       348.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "Dell XPS 13",
			Description: "Compact 13-inch laptop with InfinityEdge display.",
			Price:       1199.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "iPad Air",
			Description: "M1 chip, 10.9-inch Liquid Retina display.",
			Price:       599.00,
			Status:      "inactive",
			Version:     1,
		},
		{
			Name:        "Samsung Galaxy S24 Ultra",
			Description: "AI-powered smartphone with S Pen and 200MP camera.",
			Price:       1299.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "Logitech MX Master 3S",
			Description: "Performance wireless mouse with ultra-fast scrolling.",
			Price:       99.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "Keychron Q1 Pro",
			Description: "Wireless custom mechanical keyboard with aluminum body.",
			Price:       199.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "Herman Miller Aeron",
			Description: "Ergonomic office chair with Pellicle suspension.",
			Price:       1695.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "LG C3 OLED TV",
			Description: "42-inch 4K Smart TV with evo technology.",
			Price:       899.00,
			Status:      "inactive",
			Version:     1,
		},
		{
			Name:        "PlayStation 5",
			Description: "Next-gen gaming console with haptic feedback controller.",
			Price:       499.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "Nintendo Switch OLED",
			Description: "Gaming console with 7-inch OLED screen.",
			Price:       349.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "AirPods Pro 2",
			Description: "Active Noise Cancellation and Transparency mode.",
			Price:       249.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "Kindle Paperwhite",
			Description: "Waterproof e-reader with 6.8-inch display.",
			Price:       139.00,
			Status:      "active",
			Version:     1,
		},
		{
			Name:        "Dyson V15 Detect",
			Description: "Cordless vacuum with laser dust detection.",
			Price:       749.00,
			Status:      "inactive",
			Version:     1,
		},
	}

	// Add some random dates
	for i := range items {
		daysAgo := rand.Intn(30)
		items[i].CreatedAt = time.Now().AddDate(0, 0, -daysAgo)
		items[i].UpdatedAt = items[i].CreatedAt
	}

	// Create items
	count := 0
	for _, item := range items {
		if err := db.Create(&item).Error; err != nil {
			log.Printf("Failed to create item %s: %v", item.Name, err)
		} else {
			count++
		}
	}

	fmt.Printf("Successfully seeded %d items!\n", count)
}
