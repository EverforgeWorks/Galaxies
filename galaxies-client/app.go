package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// App struct manages the application lifecycle and API context.
type App struct {
	ctx context.Context
	// BaseURL points to your backend VPS or Cloudflare Tunnel.
	BaseURL string
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{
		BaseURL: "https://api.playburnrate.com/api",
	}
}

// startup is called when the app starts. The context is saved
// so we can call runtime methods during the app's life.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// -----------------------------------------------------------------------------
// SHIP & NAVIGATION METHODS
// -----------------------------------------------------------------------------

// GetShipState fetches the current ship status, including fuel, credits, and location.
func (a *App) GetShipState() (interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/ship", a.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetPlanets fetches the static universe definition (names, coordinates).
func (a *App) GetPlanets() (interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/planets", a.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// Travel sends a POST request to move the ship to a target destination.
func (a *App) Travel(destKey string) (interface{}, error) {
	payload, _ := json.Marshal(map[string]string{"destination_key": destKey})
	resp, err := http.Post(fmt.Sprintf("%s/travel", a.BaseURL), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("travel failed: server returned %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetTravelQuote asks the server for the cost of a trip without moving.
func (a *App) GetTravelQuote(destKey string) (interface{}, error) {
	payload, _ := json.Marshal(map[string]string{"destination_key": destKey})
	resp, err := http.Post(fmt.Sprintf("%s/travel/quote", a.BaseURL), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("quote failed: server returned %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// Refuel attempts to top off the fuel tank at the current station.
func (a *App) Refuel() (interface{}, error) {
	resp, err := http.Post(fmt.Sprintf("%s/refuel", a.BaseURL), "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refuel failed: server returned %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// -----------------------------------------------------------------------------
// ECONOMY & CONTRACT METHODS
// -----------------------------------------------------------------------------

// GetAvailableContracts fetches the job board for the current planet.
func (a *App) GetAvailableContracts() (interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/contracts", a.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// AcceptJob accepts a specific contract ID and adds it to the ship's manifest.
func (a *App) AcceptJob(jobID string) (interface{}, error) {
	payload, _ := json.Marshal(map[string]string{"contract_id": jobID})
	resp, err := http.Post(fmt.Sprintf("%s/contracts/accept", a.BaseURL), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("accept job failed: server returned %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (a *App) DropJob(jobID string) (interface{}, error) {
	payload, _ := json.Marshal(map[string]string{"contract_id": jobID})
	resp, err := http.Post(fmt.Sprintf("%s/contracts/drop", a.BaseURL), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("network error dropping job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("drop failed (Status %d)", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode drop response: %w", err)
	}
	return result, nil
}

// -----------------------------------------------------------------------------
// MODULE & UPGRADE METHODS (NEW)
// -----------------------------------------------------------------------------

// GetModules fetches the list of purchasable upgrades.
// Note: The backend logic typically returns an empty list if not at 'Prime'.
func (a *App) GetModules() (interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/modules", a.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// BuyModule attempts to purchase a ship upgrade.
func (a *App) BuyModule(key string) (interface{}, error) {
	payload, _ := json.Marshal(map[string]string{"module_key": key})
	resp, err := http.Post(fmt.Sprintf("%s/modules/buy", a.BaseURL), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("purchase failed: server returned %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
