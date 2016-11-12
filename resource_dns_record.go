package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

type domainRecordResource struct {
	Customer string
	Domain   string
	Records  []*DomainRecord
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
		Delete: resourceDomainRecordUpdate,

		Schema: map[string]*schema.Schema{
			// Required
			"customer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
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
	log.Println("Fetching domain records...")
	records, err := client.GetDomainRecords(
		d.Get("customer").(string),
		d.Get("domain").(string),
	)
	if err != nil {
		return fmt.Errorf("couldn't find domain record: ", err.Error())
	}
	return populateResourceDataFromResponse(records, d)
}

func resourceDomainRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GoDaddyClient)
	r, err := newDomainRecordResource(d)
	if err != nil {
		return err
	}
	log.Println("Updating", r.Domain, "domain records...")
	return client.UpdateDomainRecords(r.Customer, r.Domain, r.Records)
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
