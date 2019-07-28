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
	"context"
	"strings"
	"time"

	"github.com/percona/pmm/version"
	"github.com/pkg/errors"

	"github.com/percona/pmm-update/pkg/run"
)

// CheckVersions returns up-to-date versions information for a package with given name.
func CheckVersions(ctx context.Context, name string) (*version.UpdateCheckResult, error) {
	// http://man7.org/linux/man-pages/man8/yum.8.html#LIST_OPTIONS

	cmdLine := "yum --showduplicates list all " + name
	stdout, _, err := run.Run(ctx, 30*time.Second, cmdLine)
	if err != nil {
		return nil, errors.Wrapf(err, "%#q failed", cmdLine)
	}

	var res version.UpdateCheckResult
	for _, line := range stdout {
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) != 3 {
			continue
		}
		pack, ver, repo := parts[0], parts[1], parts[2]

		if !strings.HasPrefix(pack, name+".") {
			continue
		}
		if strings.HasPrefix(repo, "@") {
			if res.InstalledRPMVersion != "" {
				return nil, errors.New("failed to parse `yum list` output")
			}
			res.InstalledRPMVersion = ver
		} else {
			// always overwrite - the last item is the one we need
			res.LatestRPMVersion = ver
			res.LatestRepo = repo
		}
	}

	return &res, nil
}

// UpdatePackage updates package with given name.
func UpdatePackage(ctx context.Context, name string) error {
	cmdLine := "yum update " + name
	_, _, err := run.Run(ctx, 30*time.Second, cmdLine)
	if err != nil {
		return errors.Wrapf(err, "%#q failed", cmdLine)
	}
	return nil
}
