package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"

	"livescore/internal/db"
	"livescore/internal/models"

	"github.com/gofiber/fiber/v2"
)

type sportmonksLeagueResp struct {
	Data []models.League `json:"data"`
}

// FetchLeaguesFromSportmonks fetches leagues from Sportmonks and saves them to the DB
func FetchLeaguesFromSportmonks(c *fiber.Ctx) error {
	apiToken := os.Getenv("SPORTMONKS_API_TOKEN")
	if apiToken == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "SPORTMONKS_API_TOKEN not set")
	}
	url := fmt.Sprintf("https://api.sportmonks.com/v3/football/leagues?api_token=%s", apiToken)
	resp, err := http.Get(url)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "Failed to fetch from Sportmonks")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("Sportmonks error: %s", string(body)))
	}
	var leaguesResp sportmonksLeagueResp
	if err := json.NewDecoder(resp.Body).Decode(&leaguesResp); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to decode Sportmonks response")
	}
	if err := db.InsertLeagues(context.Background(), leaguesResp.Data); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save leagues to DB")
	}
	return c.JSON(fiber.Map{"inserted": len(leaguesResp.Data)})
}

// GetLeaguesHandler returns leagues from the DB with support for React Admin parameters
func GetLeaguesHandler(c *fiber.Ctx) error {
	// Get all leagues from DB
	allLeagues, err := db.GetLeagues(context.Background())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch leagues from DB")
	}

	filteredLeagues := allLeagues

	// Parse filter parameter (e.g., filter={"name":"Premier League"})
	filterParam := c.Query("filter", "{}")
	if filterParam != "{}" {
		var filters map[string]interface{}
		if err := json.Unmarshal([]byte(filterParam), &filters); err == nil {
			// Apply filters (simple string matching for demo)
			var filtered []models.League
			for _, league := range filteredLeagues {
				matches := true
				for key, value := range filters {
					switch key {
					case "name":
						if league.Name != value.(string) {
							matches = false
						}
					case "country_id":
						if league.CountryID != int(value.(float64)) {
							matches = false
						}
						// Add more filters as needed
					}
				}
				if matches {
					filtered = append(filtered, league)
				}
			}
			filteredLeagues = filtered
		}
	}

	// Parse sort parameter (e.g., sort=["id","ASC"])
	sortParam := c.Query("sort", "")
	if sortParam != "" {
		var sortFields []string
		if err := json.Unmarshal([]byte(sortParam), &sortFields); err == nil && len(sortFields) >= 2 {
			field, order := sortFields[0], sortFields[1]

			// Sort leagues by field
			switch field {
			case "id":
				if order == "ASC" {
					sort.Slice(filteredLeagues, func(i, j int) bool {
						return filteredLeagues[i].ID < filteredLeagues[j].ID
					})
				} else {
					sort.Slice(filteredLeagues, func(i, j int) bool {
						return filteredLeagues[i].ID > filteredLeagues[j].ID
					})
				}
			case "name":
				if order == "ASC" {
					sort.Slice(filteredLeagues, func(i, j int) bool {
						return filteredLeagues[i].Name < filteredLeagues[j].Name
					})
				} else {
					sort.Slice(filteredLeagues, func(i, j int) bool {
						return filteredLeagues[i].Name > filteredLeagues[j].Name
					})
				}
				// Add more sort fields as needed
			}
		}
	}

	// Parse range parameter (e.g., range=[0,9])
	start, end := 0, len(filteredLeagues)-1
	rangeParam := c.Query("range", "")
	if rangeParam != "" {
		fmt.Sscanf(rangeParam, "[%d,%d]", &start, &end)

		// Validate start and end
		if start < 0 {
			start = 0
		}
		if end >= len(filteredLeagues) {
			end = len(filteredLeagues) - 1
		}
		if start > end {
			start = end
		}
	}

	// Apply pagination
	var paginatedLeagues []models.League
	if start <= end && start < len(filteredLeagues) {
		paginatedLeagues = filteredLeagues[start : end+1]
	}

	// Add Content-Range header (required by React Admin)
	c.Append("Access-Control-Expose-Headers", "Content-Range")
	c.Append("Content-Range", fmt.Sprintf("leagues %d-%d/%d", start, end, len(filteredLeagues)))

	return c.JSON(paginatedLeagues)
}
