package usecases_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/johnldev/4-deploy-cloud-run/internal/services"
	usecases "github.com/johnldev/4-deploy-cloud-run/internal/useCases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const WEATHER_API_TOKEN = "123"

type TemperatureSuite struct {
	suite.Suite
	usecase *usecases.GetTemperatureUseCase
}

var ctx context.Context

func (suite *TemperatureSuite) SetupSuite() {
	os.Setenv("WEATHER_API_TOKEN", WEATHER_API_TOKEN)
	ctx = context.Background()

	suite.usecase = usecases.NewGetTemperatureUseCase(
		services.NewCepService(ctx),
		services.NewWeatherService(ctx),
	)
}

func (suite *TemperatureSuite) TearDownSuite() {
	os.Unsetenv("WEATHER_API_TOKEN")
	ctx.Done()
}

func (suite *TemperatureSuite) TestGetTemperatureUseCase_Execute_InvalidZipCode() {
	assert := assert.New(suite.T())

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipcode := "96030610"
	var cepForCdn string = zipcode[:5] + "-" + zipcode[5:]

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipcode),
		httpmock.NewStringResponder(404, ""))

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cepForCdn),
		httpmock.NewStringResponder(404, ""))

	response, err := suite.usecase.Execute(ctx, "1234567890")
	assert.NotNil(err)
	assert.Nil(response)
	assert.Equal("can not find zipcode", err.Error())
}

func (suite *TemperatureSuite) TestGetTemperatureUseCase_Execute_InvalidCity() {
	assert := assert.New(suite.T())

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipcode := "96030610"
	city := "ab231c"
	var cepForCdn string = zipcode[:5] + "-" + zipcode[5:]

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipcode),
		httpmock.NewStringResponder(200, fmt.Sprintf(`{"localidade":"%s"}`, city)))

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cepForCdn),
		httpmock.NewStringResponder(200, fmt.Sprintf(`{"city":"%s"}`, city)))

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.weatherapi.com/v1/current.json?q=%s&key=%s", city, WEATHER_API_TOKEN),
		httpmock.NewErrorResponder(fmt.Errorf("can not find city")))

	response, err := suite.usecase.Execute(ctx, zipcode)
	assert.NotNil(err)
	assert.Nil(response)
	assert.Contains(err.Error(), "can not find city")
}

func (s *TemperatureSuite) TestGetTemperatureUseCase_Success() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipcode := "96030610"
	city := "Pelotas"
	var cepForCdn string = zipcode[:5] + "-" + zipcode[5:]

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipcode),
		httpmock.NewStringResponder(200, fmt.Sprintf(`{"localidade":"%s"}`, city)))

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cepForCdn),
		httpmock.NewStringResponder(200, fmt.Sprintf(`{"city":"%s"}`, city)))

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.weatherapi.com/v1/current.json?q=%s&key=%s", strings.ToLower(strings.ReplaceAll(city, " ", "-")), WEATHER_API_TOKEN),
		httpmock.NewStringResponder(200, `{"current": {
			"temp_c": 20,
			"temp_f": 68
		}}`))

	assert := assert.New(s.T())
	response, err := s.usecase.Execute(ctx, zipcode)
	assert.Nil(err)

	assert.NotEmpty(response.Celcius)
	assert.NotEmpty(response.Fahrenheit)
	assert.NotEmpty(response.Kelvin)
	assert.Equal(response.Kelvin, 293.15)
}

func TestGetTemperatureUseCase(t *testing.T) {
	suite.Run(t, new(TemperatureSuite))
}
