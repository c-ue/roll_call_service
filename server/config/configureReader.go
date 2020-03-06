package config

type (
	// Server holds supported types by the multiconfig package
	Server struct {
		Name     string
		Port     int `default:"6060"`
		Enabled  bool
		Users    []string
		Postgres Postgres
	}

	// Postgres is here for embedded struct feature
	Postgres struct {
		Enabled           bool
		Port              int
		Hosts             []string
		DBName            string
		AvailabilityRatio float64
	}
)
