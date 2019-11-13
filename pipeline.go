package amocrm

import (
	"encoding/json"
	"strconv"
)

func (c *ClientInfo) GetPipelines(reqParams *PipelineRequestParams) (*GetPipelineResponse, error) {
	addValues := map[string]string{}
	pipelineResponse := new(GetPipelineResponse)
	if err := Validate.Struct(reqParams); err != nil {
		return nil, err
	}

	if reqParams.ID != 0 {
		addValues["id"] = strconv.Itoa(reqParams.ID)
	}

	url := c.Url + apiUrls["pipelines"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return pipelineResponse, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(body, pipelineResponse)
	if err != nil {
		return nil, err
	}

	if pipelineResponse.Response != nil {
		return nil, pipelineResponse.Response
	}

	if err := Validate.Struct(pipelineResponse); err != nil {
		return nil, err
	}

	return pipelineResponse, nil
}
