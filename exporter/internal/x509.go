package internal

import (
	"fmt"
	"os"

	"github.com/artarts36/certmetrics/x509m"
)

func InspectX509Pem(pem PEMFile) error {
	file, err := os.ReadFile(pem.Path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	x509m.InspectPEMs(file, x509m.WithID(pem.ID))

	return nil
}
