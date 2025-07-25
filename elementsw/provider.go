package elementsw

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the Terraform provider definition for ElementSW
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELEMENTSW_USERNAME", nil),
				Description: "The user name for ElementSW API operations.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELEMENTSW_PASSWORD", nil),
				Description: "The user password for ElementSW API operations.",
			},
			"elementsw_server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELEMENTSW_SERVER", nil),
				Description: "The ElementSW server name for ElementSW API operations.",
			},
			"api_version": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELEMENTSW_API_VERSION", nil),
				Description: "The ElementSW server API version.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"elementsw_volume_access_group": resourceElementSwVolumeAccessGroup(),
			"elementsw_initiator":           resourceElementSwInitiator(),
			"elementsw_volume":              resourceElementSwVolume(),
			"elementsw_account":             resourceElementSwAccount(),
			"elementsw_qos_policy":         resourceElementswQoSPolicy(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := configStuct{
		User:            d.Get("username").(string),
		Password:        d.Get("password").(string),
		ElementSwServer: d.Get("elementsw_server").(string),
		APIVersion:      d.Get("api_version").(string),
	}

	return config.clientFun()
}
