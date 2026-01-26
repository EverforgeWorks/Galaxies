package gen

import (
	"math/rand"
	"time"

	"galaxies/internal/core/domain"
	"galaxies/internal/core/entity"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateCrew creates a new crew member.
// Logic: Checks System stats to determine if they are Android, Convict, or Standard.
func GenerateCrew(role domain.CrewRole, sysStats entity.SystemStats) *entity.Crew {
	identifier := GenerateCrewID()

	// 1. Determine Crew Type based on System Facilities
	cType := domain.CrewTypeStandard

	// Logic: If a system has special facilities, there is a chance to spawn that type.
	// Priority: Androids > Convicts > Standard
	if sysStats.HasAndroidFoundry && rand.Float64() < 0.30 { // 30% chance in Foundry systems
		cType = domain.CrewTypeAndroid
	} else if sysStats.HasPrison && rand.Float64() < 0.30 { // 30% chance in Prison systems
		cType = domain.CrewTypeConvict
	}

	// 2. Calculate Skill (1-10)
	baseSkill := rand.Intn(4) + 1
	finalSkill := baseSkill + sysStats.CrewSkillBonus

	// Type Modifiers for Skill
	if cType == domain.CrewTypeAndroid {
		finalSkill += 2 // Androids are programmed for perfection
	}
	
	// Clamp Skill
	if finalSkill < 1 { finalSkill = 1 }
	if finalSkill > 10 { finalSkill = 10 }

	// 3. Calculate Salary
	baseSalary := 100
	roleMult := 1.0

	// Role Multipliers
	switch role {
	case domain.RoleEngineer, domain.RoleNavigator:
		roleMult = 1.5
	case domain.RoleFirstOfficer:
		roleMult = 2.0
	}

	salary := int(float64(baseSalary*finalSkill) * roleMult)

	// Type Modifiers for Salary
	switch cType {
	case domain.CrewTypeAndroid:
		salary *= 2 // Expensive to maintain/purchase
	case domain.CrewTypeConvict:
		salary = int(float64(salary) * 0.2) // Dirt cheap labor (Indentured servitude)
	}

	return &entity.Crew{
		ID:     uuid.New(),
		Name:   identifier,
		Role:   role,
		Type:   cType,
		Skill:  finalSkill,
		Salary: salary,
	}
}