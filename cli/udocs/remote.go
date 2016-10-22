package udocs

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRemotePage(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("%d %s %s", code, http.StatusText(code), url)
	}

	return ioutil.ReadAll(resp.Body)
}
