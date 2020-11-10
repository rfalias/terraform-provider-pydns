package pydns

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/rfalias/gopydns"

	"os"
	"time"
	"math/rand"
)

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func waitForLock(client *DNSClient) bool {
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(100)) * time.Millisecond)

	locked := fileExists(client.lockfile)

	for locked == true {
		time.Sleep(100 * time.Millisecond)
		locked = fileExists(client.lockfile)
	}

	time.Sleep(1000 * time.Millisecond)
	return true
}

func resourcePyDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourcePyDNSRecordCreate,
		Read:   resourcePyDNSRecordRead,
		Delete: resourcePyDNSRecordDelete,

		Schema: map[string]*schema.Schema{
			"zone_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"record_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"record_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv4address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hostnamealias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ptrdomainname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePyDNSRecordCreate(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*DNSClient)

	zone_name := d.Get("zone_name").(string)
	record_type := d.Get("record_type").(string)
	record_name := d.Get("record_name").(string)
	ipv4address := d.Get("ipv4address").(string)
	//hostnamealias := d.Get("hostnamealias").(string)
	//ptrdomainname := d.Get("ptrdomainname").(string)

	var id string = zone_name + "_" + record_name + "_" + record_type

	//var psCommand string

	waitForLock(client)
	
	file, err := os.Create(client.lockfile)
	if err != nil {
		return err
	}


	_, err = gopydns.RunPyDnsCommandCreate(record_name, client.username, client.password, client.server, zone_name, ipv4address, client.dnspy)

	if err != nil {
		//something bad happened
		return err
	}

	d.SetId(id)

	file.Close()
	os.Remove(client.lockfile)

	return nil
}


func resourcePyDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePyDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	  //convert the interface so we can use the variables like username, etc
        client := m.(*DNSClient)

        zone_name := d.Get("zone_name").(string)
        record_type := d.Get("record_type").(string)
        record_name := d.Get("record_name").(string)
        ipv4address := d.Get("ipv4address").(string)
        //hostnamealias := d.Get("hostnamealias").(string)
        //ptrdomainname := d.Get("ptrdomainname").(string)

        var id string = zone_name + "_" + record_name + "_" + record_type

        //var psCommand string

        waitForLock(client)

        file, err := os.Create(client.lockfile)
        if err != nil {
                return err
        }


        _, err = gopydns.RunPyDnsCommandRemove(record_name, client.username, client.password, client.server, zone_name, ipv4address, client.dnspy)

        if err != nil {
                //something bad happened
                return err
        }

        d.SetId(id)

        file.Close()
        os.Remove(client.lockfile)

        return nil

}
