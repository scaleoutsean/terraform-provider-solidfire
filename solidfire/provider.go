package solidfire

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
				DefaultFunc: schema.EnvDefaultFunc("SOLIDFIRE_USERNAME", nil),
				Description: "The user name for ElementSW API operations.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDFIRE_PASSWORD", nil),
				Description: "The user password for ElementSW API operations.",
			},
			"solidfire_server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDFIRE_SERVER", nil),
				Description: "The ElementSW server name for ElementSW API operations.",
			},
			"api_version": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDFIRE_API_VERSION", nil),
				Description: "The ElementSW server API version.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"solidfire_volume_access_group": resourceElementSwVolumeAccessGroup(),
			"solidfire_initiator":           resourceElementSwInitiator(),
			"solidfire_volume":              resourceElementSwVolume(),
			"solidfire_account":             resourceElementSwAccount(),
			"solidfire_qos_policy":          resourceElementswQoSPolicy(),
			"solidfire_schedule":            resourceElementswSchedule(),
			"solidfire_snapshot":            resourceElementswSnapshot(),
			"solidfire_cluster_pairing":     resourceElementSwClusterPairing(),
			"solidfire_volume_pairing":      resourceElementSwVolumePairing(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"solidfire_cluster":             dataSourceElementSwCluster(),
			"solidfire_account":             dataSourceElementSwAccount(),
			"solidfire_volume":              dataSourceElementSwVolume(),
			"solidfire_volume_iqn":          dataSourceElementSwVolumeIQN(),
			"solidfire_cluster_stats":       dataSourceElementSwClusterStats(),
			"solidfire_volumes_by_account":  dataSourceElementswVolumesByAccount(),
			"solidfire_qos_policy":          dataSourceElementSwQosPolicy(),
			"solidfire_initiator":           dataSourceElementSwInitiator(),
			"solidfire_volume_access_group": dataSourceElementSwVolumeAccessGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	server := d.Get("solidfire_server").(string)
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
