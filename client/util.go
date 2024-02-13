package client

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func envSearch(s string) string {
	if v := os.Getenv(s); v != "" {
		return v
	}

	return ""
}

// returns: (path or lodade content, isPath, error)
//
//nolint:all
func pathOrContents(poc string) (string, bool, error) {
	if len(poc) == 0 {
		return poc, false, nil
	}

	path := filepath.Clean(poc)
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, true, err
		}
	}

	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return string(contents), true, err
		}
		return string(contents), true, nil
	}

	return poc, false, nil
}

func toMapWithJSONTags(i interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	a, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(a, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
