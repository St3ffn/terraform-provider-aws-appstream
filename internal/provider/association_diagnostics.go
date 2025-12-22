// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

type associateDiagnosticMode string

const (
	associateDiagnosticPlan   associateDiagnosticMode = "plan"
	associateDiagnosticRead   associateDiagnosticMode = "read"
	associateDiagnosticDelete associateDiagnosticMode = "delete"
)
