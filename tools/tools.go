//go:build tools
// +build tools

package tools

import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"

	// testing
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
)
