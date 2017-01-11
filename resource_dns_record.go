package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

type domainRecordResource struct {
	Customer string
	Domain   string
	Records  []*DomainRecord
}

var defaultRecords = []*DomainRecord{
	// A Records
	&DomainRecord{Type: "A", Name: "@", Data: "50.63.202.43", TTL: 600},
	// CNAME Records
	&DomainRecord{Type: "CNAME", Name: "email", Data: "email.secureserver.net", TTL: 3600},
	&DomainRecord{Type: "CNAME", Name: "ftp", Data: "@", TTL: 3600},
	&DomainRecord{Type: "CNAME", Name: "www", Data: "@", TTL: 3600},
	&DomainRecord{Type: "CNAME", Name: "_domainconnect", Data: "_domainconnect.gd.domaincontrol.com", TTL: 3600},
	// MX Records
	&DomainRecord{Type: "MX", Name: "@", Data: "mailstore1.secureserver.net", TTL: 3600, Priority: 10},
	&DomainRecord{Type: "MX", Name: "@", Data: "smtp.secureserver.net", TTL: 3600, Priority: 0},
	// NS Records
	&DomainRecord{Type: "NS", Name: "@", Data: "ns45.domaincontrol.com", TTL: 3600},
	&DomainRecord{Type: "NS", Name: "@", Data: "ns46.domaincontrol.com", TTL: 3600},
}

func newDomainRecordResource(d *schema.ResourceData) (domainRecordResource, error) {
	var err error
	r := domainRecordResource{}

	if attr, ok := d.GetOk("customer"); ok {
		r.Customer = attr.(string)
	}

	if attr, ok := d.GetOk("domain"); ok {
		r.Domain = attr.(string)
	}

	if attr, ok := d.GetOk("record"); ok {
		records := attr.(*schema.Set).List()
		r.Records = make([]*DomainRecord, len(records))

		for i, rec := range records {
			data := rec.(map[string]interface{})
			r.Records[i], err = NewDomainRecord(
				data["name"].(string),
				data["type"].(string),
				data["data"].(string),
				data["ttl"].(int))

			if err != nil {
				return r, err
			}
		}
	}

	return r, err
}

func resourceDomainRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainRecordUpdate,
		Read:   resourceDomainRecordRead,
		Update: resourceDomainRecordUpdate,
		Delete: resourceDomainRecordRestore,

		Schema: map[string]*schema.Schema{
			// Optional
			"customer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Required
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"record": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"data": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3600,
						},
					},
				},
			},
		},
	}
}

func resourceDomainRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GoDaddyClient)
	customer := d.Get("customer").(string)
	domain := d.Get("domain").(string)

	log.Println("Fetching", domain, "records...")
	records, err := client.GetDomainRecords(customer, domain)
	if err != nil {
		return fmt.Errorf("couldn't find domain record (%s): %s", domain, err.Error())
	}

	return populateResourceDataFromResponse(records, d)
}

func resourceDomainRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GoDaddyClient)
	r, err := newDomainRecordResource(d)
	if err != nil {
		return err
	}

	if err = populateDomainInfo(client, &r, d); err != nil {
		return err
	}

	log.Println("Updating", r.Domain, "domain records...")
	return client.UpdateDomainRecords(r.Customer, r.Domain, r.Records)
}

func resourceDomainRecordRestore(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GoDaddyClient)
	r, err := newDomainRecordResource(d)
	if err != nil {
		return err
	}

	if err = populateDomainInfo(client, &r, d); err != nil {
		return err
	}

	log.Println("Restoring", r.Domain, "domain records...")
	return client.UpdateDomainRecords(r.Customer, r.Domain, defaultRecords)
}

func populateDomainInfo(client *GoDaddyClient, r *domainRecordResource, d *schema.ResourceData) error {
	var err error
	var domain *Domain

	log.Println("Fetching", r.Domain, "info...")
	domain, err = client.GetDomain(r.Customer, r.Domain)
	if err != nil {
		return fmt.Errorf("couldn't find domain (%s): %s", r.Domain, err.Error())
	}

	d.SetId(strconv.FormatInt(domain.ID, 10))
	return nil
}

func populateResourceDataFromResponse(r []*DomainRecord, d *schema.ResourceData) error {
	d.Set("record", flattenRecords(r))
	return nil
}

func flattenRecords(list []*DomainRecord) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, r := range list {
		l := map[string]interface{}{
			"name": r.Name,
			"type": r.Type,
			"data": r.Data,
			"ttl":  r.TTL,
		}
		result = append(result, l)
	}
	return result
}
