package database

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/eneba-intern/marketplace/backend/internal/crypto"
	"github.com/eneba-intern/marketplace/backend/internal/models"
)

// SeedIfEmpty populates sellers, products, multi-seller offers, encrypted digital keys,
// and dynamic homepage block definitions if the database is empty.
func (s *Store) SeedIfEmpty() error {
	var count int
	if err := s.DB.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Seed Verified Marketplace Sellers
	sellers := []models.Seller{
		{Name: "KeyMaster Baltic", Slug: "keymaster-baltic", Verified: true, RatingPercent: 99.7, TotalSales: 148920, CountryCode: "LT", Badge: "Official Partner", JoinedYear: 2019},
		{Name: "NordicGames Vault", Slug: "nordicgames-vault", Verified: true, RatingPercent: 99.4, TotalSales: 94310, CountryCode: "SE", Badge: "Top Rated Seller", JoinedYear: 2020},
		{Name: "PixelForge Keys", Slug: "pixelforge-keys", Verified: true, RatingPercent: 98.9, TotalSales: 62400, CountryCode: "DE", Badge: "Instant Delivery", JoinedYear: 2021},
		{Name: "GameSphere EU", Slug: "gamesphere-eu", Verified: true, RatingPercent: 99.1, TotalSales: 51850, CountryCode: "PL", Badge: "Verified Merchant", JoinedYear: 2020},
		{Name: "HyperLoot Digital", Slug: "hyperloot-digital", Verified: true, RatingPercent: 98.2, TotalSales: 31900, CountryCode: "NL", Badge: "Fast Dispatches", JoinedYear: 2022},
		{Name: "Vanguard Codes", Slug: "vanguard-codes", Verified: true, RatingPercent: 99.5, TotalSales: 112500, CountryCode: "GB", Badge: "Platinum Merchant", JoinedYear: 2018},
	}

	var sellerIDs []int64
	for _, sel := range sellers {
		res, err := tx.Exec(`
			INSERT INTO sellers (name, slug, verified, rating_percent, total_sales, country_code, badge, joined_year)
			VALUES (?, ?, 1, ?, ?, ?, ?, ?)
		`, sel.Name, sel.Slug, sel.RatingPercent, sel.TotalSales, sel.CountryCode, sel.Badge, sel.JoinedYear)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		sellerIDs = append(sellerIDs, id)
	}

	// 2. Seed Catalog Products (Games, DLCs, Gift Cards, Subscriptions)
	type seedProduct struct {
		Slug        string
		Title       string
		Subtitle    string
		Description string
		Platform    string
		Region      string
		ProductType string
		Genre       string
		Publisher   string
		Developer   string
		ReleaseDate string
		CoverImage  string
		BannerImage string
		Screenshots []string
		BasePrice   float64
		BestPrice   float64
		Rating      float64
		ReviewCount int
		LikesCount  int
		SalesCount  int
		SystemReqs  models.SystemReq
		Steps       []string
	}

	catalog := []seedProduct{
		{
			Slug:        "cyberpunk-2077-ultimate-edition-steam-key-global",
			Title:       "Cyberpunk 2077: Ultimate Edition",
			Subtitle:    "Includes Base Game + Phantom Liberty Spy-Thriller Expansion",
			Description: "Cyberpunk 2077: Ultimate Edition is an open-world, action-adventure RPG set in the megalopolis of Night City, where you play as V, a cyberpunk mercenary wrapped up in a do-or-die fight for survival. Includes the acclaimed Phantom Liberty expansion set in Dogtown.",
			Platform:    "Steam",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Action RPG",
			Publisher:   "CD PROJEKT RED",
			Developer:   "CD PROJEKT RED",
			ReleaseDate: "2023-12-05",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1091500/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1091500/header.jpg",
			Screenshots: []string{
				"https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1091500/ss_b529b0abc43f55fc23fe8058eddb6e37c9629a6a.1920x1080.jpg",
				"https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1091500/ss_7924f64b6e5d586a80418c9896a1c92881a7905b.1920x1080.jpg",
			},
			BasePrice:   79.99,
			BestPrice:   34.49,
			Rating:      4.88,
			ReviewCount: 4210,
			LikesCount:  18940,
			SalesCount:  54200,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "Intel Core i7-12700 or Ryzen 7 7800X3D", Memory: "16 GB RAM", Graphics: "GeForce RTX 2060 SUPER / Radeon RX 5700 XT", Storage: "70 GB NVMe SSD"},
			Steps:       []string{"Launch the Steam desktop client and sign in", "Click 'Games' in the top-left menu and choose 'Activate a Product on Steam...'", "Paste your decrypted digital key from Eneba and click Finish"},
		},
		{
			Slug:        "elden-ring-shadow-of-the-erdtree-edition-steam-key-global",
			Title:       "Elden Ring: Shadow of the Erdtree Edition",
			Subtitle:    "Game of the Year RPG + Full Realm of Shadow Expansion",
			Description: "Rise, Tarnished, and be guided by grace to brandish the power of the Elden Ring and become an Elden Lord in the Lands Between. This complete bundle includes the base game and the massive Shadow of the Erdtree expansion.",
			Platform:    "Steam",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Action RPG",
			Publisher:   "Bandai Namco Entertainment",
			Developer:   "FromSoftware Inc.",
			ReleaseDate: "2024-06-21",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1245620/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1245620/header.jpg",
			Screenshots: []string{
				"https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1245620/ss_e80a9029b61143d2c7f929b47a4e4066a40292d0.1920x1080.jpg",
			},
			BasePrice:   79.99,
			BestPrice:   46.89,
			Rating:      4.95,
			ReviewCount: 6830,
			LikesCount:  29400,
			SalesCount:  68100,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "Intel Core i7-8700K or AMD Ryzen 5 3600X", Memory: "16 GB RAM", Graphics: "NVIDIA GeForce GTX 1070 8GB", Storage: "80 GB SSD"},
			Steps:       []string{"Open Steam client", "Select 'Activate a Product on Steam'", "Enter your 15-character Steam key"},
		},
		{
			Slug:        "the-witcher-3-wild-hunt-complete-edition-gog-key-global",
			Title:       "The Witcher 3: Wild Hunt - Complete Edition",
			Subtitle:    "Next-Gen Update + Hearts of Stone & Blood and Wine DLCs",
			Description: "You are Geralt of Rivia, mercenary monster slayer. Before you stands a war-torn, monster-infested continent you can explore at will. Your current contract? Tracking down Ciri — the Child of Prophecy, a living weapon that can alter the shape of the world.",
			Platform:    "GOG.com",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Action RPG",
			Publisher:   "CD PROJEKT RED",
			Developer:   "CD PROJEKT RED",
			ReleaseDate: "2022-12-14",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/292030/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/292030/header.jpg",
			Screenshots: []string{
				"https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/292030/ss_107600c1337accc09104f7a8aa7f275f23cad096.1920x1080.jpg",
			},
			BasePrice:   49.99,
			BestPrice:   11.49,
			Rating:      4.97,
			ReviewCount: 9120,
			LikesCount:  41200,
			SalesCount:  92400,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "Intel Core i5-7400 / Ryzen 5 1600", Memory: "8 GB RAM", Graphics: "GeForce GTX 1070 / Radeon RX 480", Storage: "55 GB SSD"},
			Steps:       []string{"Open GOG Galaxy or gog.com/redeem", "Paste the 18-character GOG activation code", "Confirm and download from your library"},
		},
		{
			Slug:        "helldivers-2-super-citizen-edition-steam-key-global",
			Title:       "HELLDIVERS 2",
			Subtitle:    "Galactic Co-Op Third-Person Squad Shooter",
			Description: "The Galaxy's Last Line of Offence. Enlist in the Helldivers and join the fight for freedom across a hostile galaxy in a fast, frantic, and ferocious third-person co-op shooter for up to 4 players.",
			Platform:    "Steam",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Shooter",
			Publisher:   "PlayStation Publishing LLC",
			Developer:   "Arrowhead Game Studios",
			ReleaseDate: "2024-02-08",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/553850/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/553850/header.jpg",
			Screenshots: []string{
				"https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/553850/ss_9da74db1314b211a847b278b0e207e712a988101.1920x1080.jpg",
			},
			BasePrice:   39.99,
			BestPrice:   26.99,
			Rating:      4.82,
			ReviewCount: 3890,
			LikesCount:  15800,
			SalesCount:  47500,
			SystemReqs:  models.SystemReq{OS: "Windows 10 64-bit", Processor: "Intel Core i7-9700K or AMD Ryzen 7 3700X", Memory: "16 GB RAM", Graphics: "NVIDIA GeForce RTX 2060 or AMD Radeon RX 6600XT", Storage: "100 GB SSD"},
			Steps:       []string{"Open Steam client", "Click 'Add a Game' -> 'Activate a Product on Steam'", "Enter key and deploy for Super Earth"},
		},
		{
			Slug:        "baldurs-gate-3-digital-deluxe-steam-key-global",
			Title:       "Baldur's Gate 3",
			Subtitle:    "Dungeons & Dragons Next-Generation Party CRPG",
			Description: "Gather your party and return to the Forgotten Realms in a tale of fellowship and betrayal, sacrifice and survival, and the lure of absolute power. Mysterious abilities are awakening inside you, drawn from a mind flayer parasite planted in your brain.",
			Platform:    "Steam",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Action RPG",
			Publisher:   "Larian Studios",
			Developer:   "Larian Studios",
			ReleaseDate: "2023-08-03",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1086940/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1086940/header.jpg",
			Screenshots: []string{
				"https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1086940/ss_c30dc70ed3b6b29a6f9c7f81d187e65c8c8bc7c0.1920x1080.jpg",
			},
			BasePrice:   59.99,
			BestPrice:   42.50,
			Rating:      4.96,
			ReviewCount: 5400,
			LikesCount:  26300,
			SalesCount:  51900,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "Intel i7 8700K / AMD r5 3600", Memory: "16 GB RAM", Graphics: "Nvidia 2060 Super / RX 5700 XT", Storage: "150 GB SSD"},
			Steps:       []string{"Open Steam client", "Select 'Activate a Product on Steam'", "Enter key and install"},
		},
		{
			Slug:        "red-dead-redemption-2-ultimate-edition-rockstar-key-global",
			Title:       "Red Dead Redemption 2: Ultimate Edition",
			Subtitle:    "Epic Outlaw Open-World Adventure + Story & Online Bonuses",
			Description: "America, 1899. Arthur Morgan and the Van der Linde gang are outlaws on the run. With federal agents and the best bounty hunters in the nation massing on their heels, the gang must rob, steal and fight their way across the rugged heartland of America in order to survive.",
			Platform:    "Rockstar Games",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Adventure",
			Publisher:   "Rockstar Games",
			Developer:   "Rockstar Games",
			ReleaseDate: "2019-12-05",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1174180/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1174180/header.jpg",
			Screenshots: []string{
				"https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1174180/ss_bac60bacbf5da8945103648c08d27d5e202444ca.1920x1080.jpg",
			},
			BasePrice:   89.99,
			BestPrice:   17.89,
			Rating:      4.93,
			ReviewCount: 7420,
			LikesCount:  33100,
			SalesCount:  83600,
			SystemReqs:  models.SystemReq{OS: "Windows 10 64-bit", Processor: "Intel Core i7-4770K / AMD Ryzen 5 1500X", Memory: "12 GB RAM", Graphics: "Nvidia GeForce GTX 1060 6GB", Storage: "150 GB"},
			Steps:       []string{"Open Rockstar Games Launcher", "Click your profile icon -> 'Redeem Code'", "Enter key and confirm"},
		},
		{
			Slug:        "grand-theft-auto-v-premium-online-edition-rockstar-key-global",
			Title:       "Grand Theft Auto V: Premium Online Edition",
			Subtitle:    "Complete GTA V Story + Criminal Enterprise Starter Pack",
			Description: "Experience Rockstar Games' critically acclaimed open world game, Grand Theft Auto V. Includes the complete GTA V story experience, free access to the ever-evolving Grand Theft Auto Online, and the Criminal Enterprise Starter Pack with GTA$1,000,000 bonus cash.",
			Platform:    "Rockstar Games",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Action",
			Publisher:   "Rockstar Games",
			Developer:   "Rockstar North",
			ReleaseDate: "2015-04-14",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/271590/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/271590/header.jpg",
			Screenshots: []string{},
			BasePrice:   29.99,
			BestPrice:   9.49,
			Rating:      4.87,
			ReviewCount: 11200,
			LikesCount:  48900,
			SalesCount:  124500,
			SystemReqs:  models.SystemReq{OS: "Windows 10 64-bit", Processor: "Intel Core i5 3470 @ 3.2GHz", Memory: "8 GB RAM", Graphics: "NVIDIA GTX 660 2GB / AMD HD 7870 2GB", Storage: "110 GB"},
			Steps:       []string{"Download Rockstar Games Launcher", "Navigate to Store -> Redeem Code", "Enter key to unlock GTA V Premium"},
		},
		{
			Slug:        "ea-sports-fc-25-standard-edition-ea-app-key-global",
			Title:       "EA SPORTS FC 25 Standard Edition",
			Subtitle:    "5v5 Rush Mode, FC IQ Tactical Overhaul & Ultimate Team",
			Description: "EA SPORTS FC 25 gives you more ways to win for the club. Team up with friends in your favourite modes with the new 5v5 Rush, and manage your club to victory as FC IQ delivers more tactical control than ever before.",
			Platform:    "EA App",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Sports",
			Publisher:   "Electronic Arts",
			Developer:   "EA Canada",
			ReleaseDate: "2024-09-27",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2669320/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2669320/header.jpg",
			Screenshots: []string{},
			BasePrice:   69.99,
			BestPrice:   28.90,
			Rating:      4.45,
			ReviewCount: 2150,
			LikesCount:  9400,
			SalesCount:  39200,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "AMD Ryzen 7 2700X or Intel Core i7 6700", Memory: "12 GB RAM", Graphics: "AMD RX 5600 XT or Nvidia GTX 1660", Storage: "100 GB"},
			Steps:       []string{"Launch the EA App on PC", "Go to Library -> Click 'Redeem Code'", "Paste your key and click Next"},
		},
		{
			Slug:        "black-myth-wukong-deluxe-edition-steam-key-global",
			Title:       "Black Myth: Wukong",
			Subtitle:    "Mythological Action RPG Rooted in Journey to the West",
			Description: "Black Myth: Wukong is an action RPG rooted in Chinese mythology. You shall set out as the Destined One to venture into the challenges and marvels ahead, to uncover the obscured truth beneath the veil of a glorious legend from the past.",
			Platform:    "Steam",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Action RPG",
			Publisher:   "Game Science",
			Developer:   "Game Science",
			ReleaseDate: "2024-08-20",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2358720/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2358720/header.jpg",
			Screenshots: []string{},
			BasePrice:   59.99,
			BestPrice:   44.99,
			Rating:      4.94,
			ReviewCount: 6100,
			LikesCount:  27800,
			SalesCount:  58300,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "Intel Core i7-9700 / AMD Ryzen 5 5500", Memory: "16 GB RAM", Graphics: "NVIDIA GeForce RTX 2060 / AMD Radeon RX 5700 XT", Storage: "130 GB SSD"},
			Steps:       []string{"Open Steam client", "Select 'Activate a Product on Steam'", "Enter key and install"},
		},
		{
			Slug:        "ghost-of-tsushima-directors-cut-steam-key-europe",
			Title:       "Ghost of Tsushima DIRECTOR'S CUT",
			Subtitle:    "Includes Iki Island Expansion & Legends Co-Op Multiplayer",
			Description: "In the late 13th century, the Mongol empire has laid waste to entire nations along their campaign to conquer the East. Tsushima Island is all that stands between mainland Japan and a massive Mongol invasion fleet led by Khotun Khan.",
			Platform:    "Steam",
			Region:      "EUROPE",
			ProductType: "Game Key",
			Genre:       "Adventure",
			Publisher:   "PlayStation Publishing LLC",
			Developer:   "Sucker Punch Productions / Nixxes",
			ReleaseDate: "2024-05-16",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2215430/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2215430/header.jpg",
			Screenshots: []string{},
			BasePrice:   59.99,
			BestPrice:   37.49,
			Rating:      4.91,
			ReviewCount: 3120,
			LikesCount:  14900,
			SalesCount:  33400,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "Intel Core i5-8600 or AMD Ryzen 5 3600", Memory: "16 GB RAM", Graphics: "NVIDIA GeForce RTX 2060 or AMD Radeon RX 5600 XT", Storage: "75 GB SSD"},
			Steps:       []string{"Open Steam client in a European region account", "Select 'Activate a Product on Steam'", "Paste your key"},
		},
		{
			Slug:        "god-of-war-ragnarok-deluxe-steam-key-global",
			Title:       "God of War Ragnarök",
			Subtitle:    "Kratos & Atreus Across the Nine Realms + Valhalla DLC",
			Description: "Kratos and Atreus must journey to each of the Nine Realms in search of answers as Asgardian forces prepare for a prophesied battle that will end the world. Includes the roguelite God of War Ragnarök: Valhalla DLC.",
			Platform:    "Steam",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Adventure",
			Publisher:   "PlayStation Publishing LLC",
			Developer:   "Santa Monica Studio",
			ReleaseDate: "2024-09-19",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2322010/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2322010/header.jpg",
			Screenshots: []string{},
			BasePrice:   59.99,
			BestPrice:   39.99,
			Rating:      4.85,
			ReviewCount: 2490,
			LikesCount:  12400,
			SalesCount:  29800,
			SystemReqs:  models.SystemReq{OS: "Windows 10 64-bit", Processor: "Intel i5-8600 or AMD Ryzen 5 3600", Memory: "16 GB RAM", Graphics: "NVIDIA RTX 2060 Super or AMD RX 5700", Storage: "190 GB SSD"},
			Steps:       []string{"Open Steam client", "Select 'Activate a Product on Steam'", "Enter key and download"},
		},
		{
			Slug:        "minecraft-java-and-bedrock-edition-pc-xbox-live-key-global",
			Title:       "Minecraft: Java & Bedrock Edition",
			Subtitle:    "Both Official PC Editions Included in One Digital License",
			Description: "Create, explore, survive, repeat. Minecraft: Java & Bedrock Edition for PC gives you access to both games in one purchase and one unified launcher, making it easier than ever to switch between editions and cross-play with any Minecrafter.",
			Platform:    "Xbox Live",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Adventure",
			Publisher:   "Xbox Game Studios",
			Developer:   "Mojang Studios",
			ReleaseDate: "2022-06-07",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1672970/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1672970/header.jpg",
			Screenshots: []string{},
			BasePrice:   29.99,
			BestPrice:   13.79,
			Rating:      4.92,
			ReviewCount: 8910,
			LikesCount:  37500,
			SalesCount:  108900,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11", Processor: "Intel Core i5-4690 / AMD A10-7800", Memory: "8 GB RAM", Graphics: "GeForce 700 Series or AMD Radeon Rx 200", Storage: "4 GB SSD"},
			Steps:       []string{"Visit redeem.microsoft.com or open the Xbox App on Windows", "Sign in with your Microsoft account", "Enter your 25-character code"},
		},
		{
			Slug:        "forza-horizon-5-premium-edition-xbox-pc-key-global",
			Title:       "Forza Horizon 5: Premium Edition",
			Subtitle:    "Xbox Play Anywhere (PC + Xbox Series X|S) + Hot Wheels & Rally",
			Description: "Your Ultimate Horizon Adventure awaits! Explore the vibrant and ever-evolving open world landscapes of Mexico with limitless, fun driving action in hundreds of the world's greatest cars.",
			Platform:    "Xbox Live",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Sports",
			Publisher:   "Xbox Game Studios",
			Developer:   "Playground Games",
			ReleaseDate: "2021-11-09",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1551360/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1551360/header.jpg",
			Screenshots: []string{},
			BasePrice:   99.99,
			BestPrice:   36.40,
			Rating:      4.89,
			ReviewCount: 3420,
			LikesCount:  16200,
			SalesCount:  42300,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11", Processor: "Intel i5-8400 or AMD Ryzen 5 1500X", Memory: "16 GB RAM", Graphics: "NVidia GTX 1070 OR AMD RX 590", Storage: "110 GB SSD"},
			Steps:       []string{"Open Microsoft Store or Xbox App", "Click Redeem a Code", "Enter 25-character key"},
		},
		{
			Slug:        "hogwarts-legacy-deluxe-edition-steam-key-global",
			Title:       "Hogwarts Legacy: Digital Deluxe Edition",
			Subtitle:    "Open-World Wizarding RPG + Dark Arts Pack & Thestral Mount",
			Description: "Hogwarts Legacy is an immersive, open-world action RPG set in the world first introduced in the Harry Potter books. Embark on a journey through familiar and new locations as you explore and discover fantastic beasts, customize your character and craft potions.",
			Platform:    "Steam",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Action RPG",
			Publisher:   "Warner Bros. Games",
			Developer:   "Avalanche Software",
			ReleaseDate: "2023-02-10",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/990080/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/990080/header.jpg",
			Screenshots: []string{},
			BasePrice:   69.99,
			BestPrice:   18.49,
			Rating:      4.84,
			ReviewCount: 4950,
			LikesCount:  21800,
			SalesCount:  63400,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "Intel Core i7-8700 / AMD Ryzen 5 3600", Memory: "16 GB RAM", Graphics: "NVIDIA GeForce 1080 Ti or AMD Radeon RX 5700 XT", Storage: "85 GB SSD"},
			Steps:       []string{"Open Steam client", "Click 'Activate a Product on Steam'", "Enter key and download"},
		},
		{
			Slug:        "resident-evil-4-gold-edition-steam-key-europe",
			Title:       "Resident Evil 4 Gold Edition",
			Subtitle:    "Remake Base Game + Separate Ways DLC + Extra DLC Pack",
			Description: "Survival is just the beginning. Six years have passed since the biological disaster in Raccoon City. Leon S. Kennedy, one of the survivors, tracks the president's kidnapped daughter to a secluded European village.",
			Platform:    "Steam",
			Region:      "EUROPE",
			ProductType: "Game Key",
			Genre:       "Action",
			Publisher:   "CAPCOM Co., Ltd.",
			Developer:   "CAPCOM Co., Ltd.",
			ReleaseDate: "2024-02-09",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2050650/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2050650/header.jpg",
			Screenshots: []string{},
			BasePrice:   49.99,
			BestPrice:   22.99,
			Rating:      4.95,
			ReviewCount: 4100,
			LikesCount:  19500,
			SalesCount:  44800,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "AMD Ryzen 5 3600 / Intel Core i7 8700", Memory: "16 GB RAM", Graphics: "AMD Radeon RX 5700 / NVIDIA GeForce GTX 1070", Storage: "67 GB SSD"},
			Steps:       []string{"Open Steam client", "Select 'Activate a Product on Steam'", "Enter key"},
		},
		{
			Slug:        "xbox-game-pass-ultimate-3-months-subscription-global",
			Title:       "Xbox Game Pass Ultimate - 3 Months",
			Subtitle:    "Hundreds of Day-One Games on Console, PC & Cloud Gaming",
			Description: "Play new games on day one. Plus, enjoy hundreds of high-quality games with friends on console, PC, or cloud. Includes EA Play membership and online console multiplayer.",
			Platform:    "Xbox Live",
			Region:      "GLOBAL",
			ProductType: "Subscription",
			Genre:       "Subscription",
			Publisher:   "Microsoft Corporation",
			Developer:   "Xbox Game Studios",
			ReleaseDate: "2024-01-01",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1172470/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1172470/header.jpg",
			Screenshots: []string{},
			BasePrice:   50.99,
			BestPrice:   29.49,
			Rating:      4.91,
			ReviewCount: 6420,
			LikesCount:  28100,
			SalesCount:  79300,
			SystemReqs:  models.SystemReq{OS: "Xbox Series X|S, Xbox One, Windows 10/11", Processor: "Any Supported Device", Memory: "N/A", Graphics: "Cloud / Local", Storage: "N/A"},
			Steps:       []string{"Sign in to redeem.microsoft.com", "Paste the 25-digit subscription code", "Confirm activation on your Microsoft account"},
		},
		{
			Slug:        "playstation-network-card-50-eur-europe",
			Title:       "PlayStation Network Gift Card 50 EUR",
			Subtitle:    "Instant PSN Wallet Top-Up for PS5 & PS4 Store Europe",
			Description: "Easily top up your PlayStation Network wallet without a credit card. Purchase full games, add-ons, PlayStation Plus subscriptions, and in-game currency directly from the PS Store.",
			Platform:    "PlayStation",
			Region:      "EUROPE",
			ProductType: "Gift Card",
			Genre:       "Gift Card",
			Publisher:   "Sony Interactive Entertainment",
			Developer:   "PlayStation Network",
			ReleaseDate: "2024-01-01",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1593500/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1593500/header.jpg",
			Screenshots: []string{},
			BasePrice:   50.00,
			BestPrice:   42.99,
			Rating:      4.96,
			ReviewCount: 8300,
			LikesCount:  31000,
			SalesCount:  96500,
			SystemReqs:  models.SystemReq{OS: "PlayStation 5 / PlayStation 4 / Web Browser", Processor: "European PSN Account", Memory: "N/A", Graphics: "N/A", Storage: "N/A"},
			Steps:       []string{"Go to PlayStation Store on your console or web", "Select 'Redeem Codes' from the profile menu", "Enter the 12-digit PSN voucher code"},
		},
		{
			Slug:        "steam-gift-card-25-usd-united-states",
			Title:       "Steam Wallet Gift Card 25 USD",
			Subtitle:    "Instant Balance Top-Up for US Steam Accounts",
			Description: "Steam Gift Cards work just like a gift certificate, and can be redeemed on Steam for the purchase of games, software, hardware, and any other item you can purchase on Steam.",
			Platform:    "Steam",
			Region:      "UNITED STATES",
			ProductType: "Gift Card",
			Genre:       "Gift Card",
			Publisher:   "Valve Corporation",
			Developer:   "Valve Corporation",
			ReleaseDate: "2024-01-01",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/730/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/730/header.jpg",
			Screenshots: []string{},
			BasePrice:   25.00,
			BestPrice:   21.79,
			Rating:      4.93,
			ReviewCount: 5120,
			LikesCount:  19800,
			SalesCount:  71200,
			SystemReqs:  models.SystemReq{OS: "Windows / macOS / Linux / Steam Deck", Processor: "US Currency Steam Wallet", Memory: "N/A", Graphics: "N/A", Storage: "N/A"},
			Steps:       []string{"Open store.steampowered.com/account/redeemwalletcode", "Enter your 15-character Steam Wallet code", "Click Continue to credit your balance"},
		},
		{
			Slug:        "the-legend-of-zelda-tears-of-the-kingdom-nintendo-switch-key-global",
			Title:       "The Legend of Zelda: Tears of the Kingdom",
			Subtitle:    "Epic Sky & Surface Adventure on Nintendo Switch eShop",
			Description: "An epic adventure across the land and skies of Hyrule awaits in The Legend of Zelda: Tears of the Kingdom for Nintendo Switch. Craft weapons, vehicles, and contraptions with Ultrahand and Fuse.",
			Platform:    "Nintendo eShop",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Adventure",
			Publisher:   "Nintendo",
			Developer:   "Nintendo EPD",
			ReleaseDate: "2023-05-12",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1203220/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1203220/header.jpg",
			Screenshots: []string{},
			BasePrice:   69.99,
			BestPrice:   45.99,
			Rating:      4.96,
			ReviewCount: 4310,
			LikesCount:  23400,
			SalesCount:  48700,
			SystemReqs:  models.SystemReq{OS: "Nintendo Switch / Switch OLED / Switch Lite", Processor: "Nintendo Custom Tegra", Memory: "4 GB", Graphics: "Docked 1080p / Handheld 720p", Storage: "16.3 GB"},
			Steps:       []string{"Open Nintendo eShop on your Switch console", "Select 'Enter Code' on the left sidebar", "Input the 16-character Nintendo download code"},
		},
		{
			Slug:        "civilization-vii-founders-edition-steam-key-global",
			Title:       "Sid Meier's Civilization VII",
			Subtitle:    "Build the Greatest Empire the World Has Ever Known",
			Description: "The award-winning strategy game franchise returns with a revolutionary new chapter. Sid Meier's Civilization VII empowers you to build the greatest empire the world has ever known across evolving human Ages.",
			Platform:    "Steam",
			Region:      "GLOBAL",
			ProductType: "Game Key",
			Genre:       "Strategy",
			Publisher:   "2K",
			Developer:   "Firaxis Games",
			ReleaseDate: "2025-02-11",
			CoverImage:  "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1295660/library_600x900_2x.jpg",
			BannerImage: "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1295660/header.jpg",
			Screenshots: []string{},
			BasePrice:   69.99,
			BestPrice:   47.90,
			Rating:      4.76,
			ReviewCount: 1920,
			LikesCount:  10800,
			SalesCount:  26400,
			SystemReqs:  models.SystemReq{OS: "Windows 10/11 64-bit", Processor: "Intel Core i5-10400 / AMD Ryzen 5 3600X", Memory: "16 GB RAM", Graphics: "NVIDIA RTX 2060 / AMD RX 6600", Storage: "20 GB SSD"},
			Steps:       []string{"Open Steam client", "Select 'Activate a Product on Steam'", "Enter key and start your empire"},
		},
	}

	for idx, cp := range catalog {
		screenshotsBytes, _ := json.Marshal(cp.Screenshots)
		sysReqsBytes, _ := json.Marshal(cp.SystemReqs)
		stepsBytes, _ := json.Marshal(cp.Steps)

		res, err := tx.Exec(`
			INSERT INTO products (
				slug, title, subtitle, description, platform, region, product_type, genre,
				publisher, developer, release_date, cover_image, banner_image,
				screenshots_json, base_price, rating, review_count, likes_count, sales_count,
				system_reqs_json, activation_steps_json
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			cp.Slug, cp.Title, cp.Subtitle, cp.Description, cp.Platform, cp.Region, cp.ProductType, cp.Genre,
			cp.Publisher, cp.Developer, cp.ReleaseDate, cp.CoverImage, cp.BannerImage,
			string(screenshotsBytes), cp.BasePrice, cp.Rating, cp.ReviewCount, cp.LikesCount, cp.SalesCount,
			string(sysReqsBytes), string(stepsBytes),
		)
		if err != nil {
			return err
		}
		productID, _ := res.LastInsertId()

		// Create 4 competing offers from different verified sellers for every product!
		priceSpread := []float64{0.00, 1.40, 2.85, 4.50}
		cashbackSpread := []float64{7.5, 5.0, 6.0, 4.0}

		for offerIdx := 0; offerIdx < 4; offerIdx++ {
			sellerID := sellerIDs[(idx+offerIdx)%len(sellerIDs)]
			offerPrice := math.Round((cp.BestPrice+priceSpread[offerIdx])*100) / 100
			cashbackPct := cashbackSpread[offerIdx]
			isFeatured := 0
			if offerIdx == 0 {
				isFeatured = 1
			}

			offerRes, err := tx.Exec(`
				INSERT INTO offers (product_id, seller_id, price, cashback_percent, delivery_time, is_featured)
				VALUES (?, ?, ?, ?, 'Instant', ?)
			`, productID, sellerID, offerPrice, cashbackPct, isFeatured)
			if err != nil {
				return err
			}
			offerID, _ := offerRes.LastInsertId()

			// Seed 4 encrypted digital activation keys per offer using AES-256-GCM!
			for k := 0; k < 4; k++ {
				plainKey := generateFormattedKey(cp.Platform, cp.Slug)
				encryptedKey, err := crypto.EncryptSecret(plainKey, s.KeyEncryptionSecret)
				if err != nil {
					return fmt.Errorf("failed encrypting seed key: %w", err)
				}
				keyFmt := fmt.Sprintf("%s Official Activation Key (%s)", cp.Platform, cp.Region)
				if _, err := tx.Exec(`
					INSERT INTO digital_keys (offer_id, encrypted_key, key_format, is_sold)
					VALUES (?, ?, ?, 0)
				`, offerID, encryptedKey, keyFmt); err != nil {
					return err
				}
			}
		}
	}

	// 3. Seed Dynamic Homepage Blocks (Nothing hardcoded in frontend!)
	type seedBlock struct {
		Slug        string
		BlockType   string
		Title       string
		Subtitle    string
		AccentColor string
		BadgeText   string
		CTAText     string
		CTALink     string
		SortOrder   int
		ServiceName string
		Config      map[string]interface{}
		CustomItems []models.BlockCustomItem
	}

	blocks := []seedBlock{
		{
			Slug:        "hero-spotlight-banner",
			BlockType:   "hero_banner",
			Title:       "MEGA GAMING WEEKEND SALE",
			Subtitle:    "Save up to 80% on verified instant-delivery Steam, PlayStation, Xbox & Nintendo keys plus earn up to 7.5% Eneba Cashback on every order.",
			AccentColor: "#ffd524",
			BadgeText:   "UP TO 80% OFF + INSTANT CASHBACK",
			CTAText:     "Explore All Deals",
			CTALink:     "/store?sort=discount_desc",
			SortOrder:   1,
			ServiceName: "HeroPromotionService",
			Config: map[string]interface{}{
				"limit": 3,
				"featured_slugs": []string{
					"cyberpunk-2077-ultimate-edition-steam-key-global",
					"elden-ring-shadow-of-the-erdtree-edition-steam-key-global",
					"black-myth-wukong-deluxe-edition-steam-key-global",
				},
				"promo_code_hint": "INTERN2027",
				"promo_discount":  "20% EXTRA OFF",
			},
			CustomItems: []models.BlockCustomItem{},
		},
		{
			Slug:        "flash-deals-block",
			BlockType:   "flash_deals",
			Title:       "Flash Deals & Biggest Discounts",
			Subtitle:    "Curated dynamically by FlashDealsService — highest verified discounts vs official retail MSRP",
			AccentColor: "#ff4081",
			BadgeText:   "HOT DEALS",
			CTAText:     "View All Discounts",
			CTALink:     "/store?sort=discount_desc",
			SortOrder:   2,
			ServiceName: "FlashDealsService",
			Config: map[string]interface{}{
				"limit":            5,
				"countdown_hours":  14,
				"show_discount_tag": true,
			},
			CustomItems: []models.BlockCustomItem{},
		},
		{
			Slug:        "best-sellers-block",
			BlockType:   "product_grid",
			Title:       "Top Selling Games Right Now",
			Subtitle:    "Powered by BestSellersService — live ranking by completed key activations across global regions",
			AccentColor: "#ffd524",
			BadgeText:   "BEST SELLERS",
			CTAText:     "Browse Top Sellers",
			CTALink:     "/store?sort=popular",
			SortOrder:   3,
			ServiceName: "BestSellersService",
			Config: map[string]interface{}{
				"limit": 5,
			},
			CustomItems: []models.BlockCustomItem{},
		},
		{
			Slug:        "platform-categories-block",
			BlockType:   "category_grid",
			Title:       "Browse by Gaming Platform",
			Subtitle:    "Delivered by CategoryNavigationService — jump straight to your preferred ecosystem",
			AccentColor: "#00e676",
			BadgeText:   "ECOSYSTEMS",
			CTAText:     "Full Store Catalog",
			CTALink:     "/store",
			SortOrder:   4,
			ServiceName: "CategoryNavigationService",
			Config:      map[string]interface{}{},
			CustomItems: []models.BlockCustomItem{
				{Title: "Steam Keys", Subtitle: "PC & Steam Deck Instant Codes", Badge: "MOST POPULAR", Icon: "monitor", Link: "/store?platform=Steam", AccentColor: "#00b0ff", Metric: "Global & Europe"},
				{Title: "Xbox Live & Game Pass", Subtitle: "Series X|S, Xbox One & Windows", Badge: "PLAY ANYWHERE", Icon: "gamepad-2", Link: "/store?platform=Xbox+Live", AccentColor: "#00e676", Metric: "Up to 65% Off"},
				{Title: "PlayStation Network", Subtitle: "PS5 & PS4 Games & PSN Cards", Badge: "INSTANT TOP-UP", Icon: "trophy", Link: "/store?platform=PlayStation", AccentColor: "#536dfe", Metric: "Instant Voucher"},
				{Title: "Nintendo eShop", Subtitle: "Switch Exclusives & Memberships", Badge: "FAMILY FAVORITE", Icon: "sparkles", Link: "/store?platform=Nintendo+eShop", AccentColor: "#ff5252", Metric: "Global Keys"},
				{Title: "GOG & Rockstar", Subtitle: "DRM-Free Classics & Outlaw Epics", Badge: "DRM-FREE", Icon: "shield-check", Link: "/store?platform=GOG.com", AccentColor: "#ffd524", Metric: "Up to 80% Off"},
				{Title: "EA App Sports", Subtitle: "EA SPORTS FC 25 & Blockbusters", Badge: "SPORTS HUB", Icon: "zap", Link: "/store?platform=EA+App", AccentColor: "#ff9100", Metric: "Official Partner"},
			},
		},
		{
			Slug:        "trending-community-favorites",
			BlockType:   "product_grid",
			Title:       "Trending & Highest Rated by Gamers",
			Subtitle:    "Powered by TrendingNowService — ranked by community likes and verified buyer review scores",
			AccentColor: "#00b0ff",
			BadgeText:   "COMMUNITY PICKS",
			CTAText:     "See Highest Rated",
			CTALink:     "/store?sort=rating_desc",
			SortOrder:   5,
			ServiceName: "TrendingNowService",
			Config: map[string]interface{}{
				"limit": 5,
			},
			CustomItems: []models.BlockCustomItem{},
		},
		{
			Slug:        "budget-under-25-block",
			BlockType:   "product_grid",
			Title:       "Legendary Games Under €25",
			Subtitle:    "Filtered dynamically by BudgetDealsService — AAA masterpieces at indie prices",
			AccentColor: "#00e676",
			BadgeText:   "UNDER €25",
			CTAText:     "Shop Under €25",
			CTALink:     "/store?max_price=25",
			SortOrder:   6,
			ServiceName: "BudgetDealsService",
			Config: map[string]interface{}{
				"limit":     5,
				"max_price": 25.0,
			},
			CustomItems: []models.BlockCustomItem{},
		},
		{
			Slug:        "gift-cards-and-subscriptions-block",
			BlockType:   "product_grid",
			Title:       "Top-Up Gift Cards & Gaming Subscriptions",
			Subtitle:    "Served by GiftCardsAndSubsService — instant wallet funds for PSN, Steam & Xbox Game Pass",
			AccentColor: "#aa00ff",
			BadgeText:   "GIFT CARDS & SUBS",
			CTAText:     "View All Gift Cards",
			CTALink:     "/store?product_type=Gift+Card",
			SortOrder:   7,
			ServiceName: "GiftCardsAndSubsService",
			Config: map[string]interface{}{
				"limit": 4,
			},
			CustomItems: []models.BlockCustomItem{},
		},
		{
			Slug:        "marketplace-trust-guarantees",
			BlockType:   "value_props",
			Title:       "Why Over 15 Million Gamers Trust Our Marketplace",
			Subtitle:    "Served by TrustAndSecurityService — end-to-end encrypted key vault & verified merchant protection",
			AccentColor: "#ffd524",
			BadgeText:   "ENEBA SHIELD",
			CTAText:     "Create Free Account",
			CTALink:     "/register",
			SortOrder:   8,
			ServiceName: "TrustAndSecurityService",
			Config:      map[string]interface{}{},
			CustomItems: []models.BlockCustomItem{
				{Title: "AES-256-GCM Key Vault", Subtitle: "Every digital activation key is encrypted at rest and decrypted only upon verified checkout.", Badge: "ZERO PLAINTEXT", Icon: "lock", AccentColor: "#ffd524"},
				{Title: "Multi-Seller Price Competition", Subtitle: "Compare offers from verified merchants side-by-side to lock in the lowest price or highest cashback.", Badge: "BEST PRICE", Icon: "users", AccentColor: "#00e676"},
				{Title: "Instant Digital Delivery", Subtitle: "Keys are dispatched immediately to your screen and order history the second payment completes.", Badge: "< 2 SECONDS", Icon: "zap", AccentColor: "#00b0ff"},
				{Title: "Eneba Purple Cashback", Subtitle: "Earn up to 7.5% cashback directly into your Eneba Wallet on eligible marketplace offers.", Badge: "CASHBACK REWARDS", Icon: "coins", AccentColor: "#ff4081"},
			},
		},
	}

	for _, b := range blocks {
		cfgJSON, _ := json.Marshal(b.Config)
		itemsJSON, _ := json.Marshal(b.CustomItems)
		if _, err := tx.Exec(`
			INSERT INTO homepage_blocks (
				slug, block_type, title, subtitle, accent_color, badge_text,
				cta_text, cta_link, sort_order, service_name, query_config_json, custom_items_json, is_active
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)
		`,
			b.Slug, b.BlockType, b.Title, b.Subtitle, b.AccentColor, b.BadgeText,
			b.CTAText, b.CTALink, b.SortOrder, b.ServiceName, string(cfgJSON), string(itemsJSON),
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}
