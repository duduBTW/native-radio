package lib

import "encoding/json"

func ParseNodeOutput(output []byte) (map[string]Song, error) {
	var data map[string]Song
	if err := json.Unmarshal(output, &data); err != nil {
		return data, err
	}

	return data, nil
}
