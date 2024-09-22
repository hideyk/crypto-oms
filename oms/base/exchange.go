package base

type Exchange struct {
	Id              string
	Name            string
	Countries       []string
	Version         string
	EnableRateLimit bool
	RateLimit       int // Milliseconds
	Timeout         int // Milliseconds
	Hostname        string
	Symbols         []string
	Codes           []string
	Fees            struct{}
	BaseEndpoint    string
	ApiKey          string
}
