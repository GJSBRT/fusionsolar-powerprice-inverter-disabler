package fusionsolar

import (
	"net/http"
)

type Configuration struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Fusionsolar struct {
	configuration 	Configuration
	httpClient 		*http.Client
	roarand			string
}

func NewFusionsolar(configuration Configuration) *Fusionsolar {
	return &Fusionsolar{
		configuration: 	configuration,
		httpClient:		&http.Client{},
		roarand:		"",
	}
}
