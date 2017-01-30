package digitalocean

import (
	"fmt"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"os"
)

type Do struct {
	token  string
	client *godo.Client
}
type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
func NewFromEnv() Do {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		panic("No Digitalocean token was found")
	}
	tokenSource := &TokenSource{
		AccessToken: token,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	return Do{token: token, client: client}
}
func (d *Do) FetchDnsAddrs() map[string][]string {
	result := make(map[string][]string)
	domains, _, err := d.client.Domains.List(&godo.ListOptions{
		Page:    1,
		PerPage: 500,
	})
	if err != nil {
		panic(err)
	}
	for _, domain := range domains {
		opt := &godo.ListOptions{Page: 1, PerPage: 300}
		records, _, err := d.client.Domains.Records(domain.Name, opt)
		if err != nil {
			fmt.Println(err)
		}
		for _, record := range records {
			if record.Type == "A" {
				var fulldomain string
				if record.Name == "@" {
					fulldomain = domain.Name
				} else {
					fulldomain = record.Name + "." + domain.Name
				}
				addr := record.Data
				result[addr] = append(result[addr], fulldomain)
			}
		}
	}
	return result
}
