package amocrm

import "encoding/json"

func (c *ClientInfo) GetPipelines(reqParams *PipelineRequestParams) (*PipelineResponse, error) {
	addValues := map[string]string{}
	pipeline := new(PipelineResponse)
	var err error
	if reqParams.ID != "" {
		addValues["id"] = reqParams.ID
	}
	url := c.Url + apiUrls["pipelines"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return pipeline, err
	}
	err = json.Unmarshal(body, pipeline)
	if err != nil {
		return nil, err
	}

	if pipeline.Response != nil {
		return nil, pipeline.Response
	}

	return pipeline, nil
}
