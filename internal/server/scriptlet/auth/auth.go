package auth

import (
	"context"
	"fmt"

	"go.starlark.net/starlark"

	"github.com/lxc/incus/v6/internal/server/auth/common"
	scriptletLoad "github.com/lxc/incus/v6/internal/server/scriptlet/load"
	"github.com/lxc/incus/v6/internal/server/scriptlet/log"
	"github.com/lxc/incus/v6/internal/server/scriptlet/marshal"
	"github.com/lxc/incus/v6/shared/logger"
)

// AuthorizationRun runs the authorization scriptlet.
func AuthorizationRun(l logger.Logger, details *common.RequestDetails, object string, entitlement string) (bool, error) {
	logFunc := log.CreateLogger(l, "Authorization scriptlet")

	// Remember to match the entries in scriptletLoad.AuthorizationCompile() with this list so Starlark can
	// perform compile time validation of functions used.
	env := starlark.StringDict{
		"log_info":  starlark.NewBuiltin("log_info", logFunc),
		"log_warn":  starlark.NewBuiltin("log_warn", logFunc),
		"log_error": starlark.NewBuiltin("log_error", logFunc),
	}

	prog, thread, err := scriptletLoad.AuthorizationProgram()
	if err != nil {
		return false, err
	}

	globals, err := prog.Init(thread, env)
	if err != nil {
		return false, fmt.Errorf("Failed initializing: %w", err)
	}

	globals.Freeze()

	// Retrieve a global variable from starlark environment.
	authorizer := globals["authorize"]
	if authorizer == nil {
		return false, fmt.Errorf("Scriptlet missing authorize function")
	}

	detailsv, err := marshal.StarlarkMarshal(details)
	if err != nil {
		return false, fmt.Errorf("Marshalling details failed: %w", err)
	}

	// Call starlark function from Go.
	v, err := starlark.Call(thread, authorizer, nil, []starlark.Tuple{
		{
			starlark.String("details"),
			detailsv,
		}, {
			starlark.String("object"),
			starlark.String(object),
		}, {
			starlark.String("entitlement"),
			starlark.String(entitlement),
		},
	})
	if err != nil {
		return false, fmt.Errorf("Failed to run: %w", err)
	}

	if v.Type() != "bool" {
		return false, fmt.Errorf("Failed with unexpected return value: %v", v)
	}

	return bool(v.(starlark.Bool)), nil
}
