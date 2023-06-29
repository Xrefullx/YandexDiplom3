package models

type Config struct {
	Address        string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DataBaseURI    string `env:"DATABASE_URI" envDefault:"user=postgres password=123qwe dbname=loyality sslmode=disable"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"localhost:8080"`
	SecretKey      string `env:"SECRET_KEY" envDefault:"Xrefullx"`
	ReleaseMOD     bool   `env:"RELEASE_MODE" envDefault:"false"`
}
