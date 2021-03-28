package skycmd

import (
	"encoding/json"
	"errors"

	"github.com/concourse/dex/connector/google"
	"github.com/hashicorp/go-multierror"
)

func init() {
	RegisterConnector(&Connector{
		id:         "google",
		config:     &GoogleFlags{},
		teamConfig: &GoogleTeamFlags{},
	})
}

type GoogleFlags struct {
	DisplayName            string   `long:"display-name" description:"The auth provider name displayed to users on the login page"`
	ClientID               string   `long:"client-id" description:"(Required) Client id"`
	ClientSecret           string   `long:"client-secret" description:"(Required) Client secret"`
	Scopes                 []string `long:"scope" description:"Any additional scopes that need to be requested during authorization. Default to [profile, email]."`
	HostedDomains          []string `long:"hosted-domains" description:"List of whitelisted domains, only users from a listed domain will be allowed to log in"`
	Groups                 []string `long:"groups" descripton:"If this field is nonempty, only users from a listed group will be allowed to log in"`
	ServiceAccountFilePath string   `long:"service-account-file-path" description:"If nonempty, and groups claim is made, will use authentication from file to check groups with the admin directory api"`
	AdminEmail             string   `long:"admin-email" descripton:"The email of a GSuite super user which the service account will impersonate when listing groups"`
}

func (flag *GoogleFlags) Name() string {
	if flag.DisplayName != "" {
		return flag.DisplayName
	}
	return "Google"
}

func (flag *GoogleFlags) Validate() error {
	var errs *multierror.Error

	if flag.ClientID == "" {
		errs = multierror.Append(errs, errors.New("Missing client-id"))
	}

	if flag.ClientSecret == "" {
		errs = multierror.Append(errs, errors.New("Missing client-secret"))
	}

	return errs.ErrorOrNil()
}

func (flag *GoogleFlags) Serialize(redirectURI string) ([]byte, error) {
	if err := flag.Validate(); err != nil {
		return nil, err
	}

	config := google.Config{
		ClientID:               flag.ClientID,
		ClientSecret:           flag.ClientSecret,
		RedirectURI:            redirectURI,
		Scopes:                 flag.Scopes,
		HostedDomains:          flag.HostedDomains,
		Groups:                 flag.Groups,
		ServiceAccountFilePath: flag.ServiceAccountFilePath,
		AdminEmail:             flag.AdminEmail,
	}

	return json.Marshal(config)
}

type GoogleTeamFlags struct {
	Users  []string `json:"users" long:"user" description:"A whitelisted Google user" value-name:"USERNAME"`
	Groups []string `json:"groups" long:"group" description:"A whitelisted Google group" value-name:"GROUP_NAME"`
}

func (flag *GoogleTeamFlags) GetUsers() []string {
	return flag.Users
}

func (flag *GoogleTeamFlags) GetGroups() []string {
	return flag.Groups
}
