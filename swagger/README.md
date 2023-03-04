# swagger

This is for swaggo extension library to set extra variables.

```sh
go get github.com/worldline-go/utility/swagger
```

__OAuth2Application__ for service authentication.  
__OAuth2Implicit__ for user authentication but need to set implicit flow!  
__OAuth2Password__ for user authentication with directly with username and password.  
__OAuth2AccessCode__ for user authentication with code, it will redirect to the authUrl to login.  


## Usage

Create a file in the __/docs/__ folder, add to that folder to not include again generated file of swaggo.

Example __/docs/info.go__:

```go
package docs

import (
	"github.com/worldline-go/auth"
	"github.com/worldline-go/utility/swagger"
)

func Info(version string, provider auth.InfProvider) error {
	return swagger.SetInfo(
		swagger.WithVersion(version),
		swagger.WithCustom(map[string]interface{}{
			"tokenUrl": provider.GetTokenURL(),
			"authUrl":  provider.GetAuthURL(),
		}),
	)
}
```

Inside of the general information of server:

```go
// Echo Server
//
// @title Auth Test API
// @description This is a sample server for out Auth library.
//
// @contant.name worldline-go
// @contant.url https://github.com/worldline-go
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Description for what is this security definition being used

// @securitydefinitions.oauth2.application	OAuth2Application
// @tokenUrl								[[ .Custom.tokenUrl ]]

// @securitydefinitions.oauth2.implicit	OAuth2Implicit
// @authorizationUrl						[[ .Custom.authUrl ]]

// @securitydefinitions.oauth2.password	OAuth2Password
// @tokenUrl								[[ .Custom.tokenUrl ]]

// @securitydefinitions.oauth2.accessCode	OAuth2AccessCode
// @tokenUrl								[[ .Custom.tokenUrl ]]
// @authorizationUrl						[[ .Custom.authUrl ]]
func echoServer(ctx context.Context) error {
    // ...
}
```

Use OAuth2 tags in the requests

```go
// PostValue return the body
//
// @Summary Post value
// @Description Post value
// @Tags restricted
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /value [post]
// @Security ApiKeyAuth || OAuth2Application || OAuth2Implicit || OAuth2Password || OAuth2AccessCode
func (API) PostValue(c echo.Context) error {
    // ...
}
```
