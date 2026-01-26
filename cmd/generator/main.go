package main

import (
	"bufio"
	"fmt"
	"galaxies/internal/core"
	"galaxies/internal/core/enums"
	"galaxies/internal/gen"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Global state to track where the player is currently docked
var currentSystem *core.System
var lastDiscoveredSystem *core.System // Track the last 'gen' result to allow jumping

func main() {
	// 1. Init
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	// 2. Establish "Home Base" (Coordinates 0, 0)
	currentSystem = gen.GenerateSystem(gen.GenerateSystemConfig{
		Name:      "Home Station Alpha", // Force a name
		X:         0,
		Y:         0,
		Political: enums.PoliticalStatus(-1),
		Economic:  enums.EconomicStatus(-1),
		Social:    enums.SocialStatus(-1),
	})

	fmt.Println("================================================")
	fmt.Println("   GALACTIC SURVEY DATABASE :: CORE WORLDS SCAN")
	fmt.Println("================================================")
	fmt.Printf("STATUS: DOCKED AT %s [%d, %d]\n", currentSystem.Name, currentSystem.X, currentSystem.Y)
	fmt.Println("")

	// 3. Generate initial batch for flavor
	generateBatch(5)

	fmt.Println("")
	fmt.Println("SYSTEM READY. WAITING FOR INPUT.")
	fmt.Println("Commands: [gen] scan new sector | [jump] travel to scan | [exit]")
	fmt.Println("------------------------------------------------")

	// 4. Interactive Command Loop
	for {
		fmt.Print("COMMAND > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "gen", "generate":
			// Generate a random target in the Outer Rim (-20 to 20 range)
			x := rand.Intn(41) - 20
			y := rand.Intn(41) - 20

			config := gen.GenerateSystemConfig{
				X:         x,
				Y:         y,
				Political: enums.PoliticalStatus(-1),
				Economic:  enums.EconomicStatus(-1),
				Social:    enums.SocialStatus(-1),
			}

			sys := gen.GenerateSystem(config)
			lastDiscoveredSystem = sys // Save it so we can jump to it

			fmt.Println("\n>> NEW SIGNAL DETECTED...")
			PrintSystemDossier(sys)
			PrintNavigationData(sys)

		case "jump", "travel":
			if lastDiscoveredSystem == nil {
				fmt.Println("ERROR: No navigation target set. Run 'gen' first.")
				continue
			}
			
			// Simple "travel" logic
			currentSystem = lastDiscoveredSystem
			fmt.Printf("\n>> JUMP DRIVE ENGAGED... TRAVEL COMPLETE.\n")
			fmt.Printf(">> WELCOME TO: %s\n", currentSystem.Name)
			// Reset target
			lastDiscoveredSystem = nil

		case "exit", "quit":
			fmt.Println("Shutting down survey link...")
			return

		default:
			if input != "" {
				fmt.Println("Unknown command. Try 'gen', 'jump', or 'exit'.")
			}
		}
	}
}

// Helper to calculate and print navigation physics
func PrintNavigationData(target *core.System) {
	// Use the DefaultShip we defined in core/navigation.go
	nav := core.CalculateJump(currentSystem.X, currentSystem.Y, target.X, target.Y, core.DefaultShip)

	fmt.Println(">> NAVIGATION COMPUTER CALCULATION")
	fmt.Printf("   Origin:       %s [%d, %d]\n", currentSystem.Name, currentSystem.X, currentSystem.Y)
	fmt.Printf("   Target:       %s [%d, %d]\n", target.Name, target.X, target.Y)
	fmt.Printf("   Distance:     %.2f Light Years\n", nav.Distance)
	fmt.Printf("   Fuel Need:    %.1f Units\n", nav.FuelRequired)
	fmt.Printf("   Est. Time:    %.1f Hours\n", nav.TimeHours)

	// Calculate Cost based on CURRENT system's fuel prices
	// Base Price (10cr) * Local Multiplier * Fuel Needed
	baseFuelPrice := 10.0
	
	if currentSystem.Stats.HasRefueling {
		cost := nav.FuelRequired * baseFuelPrice * currentSystem.Stats.FuelCostMult
		fmt.Printf("   Refuel Cost:  %.0f Credits (At current station prices)\n", cost)
	} else {
		fmt.Println("   Refuel Cost:  WARNING - NO FUEL AVAILABLE AT CURRENT STATION")
	}

	if !nav.IsReachable {
		fmt.Println("   [!] WARNING: TARGET EXCEEDS MAX FUEL RANGE")
	}
	fmt.Println("")
}

func generateBatch(count int) {
	for i := 0; i < count; i++ {
		x := rand.Intn(21) - 10
		y := rand.Intn(21) - 10
		config := gen.GenerateSystemConfig{
			X: x, Y: y,
			Political: enums.PoliticalStatus(-1),
			Economic:  enums.EconomicStatus(-1),
			Social:    enums.SocialStatus(-1),
		}
		sys := gen.GenerateSystem(config)
		PrintSystemDossier(sys)
	}
}

// Reuse your existing PrintSystemDossier function here...
func PrintSystemDossier(sys *core.System) {
	fmt.Printf("SYSTEM: %-20s COORDS: [%d, %d]\n", sys.Name, sys.X, sys.Y)
	fmt.Println("------------------------------------------------")
	fmt.Printf("POL: %-22s ECO: %-22s SOC: %s\n", sys.Political, sys.Economic, sys.Social)
	fmt.Println("------------------------------------------------")

	fmt.Printf("MARKET  >> Tax: %.0f%% | Docking: %dcr | Buy: x%.2f | Sell: x%.2f\n",
		sys.Stats.TaxRate*100, sys.Stats.DockingFee, sys.Stats.MarketBuyMult, sys.Stats.MarketSellMult)

	if sys.Stats.HasBlackMarket {
		fmt.Printf("ILLEGAL >> [OPEN] Fence Buys: x%.2f | Smuggling Profit: x%.2f\n",
			sys.Stats.BlackMarketBuyMult, sys.Stats.ContrabandProfit)
	}

	fmt.Printf("RISK    >> Piracy: %.0f%% | Inspection: %.0f%% | Bribe Cost: x%.2f\n",
		sys.Stats.PiracyChance*100, sys.Stats.InspectionChance*100, sys.Stats.BribeCostMult)

	fmt.Printf("POPL    >> Pax: %d (Wealth x%.1f) | Crew Pool: %d (Skill %d)\n",
		sys.Stats.PassengerCount, sys.Stats.PassengerWealth, sys.Stats.CrewPoolCount, sys.Stats.CrewSkillAvg)

	if sys.Stats.VIPCount > 0 {
		fmt.Printf("SPECIAL >> VIPs: %d (Paying x%.1f)\n", sys.Stats.VIPCount, sys.Stats.VIPWealth)
	}
	if sys.Stats.SlumsCount > 0 {
		fmt.Printf("SPECIAL >> Slum Dwellers: %d\n", sys.Stats.SlumsCount)
	}
	if sys.Stats.PrisonerCount > 0 {
		fmt.Printf("SPECIAL >> Prisoners: %d (Bail: %dcr)\n", sys.Stats.PrisonerCount, sys.Stats.PrisonerCost)
	}
	if sys.Stats.AndroidCount > 0 {
		fmt.Printf("SPECIAL >> Androids: %d (Cost: %dcr)\n", sys.Stats.AndroidCount, sys.Stats.AndroidCost)
	}

	fmt.Print("FACILS  >> ")
	if sys.Stats.HasRefueling { fmt.Print("[FUEL] ") }
	if sys.Stats.HasShipyard { fmt.Print("[SHIPYARD] ") }
	if sys.Stats.HasOutfitter { fmt.Print("[OUTFITTER] ") }
	if sys.Stats.HasCantina { fmt.Print("[CANTINA] ") }
	if sys.Stats.HasHospital { fmt.Print("[HOSPITAL] ") }
	if sys.Stats.HasMissionBoard { fmt.Print("[MISSIONS] ") }
	if sys.Stats.HasAndroidFoundry { fmt.Print("[FOUNDRY] ") }
	if sys.Stats.HasPrison { fmt.Print("[PRISON] ") }
	fmt.Println("\n")
}