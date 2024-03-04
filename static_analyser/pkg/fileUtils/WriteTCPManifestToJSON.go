package file_utils

import (
	"encoding/json"
	"fmt"
	"os"
	t "static_analyser/pkg/types"
)

func WriteTCPManifestToJSON(
	manifest t.TCPManifest,
	serviceName string,
	outputPrefix string,
) error {
	// WriteTCPManifestToJSON converts a TCPManifest to JSON and writes it to a file.
	//
	// manifest: The TCPManifest to convert to JSON.
	// serviceName: The name of the service for error reporting.
	// outputPrefix: The prefix for the output file name.
	//
	// Returns:
	// An error if there was a problem converting the TCPManifest to JSON or writing the file.

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
