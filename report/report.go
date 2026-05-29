package report

import (
	"encoding/json"
	"os"
	"slices"

	"github.com/alexis-wizeline/boot_crawler/scrapper"
)

func WriteJsonReport(pages map[string]scrapper.PageData, filename string) error {

	keys := make([]string, 0, len(pages))
	for key := range pages {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	sorted := make([]scrapper.PageData, len(keys))
	for i, key := range keys {
		sorted[i] = pages[key]
	}

	data, err := json.MarshalIndent(sorted, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
