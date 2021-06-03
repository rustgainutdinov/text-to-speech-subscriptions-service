package infrastructure

import "github.com/kelseyhightower/envconfig"

const appID = "balanceservice"

type config struct {
	ServeRESTAddress string `envconfig:"serve_rest_address" default:":8000"`
	DBUser           string `envconfig:"db_user"`
	DBName           string `envconfig:"db_name"`
	DBPass           string `envconfig:"db_pass"`
}

func ParseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, err
	}
	return c, nil
}
