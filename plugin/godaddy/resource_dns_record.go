package godaddy

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/n3integration/terraform-provider-godaddy/api"
)

const (
	attrCustomer    = "customer"
	attrDomain      = "domain"
	attrRecord      = "record"
	attrAddresses   = "addresses"
	attrNameservers = "nameservers"

	recName     = "name"
	recType     = "type"
	recData     = "data"
	recTTL      = "ttl"
	recPriority = "priority"
	recWeight   = "weight"
	recProto    = "protocol"
	recService  = "service"
	recPort     = "port"
)

type domainRecordResource struct {
	Customer         string
	Domain           string
	Records          []*api.DomainRecord
	ARecords         []string
	NSRecords        []string
	ReplaceNSRecords bool
}

var defaultRecords = []*api.DomainRecord{
	// CNAME Records
	{Type: api.CNameType, Name: "www", Data: "@", TTL: api.DefaultTTL},
	{Type: api.CNameType, Name: "_domainconnect", Data: "_domainconnect.gd.domaincontrol.com", TTL: api.DefaultTTL},
}

func newDomainRecordResource(d *schema.ResourceData) (*domainRecordResource, error) {
	var err error
	r := &domainRecordResource{}
	nsCount := 0

	if attr, ok := d.GetOk(attrCustomer); ok {
		r.Customer = attr.(string)
	}

	if attr, ok := d.GetOk(attrDomain); ok {
		r.Domain = attr.(string)
	}

	if attr, ok := d.GetOk(attrRecord); ok {
		records := attr.(*schema.Set).List()
		r.Records = make([]*api.DomainRecord, len(records))

		for i, rec := range records {
			data := rec.(map[string]interface{})
			t := data[recType].(string)
			if strings.EqualFold(t, api.NSType) {
				nsCount++
			}
			r.Records[i], err = api.NewDomainRecord(
				data[recName].(string),
				t,
				data[recData].(string),
				data[recTTL].(int),
				api.Priority(data[recPriority].(int)),
				api.Weight(data[recWeight].(int)),
				api.Port(data[recPort].(int)),
				api.Service(data[recService].(string)),
				api.Protocol(data[recProto].(string)))

			if err != nil {
				return r, err
			}
		}
	}

	if attr, ok := d.GetOk(attrNameservers); ok {
		records := attr.([]interface{})
		nsCount += len(records)
		r.NSRecords = make([]string, len(records))
		for i, rec := range records {
			if err = api.ValidateData(api.NSType, rec.(string)); err != nil {
				return r, err
			}
			r.NSRecords[i] = rec.(string)
		}
	}

	if attr, ok := d.GetOk(attrAddresses); ok {
		records := attr.([]interface{})
		r.ARecords = make([]string, len(records))
		for i, rec := range records {
			if err = api.ValidateData(api.AType, rec.(string)); err != nil {
				return r, err
			}
			r.ARecords[i] = rec.(string)
		}
	}

	if nsCount > 0 {
		r.ReplaceNSRecords = true
	}

	return r, err
}

func (r *domainRecordResource) converge() {
	r.mergeRecords(r.ARecords, api.NewARecord)
	r.mergeRecords(r.NSRecords, api.NewNSRecord)
}

func (r *domainRecordResource) mergeRecords(list []string, factory api.RecordFactory) {
	for _, data := range list {
		record, _ := factory(data)
		r.Records = append(r.Records, record)
	}
}

func resourceDomainRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainRecordCreate,
		ReadContext:   resourceDomainRecordRead,
		UpdateContext: resourceDomainRecordUpdate,
		DeleteContext: resourceDomainRecordRestore,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
			attrDomain: {
				Type:     schema.TypeString,
				Required: true,
			},
			// Optional
			attrAddresses: {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			attrCustomer: {
				Type:     schema.TypeString,
				Optional: true,
			},
			attrNameservers: {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			attrRecord: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						recName: {
							Type:     schema.TypeString,
							Required: true,
						},
						recType: {
							Type:     schema.TypeString,
							Required: true,
						},
						recData: {
							Type:     schema.TypeString,
							Required: true,
						},
						recTTL: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  api.DefaultTTL,
						},
						recPriority: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  api.DefaultPriority,
						},
						recWeight: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  api.DefaultWeight,
						},
						recService: {
							Type:     schema.TypeString,
							Optional: true,
						},
						recProto: {
							Type:     schema.TypeString,
							Optional: true,
						},
						recPort: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  api.DefaultPort,
						},
					},
				},
			},
		},
	}
}

func resourceDomainRecordRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)
	customer := d.Get(attrCustomer).(string)
	domain := d.Get(attrDomain).(string)
	r, err := newDomainRecordResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Importer support
	if domain == "" {
		r.Domain = d.Id()
		domain = r.Domain
	}

	log.Println("Fetching", domain, "records...")
	records, err := client.GetDomainRecords(customer, domain)
	if err != nil {
		return diag.FromErr(fmt.Errorf("couldn't find domain record (%s): %s", domain, err.Error()))
	}

	r.converge()
	if err := populateResourceDataFromResponse(records, r, d); err != nil {
		return diag.FromErr(err)
	}

	if err := populateDomainInfo(client, r, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDomainRecordCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)
	r, err := newDomainRecordResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = populateDomainInfo(client, r, d); err != nil {
		return diag.FromErr(err)
	}

<<<<<<< Updated upstream
	log.Println("Updating", r.Domain, "domain records...")
	r.converge()
	return client.UpdateDomainRecords(r.Customer, r.Domain, r.Records)
=======
	log.Println("Creating", r.Domain, "domain records...")
	r.converge()
	if err := client.UpdateDomainRecords(r.Customer, r.Domain, r.Records); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDomainRecordUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)
	r, err := newDomainRecordResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = populateDomainInfo(client, r, d); err != nil {
		return diag.FromErr(err)
	}

	log.Println("Updating", r.Domain, "domain records...")
	r.converge()
	if err := client.UpdateDomainRecords(r.Customer, r.Domain, r.Records); err != nil {
		return diag.FromErr(err)
	}
	return nil
>>>>>>> Stashed changes
}

func resourceDomainRecordRestore(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)
	r, err := newDomainRecordResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = populateDomainInfo(client, r, d); err != nil {
		return diag.FromErr(err)
	}

	log.Println("Restoring", r.Domain, "domain records...")
<<<<<<< Updated upstream
	return client.UpdateDomainRecords(r.Customer, r.Domain, defaultRecords)
=======
	if err := client.UpdateDomainRecords(r.Customer, r.Domain, defaultRecords); err != nil {
		return diag.FromErr(err)
	}
	return nil
>>>>>>> Stashed changes
}

func populateDomainInfo(client *api.Client, r *domainRecordResource, d *schema.ResourceData) error {
	var err error
	var domain *api.Domain

	log.Println("Fetching", r.Domain, "info...")
	domain, err = client.GetDomain(r.Customer, r.Domain)
	if err != nil {
		return fmt.Errorf("couldn't find domain (%s): %s", r.Domain, err.Error())
	}

	d.SetId(strconv.FormatInt(domain.ID, 10))
	return nil
}

func populateResourceDataFromResponse(recs []*api.DomainRecord, r *domainRecordResource, d *schema.ResourceData) error {
	aRecords := make([]string, 0)
	nsRecords := make([]string, 0)
	domain := d.Get(attrDomain).(string)
	records := make([]*api.DomainRecord, 0)

	for _, rec := range recs {
		switch {
		case api.IsDefaultNSRecord(rec):
			nsRecords = append(nsRecords, rec.Data)
		case api.IsDefaultARecord(rec):
			aRecords = append(aRecords, rec.Data)
		default:
			records = append(records, rec)
		}
	}

	if err := d.Set(attrAddresses, aRecords); err != nil {
		return err
	}

	if r.ReplaceNSRecords || domain == "" {
		if err := d.Set(attrNameservers, nsRecords); err != nil {
			return err
		}
	}

	if err := d.Set(attrRecord, flattenRecords(records)); err != nil {
		return err
	}

	if domain == "" {
		d.Set(attrDomain, d.Id())
	}

	return nil
}

func flattenRecords(list []*api.DomainRecord) []map[string]interface{} {
	result := make([]map[string]interface{}, len(list))
	for i, r := range list {
		result[i] = map[string]interface{}{
			recName:     r.Name,
			recType:     r.Type,
			recData:     r.Data,
			recTTL:      r.TTL,
			recPriority: r.Priority,
			recWeight:   r.Weight,
			recPort:     r.Port,
			recService:  r.Service,
			recProto:    r.Protocol,
		}
	}
	return result
}
