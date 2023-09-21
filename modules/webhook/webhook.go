package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gennesseaux/NotionWatcher/modules/event"
	log "github.com/go-mods/zerolog-quick"
	"io"
	"net/http"
	"strings"
)

// PostHook : post event to webhook
func PostHook(url string, event event.Event) error {
	payload := new(bytes.Buffer)

	err := json.NewEncoder(payload).Encode(event)
	if err != nil {
		return err
	}

	log.Debug().Msgf("POST %v", strings.Trim(payload.String(), " \r\n"))

	// #nosec:G107
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(responseBody))
	}

	return nil

}
