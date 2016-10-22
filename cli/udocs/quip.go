package udocs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var DefaultQuipClient = NewQuipClient(os.Getenv("UDOCS_QUIP_ACCESS_TOKEN"))

type QuipClient struct {
	accessToken string
}

func NewQuipClient(accessToken string) *QuipClient {
	return &QuipClient{accessToken: accessToken}
}

func (qc *QuipClient) GetThread(id string) (*Thread, error) {
	req, err := http.NewRequest("GET", "https://platform.quip.com/1/threads/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+qc.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("thread not found: " + id)
	}

	var thread Thread
	if err := json.NewDecoder(resp.Body).Decode(&thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

func (qc *QuipClient) GetBlob(threadID, blobID string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://platform.quip.com/1/blob/%s/%s", threadID, blobID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+qc.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("blob not found: " + threadID + "/" + blobID)
	}
	return ioutil.ReadAll(resp.Body)
}

type Thread struct {
	ExpandedUserIds []string `json:"expanded_user_ids"`
	UserIds         []string `json:"user_ids"`
	SharedFolderIds []string `json:"shared_folder_ids"`
	HTML            string   `json:"html"`
	ThreadMetadata  struct {
		CreatedUsec int64  `json:"created_usec"`
		UpdatedUsec int64  `json:"updated_usec"`
		ID          string `json:"id"`
		Title       string `json:"title"`
		Link        string `json:"link"`
	} `json:"thread"`
}

func IsQuipBlob(filename string) bool {
	return strings.HasPrefix(filename, "/blob/")
}

func IsQuipThread(filename string) bool {
	return filepath.Ext(filename) == ".quip"
}
