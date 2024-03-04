package file_utils

import (
	"encoding/json"
	"fmt"
	"os"
	t "static_analyser/pkg/types"
)

// WriteTCPManifestToJSON writes the TCPManifest to a JSON file
func WriteTCPManifestToJSON(
	manifest t.TCPManifest, // The TCPManifest to write to a file
	serviceName string, // The name of the service
	outputPrefix string, // The prefix for the output file
) error { // Returns an error if marshalling or writing the file fails

	// Convert the manifest to JSON
	jsonData, err := json.MarshalIndent(manifest, "", " ")

	if err != nil {
		return fmt.Errorf("failed to marshal TCPManifest for service '%s': %w", serviceName, err)
	}

	// Write the JSON to a file
	filename := outputPrefix + manifest.Service + ".json"
	err = os.WriteFile(filename, jsonData, 0777) // consider using 0644 in future for more secure permissions

	if err != nil {
		return fmt.Errorf("failed to write TCPManifest to file '%s': %w", filename, err)
	}

	return nil
}
