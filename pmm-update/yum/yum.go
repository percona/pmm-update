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
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
)

func run(ctx context.Context, cmdLine string) ([]string, error) {
	// TODO graceful cancelation with ctx

	args := strings.Fields(cmdLine)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	setSysProcAttr(cmd)
	var stdout bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.WithStack(err)
	}
	return strings.Split(stdout.String(), "\n"), nil
}

func CheckVersions(ctx context.Context) (map[string]version.Info, error) {
	stdout, err := run(ctx, "yum list --showduplicates pmm-update")
	if err != nil {
		return nil, err
	}

	res := make(map[string]version.Info)
	for _, line := range stdout {
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) != 3 {
			continue
		}
		pack, ver, repo := parts[0], parts[1], parts[2]

		if pack != "pmm-update.noarch" {
			continue
		}
		if repo == "@local" || strings.HasPrefix(repo, "pmm2-") {
			v, err := version.Parse(ver)
			if err != nil {
				continue
			}
			res[repo] = v
		}
	}
	return res, nil
}
