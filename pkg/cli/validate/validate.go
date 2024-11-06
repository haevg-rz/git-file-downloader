package validate

import (
	"errors"
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"strings"
)

func Flags(flagToValue map[string]interface{}) error {
	missingFlags := make([]string, 0)

	// extend value switch according to requirements
	for flag, value := range flagToValue {
		missing := false

		switch value.(type) {
		case string:
			if value.(string) == "" {
				missing = true
			}
		case int:
			if value.(int) == -1 {
				missing = true
			}
		}

		if missing {
			missingFlags = append(missingFlags, flag)
		}
	}

	if len(missingFlags) > 0 {
		exit.Code = exit.MissingFlags
		return errors.New(fmt.Sprintf("missing flags: %s\n", strings.Join(missingFlags, ", ")))
	}
	return nil
}
