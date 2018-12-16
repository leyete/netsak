package nsconfig

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	viper.Reset()

	v := viper.New()
	v.SetConfigFile("./nsconfig.yaml")
	err := v.ReadInConfig()

	c := ConfigurationNew()
	c.ReadFile("./nsconfig.yaml")

	assert.NoError(t, err)

	fmt.Println(c)
	
	assert.Equal(t, c.Server.Url , v.GetString("server.url"))
	assert.Equal(t, c.Server.Port , v.GetString("server.port"))
	assert.Equal(t, c.Server.Swagger , v.GetString("server.swagger"))
	assert.Equal(t, c.Server.User , v.GetString("server.user"))
	assert.Equal(t, c.Server.Password , v.GetString("server.password"))

}

func TestConfig2(t *testing.T) {
	viper.Reset()

	v := viper.New()
	v.SetConfigFile("./nsconfig.1.yaml")
	err := v.ReadInConfig()

	c := ConfigurationNew()
	c.ReadFile("./nsconfig.1.yaml")

	assert.NoError(t, err)

	fmt.Println(c)
	
	assert.Equal(t, c.Server.Url , v.GetString("server.url"))
	assert.Equal(t, c.Server.Port , v.GetString("server.port"))
	assert.Equal(t, c.Server.Swagger , v.GetString("server.swagger"))
	assert.Equal(t, c.Server.User , v.GetString("server.user"))
	assert.Equal(t, c.Server.Password , v.GetString("server.password"))

}