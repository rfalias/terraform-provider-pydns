package pydns

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"fmt"
	"os"
	"io/ioutil"
)

// Provider allows making changes to Windows DNS server
// Utilises Powershell to connect to domain controller
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("USERNAME", nil),
				Description: "Username to connect to AD.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PASSWORD", nil),
				Description: "The password to connect to AD.",
			},
			"server": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVER", nil),
				Description: "The AD server to connect to.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NAME", false),
				Description: "DNS Record Name",
			},
                        "zone": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("ZONE", false),
                                Description: "DNS Zone",
                        },
                        "ip": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("IP", false),
                                Description: "IP Of the A record",
                        },
                        "dnspy": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("DNSPY", false),
                                Description: "Path to the python script file",
                        },
		},
		ResourcesMap: map[string]*schema.Resource{
			"pydns": resourcePyDNSRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, fmt.Errorf("The 'username' property was not specified.")
	}
	
        password := d.Get("password").(string)
        if password == ""{
                return nil, fmt.Errorf("The 'password' property was not specified.")
        }

	server := d.Get("server").(string)
	if server == "" {
		return nil, fmt.Errorf("The 'server' property was not specified.")
	}

	name := d.Get("name").(string)
	zone := d.Get("zone").(string)
	ip := d.Get("ip").(string)
	dnspy := d.Get("dnspy").(string)

	f, err := ioutil.TempFile("", "terraform-pydns")
	lockfile := f.Name()
	os.Remove(f.Name())

	client := DNSClient {
		username:	username,
		password:	password,
		server:		server,
		name:		name,
		zone:           zone,
                ip:             ip,
		lockfile:       lockfile,
	}

	return &client, err
}

type DNSClient struct {
	username	string
	password	string
	server		string
	name		string
	zone            string
        ip              string
	lockfile	string
}
