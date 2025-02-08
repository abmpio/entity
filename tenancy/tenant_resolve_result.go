package tenancy

// tenant resolve result
type TenantResolveResult struct {
	// TenantId
	TenantId         string
	AppliedResolvers []string
}
