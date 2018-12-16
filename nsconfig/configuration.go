package nsconfig

import (
	"github.com/spf13/viper"
	"github.com/mattrbeam/netsak/nsconfig/models"
	"os"
	"errors"
)

type Configuration struct {
	v *viper.Viper

	Server     models.ServerConfiguration
}

func ConfigurationNew() *Configuration {
	viper.Reset()

	configuration := new(Configuration)
	configuration.v = viper.New()

	return configuration
}

// ReadFile reads a config file from nsconfig directory
// fileName should be without the extension
func (config *Configuration) ReadFile(fileName string) error {
	err := existsFile(fileName)

	config.v.SetConfigType("yaml")
	config.v.SetConfigFile(fileName)
	config.v.ReadInConfig()

	if config.v.IsSet("server") {
		err:= assignValues(config,"server")
		if err != nil{
			return err
		}
	}
	return err
}

func assignValues(config *Configuration, typeOf string) error {
switch typeOf{
	case "server":
		var server models.ServerConfiguration
		vAuxServer := config.v.Sub("server")
		if err := vAuxServer.Unmarshal(&server); err != nil {
			return err
		}
		config.Server = server
		serverErr := config.Server.ServerValidator()
		if serverErr != nil {
			return serverErr
		}
	default:
		errors.New("typeOf undefined")

}

return nil

}

//This function check if the file exists
func existsFile(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return err
	}
	return nil
}


