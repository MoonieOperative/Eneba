package search_test

import (
	"testing"

	"github.com/eneba-intern/marketplace/backend/internal/models"
	"github.com/eneba-intern/marketplace/backend/internal/search"
)

func sampleCatalog() []models.Product {
	return []models.Product{
		{ID: 1, Title: "The Witcher 3: Wild Hunt - Complete Edition", Platform: "GOG.com", Genre: "Action RPG", SalesCount: 90000},
		{ID: 2, Title: "Cyberpunk 2077: Ultimate Edition", Platform: "Steam", Genre: "Action RPG", SalesCount: 54000},
		{ID: 3, Title: "Elden Ring: Shadow of the Erdtree Edition", Platform: "Steam", Genre: "Action RPG", SalesCount: 68000},
		{ID: 4, Title: "HELLDIVERS 2", Platform: "Steam", Genre: "Shooter", SalesCount: 47000},
		{ID: 5, Title: "Baldur's Gate 3", Platform: "Steam", Genre: "Action RPG", SalesCount: 51000},
		{ID: 6, Title: "Red Dead Redemption 2: Ultimate Edition", Platform: "Rockstar Games", Genre: "Adventure", SalesCount: 83000},
	}
}

func TestSearchTypoToleranceAndPartialNames(t *testing.T) {
	products := sampleCatalog()

	cases := []struct {
		query         string
		expectedTopID int64
	}{
		{query: "witcer 3", expectedTopID: 1},       // Missing 'h' typo
		{query: "wticher", expectedTopID: 1},        // Transposition typo
		{query: "cyberpnk 2077", expectedTopID: 2},  // Missing 'u' typo
		{query: "cyber", expectedTopID: 2},          // Partial word prefix
		{query: "eldn ring", expectedTopID: 3},      // Missing 'e' typo
		{query: "helldver", expectedTopID: 4},       // Missing 'i' & 's'
		{query: "baldur gat", expectedTopID: 5},     // Partial + missing 'e'
		{query: "red ded", expectedTopID: 6},        // Phonetic typo
	}

	for _, tc := range cases {
		results, _ := search.RankProducts(products, tc.query)
		if len(results) == 0 {
			t.Fatalf("expected results for query %q, got none", tc.query)
		}
		if results[0].ID != tc.expectedTopID {
			t.Errorf("for query %q: expected top product ID %d (%s), got ID %d (%s)",
				tc.query, tc.expectedTopID, "", results[0].ID, results[0].Title)
		}
	}
}
