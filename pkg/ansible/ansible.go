package ansible

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/percona/pmm-update/pkg/run"
)

// RunPlaybook runs ansible-playbook.
func RunPlaybook(ctx context.Context, playbook string, v int) error {
	cmdLine := fmt.Sprintf(
		`ansible-playbook --flush-cache --inventory='localhost,' -%s --connection=local %s`,
		strings.Repeat("v", v), playbook,
	)
	_, _, err := run.Run(ctx, 30*time.Second, cmdLine)
	return err
}
