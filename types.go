package iam

import (
	"github.com/imroc/req/v3"
	http_mw "github.com/zitadel/zitadel-go/v2/pkg/api/middleware/http"
)

type AccountClient struct {
	*req.Client
	*http_mw.IntrospectionInterceptor
	*ZitadelUser
	UserId *string
}
type ZitadelUser struct {
	Email              string      `json:"email,omitempty"`
	EmailVerified      bool        `json:"email_verified,omitempty"`
	FamilyName         string      `json:"family_name,omitempty"`
	GivenName          string      `json:"given_name,omitempty"`
	Locale             string      `json:"locale,omitempty"`
	Name               string      `json:"name,omitempty"`
	PreferredUsername  string      `json:"preferred_username,omitempty"`
	Sub                string      `json:"sub,omitempty"`
	UpdatedAt          int         `json:"updated_at,omitempty"`
	UrnZitadelIamRoles interface{} `json:"urn:zitadel:iam:org:project:roles,omitempty"`
}
