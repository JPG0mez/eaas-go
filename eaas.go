package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

var spfrecords = map[string]string{
	"_spf.salesforce.com":      "Domain is allowing emails to be sent from Salesforce.com. This indicates a high liklihood of subscription to Salesforce services.",
	"_spf.google.com":          "Domain is allowing emails to be sent from Google.com. This indicates the domain may utilize Gmail or other G-Suite product offerings.",
	"protection.outlook.com":   "Domain is allowing emails to be sent from Microsoft.com. This strongly indicates the domain is utilzing Microsoft Hosted Exchange.",
	"service-now.com":          "Domain is allowing emails to be sent from service-now.com. This strongly indicates the domain is using the Service Now helpdesk platform.",
	"mailsenders.netsuite.com": "Domain is allowing emails to be sent from NetSuite. This strongly indicates the domain is using NetSuite product offerings (ERP / Cloud Accounting).",
	"mktomail.com":             "Domain is allowing emails to be sent from Marketo. This strongly indicates the domain is using the Marketo Marketing and Lead Generation platform.",
	"spf.mandrillapp.com":      "Domain is allowing emails to be sent from Mandrill (MailChimp). This strongly indicates the domain is using the Mandrill product offering for transactional email.",
	"pphosted.com":             "Domain is allowing emails to be sent from Proof Point. This strongly indicates the domain is utilizing Proof Point Managed email services.",
	"zendesk.com":              "Domain is allowing emails to be sent from Zendesk. This strongly indicates the domain is utilizing Zendesk for help desk and ticketing purposes.",
	"mcsv.net":                 "Domain is allowing emails to be sent from MailChimp. This strongly indicates the domain is utilizing MailChimp for email marketing.",
	"freshdesk.com":            "Domain is allowing emails to be sent from Freshdesk. This strongly indicates the domain is utilizing FreshDesk for helpdesk and ticketing services.",
}

var txtrecords = map[string]string{
	"docusign":                      "This record is used as proof of domain ownership for DocuSign product offerings. This indicates the domain likely uses DocuSign as an e-signature solution.",
	"facebook-domain-verification":  "This record is used as proof of domain ownership for use in the Facebook Business Manager.",
	"google-site-verification":      "This record is used as proof of ownership for Google G Suite product offerings, however it might simply be verification for Google Analytics.",
	"adobe-sign-verification":       "This record is used as proof of ownership of a domain for the Adobe Sign product offering. This indicates the domain likely uses Adobe Sign as an e-signature solution.",
	"atlassian-domain-verification": "This record verifies domain ownership with Atlassian. This indicates the domain might be sending managed emails from an Atlassian property and likely utilizes Atalassian product offerings such as Jira of Confluence.",
	"MS":                            "This record is used as proof of domain ownership by the Microsoft Office 365 product offering. This indicates a probable usage of at least some Office 365 products.",
	"adobe-idp-site-verification":   "This record is used as proof of ownership for use in the Adobe Enterprise Products product offerings.",
	"yandex-verification":           "This record is used as proof of ownership of a domain for Yandex. This indicates a probable usage of the Yandex Webmaster Tools.",
	"_amazonses":                    "Amazon Simple Email Services",
	"logmein-verification-code":     "This record is used as proof of ownership for a domain for LogMeIn services. The presence of this record indicates the owner of the domain is likely using LogMeIn for remote troubleshooting.",
	"citrix-verification-code":      "This record indicates the domain might be associated with utilizing Citrix Services.",
	"pardot":                        "This record indicates the domain might be utilizing Pardot B2B Marketing tools from Salesforce.com.",
	"zuora":                         "This record indicates the domain might be utilizing Zuora subscription management software.",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("EaaS - Enumeration as a Service.")
		fmt.Println("Usage : ./eaas [domain]")
		os.Exit(1)
	}

	domain := os.Args[1]

//	fmt.Printf("[*] EaaS - Enumeration as a Service script started.\n")
//	fmt.Printf("[*] Performing queries on domain %s\n", domain)

	queryTXT(domain)
}

func queryTXT(domain string) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, net.JoinHostPort("8.8.8.8", "53"))

	if r == nil {
		fmt.Printf("[-] Error: %s\n", err.Error())
		return
	}

	if r.Rcode != dns.RcodeSuccess {
		fmt.Printf("[-] Invalid answer name %s after MX query for %s\n", os.Args[1], domain)
	}

	// Stuff must be in the answer section
	for _, a := range r.Answer {
		if txt, ok := a.(*dns.TXT); ok {
			for key, value := range txtrecords {
				if strings.Contains(txt.Txt[0], key) {
					info := map[string]interface{}{
						"Domain": domain,
						"Record": key,
						"Info":   value,
						"TXT":    txt.Txt,
					}
					jsonInfo, _ := json.Marshal(info)
					fmt.Println(string(jsonInfo))
				}
			}
			for spfkey, spfvalue := range spfrecords {
				if strings.Contains(txt.Txt[0], spfkey) {
					spfInfo := map[string]interface{}{
						"Domain": domain,
						"Record": spfkey,
						"Info":   spfvalue,
						"TXT":    txt.Txt,
					}
					jsonSpfInfo, _ := json.Marshal(spfInfo)
					fmt.Println(string(jsonSpfInfo))
				}
			}
		}
	}
}

