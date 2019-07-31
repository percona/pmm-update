// pmm-update
// Copyright (C) 2019 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package yum

import (
	"strings"

	"github.com/pkg/errors"
)

// parseInfo parses `yum info` stdout for a single version of a single package.
func parseInfo(b []byte) (map[string]string, error) {
	res := make(map[string]string)
	var prevKey string
	var nameFound bool
	for _, line := range strings.Split(string(b), "\n") {
		// separate progress output from data
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key == "Name" {
			// sanity check that we do not try to parse multiple packages
			if nameFound {
				return res, errors.New("second `Name` encountered")
			}
			nameFound = true
		}
		if !nameFound {
			continue
		}

		// parse data while handling multiline values
		value := strings.TrimSpace(parts[1])
		if key == "" {
			if prevKey != "" {
				res[prevKey] += " " + value
			}
			continue
		}
		res[key] = value
		prevKey = key
	}
	return res, nil
}
