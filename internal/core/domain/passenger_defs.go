package domain

type PassengerType int

const (
	PaxTourist  PassengerType = iota // Standard, low risk, medium pay
	PaxWorker                        // Low pay, bulk transport, very common
	PaxBusiness                      // Good pay, wants speed/safety
	PaxVIP                           // High pay, high demands (Luxury required)
	PaxRefugee                       // Very low pay (or charity), appears in bad systems
	PaxPrisoner                      // High pay, requires Brig/Secure hold, high risk
	PaxFugitive                      // Massive pay, Illegal status
)

func (p PassengerType) String() string {
	return [...]string{
		"Tourist", "Migrant Worker", "Business Executive", 
		"High-Value VIP", "Refugee", "Prisoner Transport", "Fugitive",
	}[p]
}