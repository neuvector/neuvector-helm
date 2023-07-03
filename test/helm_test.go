package test

import (
	"strings"
)

const nvRel = "nv"

func splitYaml(out string) []string {
	outputs := make([]string, 0)

	outs := strings.Split(out, "---")
	for _, out := range outs {
		out := strings.TrimSpace(out)

		// split section into lines, if all lines are empty or comment, ignore the section
		lines := strings.Split(out, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") || len(line) == 0 {
				continue
			}

			outputs = append(outputs, out)
			break
		}
	}
	return outputs
}
