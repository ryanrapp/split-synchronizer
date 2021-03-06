package recorder

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ImpressionListenerSubmitter struct {
	Endpoint string
}

// This struct is used to put together all the data posted by the impression's listener
type ImpressionListenerPostBody struct {
	Impressions json.RawMessage `json:"impressions"`
	SdkVersion  string          `json:"sdkVersion"`
	MachineIP   string          `json:"machineIP"`
	MachineName string          `json:"machineName"`
}

// Impression queue sizes
var ImpressionListenerMainQueueSize int = 10
var ImpressionListenerFailedQueueSize int = 10

func (r ImpressionListenerSubmitter) Post(
	impressions json.RawMessage,
	sdkVersion string,
	machineIP string,
	machineName string) error {

	client := &http.Client{}

	bundle := &ImpressionListenerPostBody{
		Impressions: impressions,
		SdkVersion:  sdkVersion,
		MachineIP:   machineIP,
		MachineName: machineName,
	}

	data, err := json.Marshal(bundle)
	if err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", r.Endpoint, bytes.NewBuffer(data))
	request.Close = true
	response, err := client.Do(request)
	if err != nil {
		return err
	} else {
		defer response.Body.Close()
	}

	return nil
}
