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
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
)

func run(cmdLine string) ([]string, error) {
	args := strings.Fields(cmdLine)
	cmd := exec.Command(args[0], args[1:]...)
	setSysProcAttr(cmd)
	var stdout bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.WithStack(err)
	}
	return strings.Split(stdout.String(), "\n"), nil
}

func CheckLatestVersions() (installed, remote version.Info, err error) {
	var stdout []string
	if stdout, err = run("yum list --showduplicates pmm-update"); err != nil {
		return
	}

	for _, line := range stdout {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "pmm-update.noarch") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}
		switch {
		case parts[2] == "@local":
			installed, err = version.Parse(parts[1])
			if err != nil {
				err = errors.WithStack(err)
				return
			}
		case strings.HasPrefix(parts[2], "pmm2-"):
			r, err := version.Parse(parts[1])
			if err != nil {
				continue
			}
			remote = r
		}
	}

	return
}
