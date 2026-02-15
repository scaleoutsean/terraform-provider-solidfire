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
			"elementsw_qos_policy":          resourceElementswQoSPolicy(),
			"elementsw_schedule":            resourceElementswSchedule(),
			"elementsw_snapshot":            resourceElementswSnapshot(),
			"elementsw_cluster_pairing":     resourceElementSwClusterPairing(),
			"elementsw_volume_pairing":      resourceElementSwVolumePairing(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"elementsw_cluster":             dataSourceElementSwCluster(),
			"elementsw_account":             dataSourceElementSwAccount(),
			"elementsw_volume":              dataSourceElementSwVolume(),
			"elementsw_volume_iqn":          dataSourceElementSwVolumeIQN(),
			"elementsw_cluster_stats":       dataSourceElementSwClusterStats(),
			"elementsw_volumes_by_account":  dataSourceElementswVolumesByAccount(),
			"elementsw_qos_policy":          dataSourceElementSwQosPolicy(),
			"elementsw_initiator":           dataSourceElementSwInitiator(),
			"elementsw_volume_access_group": dataSourceElementSwVolumeAccessGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	server := d.Get("elementsw_server").(string)
	version := d.Get("api_version").(string)
	user := d.Get("username").(string)
	config := configStuct{
		User:            user,
		Password:        d.Get("password").(string),
		ElementSwServer: server,
		APIVersion:      version,
	}

	return config.clientFun()
}
