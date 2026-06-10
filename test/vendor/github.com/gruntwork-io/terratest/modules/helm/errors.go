package helm

import (
	"errors"
	"fmt"
)

// ErrNoMatchingData is returned when no matching raw data are found for the concrete type during YAML unmarshalling.
var ErrNoMatchingData = errors.New("no matching raw data were found for the concrete type")

// ErrDestinationNotPointer is returned when the destination object passed to an unmarshal function is not a pointer.
var ErrDestinationNotPointer = errors.New("destinationObj must be a pointer")

// ErrYAMLRootIsArray is returned when the YAML root is an array but the destination object is a single object.
var ErrYAMLRootIsArray = errors.New("YAML root is an array, but destinationObj is a single object")

// ValuesFileNotFoundError is returned when a provided values file input is not found on the host path.
type ValuesFileNotFoundError struct {
	Path string
}

// Error implements the error interface for ValuesFileNotFoundError.
func (err ValuesFileNotFoundError) Error() string {
	return "could not resolve values file " + err.Path
}

// SetFileNotFoundError is returned when a provided set file input is not found on the host path.
type SetFileNotFoundError struct {
	Path string
}

// Error implements the error interface for SetFileNotFoundError.
func (err SetFileNotFoundError) Error() string {
	return "could not resolve set file path " + err.Path
}

// TemplateFileNotFoundError is returned when a provided template file input is not found in the chart.
type TemplateFileNotFoundError struct {
	Path     string
	ChartDir string
}

// Error implements the error interface for TemplateFileNotFoundError.
func (err TemplateFileNotFoundError) Error() string {
	return fmt.Sprintf("could not resolve template file %s relative to chart path %s", err.Path, err.ChartDir)
}

// ChartNotFoundError is returned when a provided chart dir is not found.
type ChartNotFoundError struct {
	Path string
}

// Error implements the error interface for ChartNotFoundError.
func (err ChartNotFoundError) Error() string {
	return "could not find chart path " + err.Path
}
