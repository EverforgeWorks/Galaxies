package domain

// --- CREW ROLES (Job) ---
type CrewRole int

const (
	RoleFirstOfficer CrewRole = iota
	RoleWeaponsOfficer
	RoleNavigator
	RoleEngineer
	RoleSteward
	RoleSpecialist
)

func (r CrewRole) String() string {
	return [...]string{
		"First Officer", "Weapons Officer", "Navigator", 
		"Engineer", "Steward", "Specialist",
	}[r]
}

// --- CREW TYPES ---
type CrewType int

const (
	CrewTypeStandard CrewType = iota // Regular humans/aliens
	CrewTypeAndroid                  // Crafted, high skill, high cost
	CrewTypeConvict                  // Cheap, skilled, high risk
)

func (t CrewType) String() string {
	return [...]string{"Standard", "Android", "Convict"}[t]
}