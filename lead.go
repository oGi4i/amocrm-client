package amocrm

import "errors"

func (c *ClientInfo) AddLead(lead Lead) (int, error) {
	if lead.Name == "" {
		return 0, errors.New("Name is empty")
	}
	if lead.StatusID == "" {
		return 0, errors.New("StatusID is empty")
	}
	url := c.Url + apiUrls["leads"]
	resp, err := c.DoPost(url, LeadSetRequest{Add: []Lead{lead}})
	if err != nil {
		return 0, err
	}
	return c.GetResponseID(resp)
}
