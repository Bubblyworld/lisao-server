package github

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// StatusMap represents a map of HTTP status codes to JSON structs - the client
// can automatically unmarshall a response into the relevant struct based on the
// status code of the response.
type StatusMap map[int]interface{}

var ErrUnknownStatusCode = errors.New("github: unknown status code")

func (c *Client) DoToJSON(req *http.Request, statusMap StatusMap) (interface{}, error) {
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonStruct, ok := statusMap[res.StatusCode]
	if !ok {
		return nil, ErrUnknownStatusCode
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, jsonStruct)
	return jsonStruct, err
}
