package search

import (
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/eneba-intern/marketplace/backend/internal/models"
)

// Result wraps a matched Product with its relevance score.
type Result struct {
	Product models.Product
	Score   float64
}

// RankProducts filters and ranks products against a user search query using a
// hybrid scoring algorithm:
// 1. Exact & substring phrase matching
// 2. Token prefix matching (partial word support)
// 3. Damerau-Levenshtein edit distance per word (typo tolerance for insertions, deletions, substitutions, transpositions)
// 4. Character trigram Jaccard similarity
func RankProducts(products []models.Product, rawQuery string) ([]models.Product, string) {
	query := Normalize(rawQuery)
	if query == "" {
		return products, ""
	}

	queryTokens := strings.Fields(query)
	queryTrigrams := Trigrams(query)

	var scored []Result
	var bestCorrection string
	var bestCorrectionScore float64

	for _, p := range products {
		score, matchedTitle := ScoreProduct(p, query, queryTokens, queryTrigrams)
		if score >= 0.26 {
			pCopy := p
			pCopy.MatchScore = math.Round(score*1000) / 1000
			if matchedTitle != "" && !strings.Contains(Normalize(p.Title), query) {
				pCopy.MatchedCorrection = p.Title
			}
			scored = append(scored, Result{
				Product: pCopy,
				Score:   score,
			})
			if score > bestCorrectionScore {
				bestCorrectionScore = score
				if !strings.Contains(Normalize(p.Title), query) {
					bestCorrection = p.Title
				}
			}
		}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if math.Abs(scored[i].Score-scored[j].Score) > 0.02 {
			return scored[i].Score > scored[j].Score
		}
		return scored[i].Product.SalesCount > scored[j].Product.SalesCount
	})

	out := make([]models.Product, 0, len(scored))
	for _, r := range scored {
		out = append(out, r.Product)
	}

	return out, bestCorrection
}

// ScoreProduct computes a [0.0, 1.0] relevance score for a single product.
func ScoreProduct(p models.Product, query string, queryTokens []string, queryTrigrams map[string]struct{}) (float64, string) {
	normTitle := Normalize(p.Title)
	normSub := Normalize(p.Subtitle)
	normPlatform := Normalize(p.Platform)
	normGenre := Normalize(p.Genre)

	// 1. Exact title match
	if normTitle == query {
		return 1.0, p.Title
	}

	// 2. Direct prefix or substring match on Title
	if strings.HasPrefix(normTitle, query) {
		return 0.96, p.Title
	}
	if strings.Contains(normTitle, query) {
		return 0.90, p.Title
	}

	// 3. Initialism / Acronym match (e.g. "gta", "rdr", "bg3")
	acronym := buildAcronym(normTitle)
	if acronym != "" && (acronym == query || strings.HasPrefix(acronym, query)) {
		return 0.88, p.Title
	}

	titleTokens := strings.Fields(normTitle + " " + normPlatform + " " + normGenre + " " + normSub)
	primaryTitleTokens := strings.Fields(normTitle)

	// 4. Token-level fuzzy + partial matching
	if len(queryTokens) > 0 {
		totalTokenScore := 0.0
		matchedTokens := 0

		for _, qt := range queryTokens {
			bestWordScore := 0.0

			for _, tt := range titleTokens {
				isPrimary := containsToken(primaryTitleTokens, tt)
				weight := 1.0
				if !isPrimary {
					weight = 0.75
				}

				if tt == qt {
					if 1.0*weight > bestWordScore {
						bestWordScore = 1.0 * weight
					}
					continue
				}

				// Partial prefix match (e.g., "cyber" -> "cyberpunk", "witch" -> "witcher")
				if len(qt) >= 2 && strings.HasPrefix(tt, qt) {
					prefixRatio := float64(len(qt)) / float64(len(tt))
					s := (0.78 + 0.20*prefixRatio) * weight
					if s > bestWordScore {
						bestWordScore = s
					}
					continue
				}

				// Substring inside a token (e.g., "punk" -> "cyberpunk")
				if len(qt) >= 3 && strings.Contains(tt, qt) {
					s := 0.72 * weight
					if s > bestWordScore {
						bestWordScore = s
					}
					continue
				}

				// Damerau-Levenshtein typo tolerance on words
				if len(qt) >= 3 && len(tt) >= 3 {
					dist := DamerauLevenshtein(qt, tt)
					maxLen := maxInt(len(qt), len(tt))
					allowedDist := 1
					if maxLen >= 6 {
						allowedDist = 2
					}
					if maxLen >= 10 {
						allowedDist = 3
					}

					// Also check edit distance against the prefix of a longer title word
					// (e.g. "cyberpnk" vs "cyberpunk", or "helldvr" vs "helldivers")
					if len(tt) > len(qt) {
						prefixCandidate := tt[:minInt(len(tt), len(qt)+1)]
						prefixDist := DamerauLevenshtein(qt, prefixCandidate)
						if prefixDist < dist {
							dist = prefixDist
						}
					}

					if dist <= allowedDist {
						sim := 1.0 - (float64(dist) / float64(maxLen))
						s := (0.45 + 0.45*sim) * weight
						if s > bestWordScore {
							bestWordScore = s
						}
					}
				}
			}

			if bestWordScore >= 0.42 {
				matchedTokens++
			}
			totalTokenScore += bestWordScore
		}

		avgTokenScore := totalTokenScore / float64(len(queryTokens))
		coverage := float64(matchedTokens) / float64(len(queryTokens))

		// Combine token fuzzy score with trigram similarity on the title
		trigramSim := TrigramSimilarity(queryTrigrams, Trigrams(normTitle))
		combined := avgTokenScore*0.78 + trigramSim*0.22

		if coverage >= 0.5 && combined >= 0.30 {
			return math.Min(0.89, combined), p.Title
		}
	}

	// 5. Fallback whole-title trigram similarity for fused words (e.g. "eldenring", "reddead")
	trigramSim := TrigramSimilarity(queryTrigrams, Trigrams(normTitle))
	if trigramSim >= 0.32 {
		return trigramSim * 0.85, p.Title
	}

	return 0.0, ""
}

// Normalize lowercases text, strips punctuation, and collapses spaces.
func Normalize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	lastSpace := true
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastSpace = false
		} else if !lastSpace {
			b.WriteByte(' ')
			lastSpace = true
		}
	}
	return strings.TrimSpace(b.String())
}

// Trigrams generates a set of 3-character sliding windows from a padded string.
func Trigrams(s string) map[string]struct{} {
	clean := strings.ReplaceAll(s, " ", "")
	res := make(map[string]struct{})
	if len(clean) == 0 {
		return res
	}
	padded := "  " + clean + " "
	runes := []rune(padded)
	for i := 0; i <= len(runes)-3; i++ {
		tri := string(runes[i : i+3])
		res[tri] = struct{}{}
	}
	return res
}

// TrigramSimilarity returns the Jaccard similarity coefficient between two trigram sets.
func TrigramSimilarity(a, b map[string]struct{}) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	intersection := 0
	for k := range a {
		if _, ok := b[k]; ok {
			intersection++
		}
	}
	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

// DamerauLevenshtein calculates the optimal string alignment distance between two strings
// (supporting insertions, deletions, substitutions, and adjacent transpositions).
func DamerauLevenshtein(a, b string) int {
	ra := []rune(a)
	rb := []rune(b)
	la, lb := len(ra), len(rb)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
		d[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if ra[i-1] != rb[j-1] {
				cost = 1
			}
			d[i][j] = minInt(
				d[i-1][j]+1,      // deletion
				minInt(
					d[i][j-1]+1,  // insertion
					d[i-1][j-1]+cost, // substitution
				),
			)
			if i > 1 && j > 1 && ra[i-1] == rb[j-2] && ra[i-2] == rb[j-1] {
				d[i][j] = minInt(d[i][j], d[i-2][j-2]+cost) // transposition
			}
		}
	}
	return d[la][lb]
}

func buildAcronym(normTitle string) string {
	words := strings.Fields(normTitle)
	if len(words) < 2 {
		return ""
	}
	var b strings.Builder
	for _, w := range words {
		if len(w) > 0 {
			// Include digits if whole word is a number (e.g. "gta 5" -> "gta5")
			if unicode.IsDigit(rune(w[0])) {
				b.WriteString(w)
			} else {
				b.WriteByte(w[0])
			}
		}
	}
	return b.String()
}

func containsToken(tokens []string, target string) bool {
	for _, t := range tokens {
		if t == target {
			return true
		}
	}
	return false
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
