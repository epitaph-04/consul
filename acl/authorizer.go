package acl

import (
	"fmt"
	"strings"
)

type EnforcementDecision int

const (
	// Deny returned from an Authorizer enforcement method indicates
	// that a corresponding rule was found and that access should be denied
	Deny EnforcementDecision = iota
	// Allow returned from an Authorizer enforcement method indicates
	// that a corresponding rule was found and that access should be allowed
	Allow
	// Default returned from an Authorizer enforcement method indicates
	// that a corresponding rule was not found and that whether access
	// should be granted or denied should be deferred to the default
	// access level
	Default
)

func (d EnforcementDecision) String() string {
	switch d {
	case Allow:
		return "Allow"
	case Deny:
		return "Deny"
	case Default:
		return "Default"
	default:
		return "Unknown"
	}
}

type Resource string

const (
	ResourceACL       Resource = "acl"
	ResourceAgent     Resource = "agent"
	ResourceEvent     Resource = "event"
	ResourceIntention Resource = "intention"
	ResourceKey       Resource = "key"
	ResourceKeyring   Resource = "keyring"
	ResourceNode      Resource = "node"
	ResourceOperator  Resource = "operator"
	ResourceQuery     Resource = "query"
	ResourceService   Resource = "service"
	ResourceSession   Resource = "session"
)

// Authorizer is the interface for policy enforcement.
type Authorizer interface {
	// ACLRead checks for permission to list all the ACLs
	ACLRead(*EnterpriseAuthorizerContext) EnforcementDecision

	// ACLWrite checks for permission to manipulate ACLs
	ACLWrite(*EnterpriseAuthorizerContext) EnforcementDecision

	// AgentRead checks for permission to read from agent endpoints for a
	// given node.
	AgentRead(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// AgentWrite checks for permission to make changes via agent endpoints
	// for a given node.
	AgentWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// EventRead determines if a specific event can be queried.
	EventRead(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// EventWrite determines if a specific event may be fired.
	EventWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// IntentionDefaultAllow determines the default authorized behavior
	// when no intentions match a Connect request.
	IntentionDefaultAllow(*EnterpriseAuthorizerContext) EnforcementDecision

	// IntentionRead determines if a specific intention can be read.
	IntentionRead(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// IntentionWrite determines if a specific intention can be
	// created, modified, or deleted.
	IntentionWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// KeyList checks for permission to list keys under a prefix
	KeyList(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// KeyRead checks for permission to read a given key
	KeyRead(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// KeyWrite checks for permission to write a given key
	KeyWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// KeyWritePrefix checks for permission to write to an
	// entire key prefix. This means there must be no sub-policies
	// that deny a write.
	KeyWritePrefix(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// KeyringRead determines if the encryption keyring used in
	// the gossip layer can be read.
	KeyringRead(*EnterpriseAuthorizerContext) EnforcementDecision

	// KeyringWrite determines if the keyring can be manipulated
	KeyringWrite(*EnterpriseAuthorizerContext) EnforcementDecision

	// NodeRead checks for permission to read (discover) a given node.
	NodeRead(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// NodeWrite checks for permission to create or update (register) a
	// given node.
	NodeWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// OperatorRead determines if the read-only Consul operator functions
	// can be used.
	OperatorRead(*EnterpriseAuthorizerContext) EnforcementDecision

	// OperatorWrite determines if the state-changing Consul operator
	// functions can be used.
	OperatorWrite(*EnterpriseAuthorizerContext) EnforcementDecision

	// PreparedQueryRead determines if a specific prepared query can be read
	// to show its contents (this is not used for execution).
	PreparedQueryRead(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// PreparedQueryWrite determines if a specific prepared query can be
	// created, modified, or deleted.
	PreparedQueryWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// ServiceRead checks for permission to read a given service
	ServiceRead(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// ServiceWrite checks for permission to create or update a given
	// service
	ServiceWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// SessionRead checks for permission to read sessions for a given node.
	SessionRead(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// SessionWrite checks for permission to create sessions for a given
	// node.
	SessionWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision

	// Snapshot checks for permission to take and restore snapshots.
	Snapshot(*EnterpriseAuthorizerContext) EnforcementDecision

	// Embedded Interface for Consul Enterprise specific ACL enforcement
	EnterpriseAuthorizer
}

func Enforce(authz Authorizer, rsc Resource, segment string, access string, ctx *EnterpriseAuthorizerContext) (EnforcementDecision, error) {
	lowerAccess := strings.ToLower(access)

	switch rsc {
	case ResourceACL:
		switch lowerAccess {
		case "read":
			return authz.ACLRead(ctx), nil
		case "write":
			return authz.ACLWrite(ctx), nil
		}
	case ResourceAgent:
		switch lowerAccess {
		case "read":
			return authz.AgentRead(segment, ctx), nil
		case "write":
			return authz.AgentWrite(segment, ctx), nil
		}
	case ResourceEvent:
		switch lowerAccess {
		case "read":
			return authz.EventRead(segment, ctx), nil
		case "write":
			return authz.EventWrite(segment, ctx), nil
		}
	case ResourceIntention:
		switch lowerAccess {
		case "read":
			return authz.IntentionRead(segment, ctx), nil
		case "write":
			return authz.IntentionWrite(segment, ctx), nil
		}
	case ResourceKey:
		switch lowerAccess {
		case "read":
			return authz.KeyRead(segment, ctx), nil
		case "list":
			return authz.KeyList(segment, ctx), nil
		case "write":
			return authz.KeyWrite(segment, ctx), nil
		}
	case ResourceKeyring:
		switch lowerAccess {
		case "read":
			return authz.KeyringRead(ctx), nil
		case "write":
			return authz.KeyringWrite(ctx), nil
		}
	case ResourceNode:
		switch lowerAccess {
		case "read":
			return authz.NodeRead(segment, ctx), nil
		case "write":
			return authz.NodeWrite(segment, ctx), nil
		}
	case ResourceOperator:
		switch lowerAccess {
		case "read":
			return authz.OperatorRead(ctx), nil
		case "write":
			return authz.OperatorWrite(ctx), nil
		}
	case ResourceQuery:
		switch lowerAccess {
		case "read":
			return authz.PreparedQueryRead(segment, ctx), nil
		case "write":
			return authz.PreparedQueryWrite(segment, ctx), nil
		}
	case ResourceService:
		switch lowerAccess {
		case "read":
			return authz.ServiceRead(segment, ctx), nil
		case "write":
			return authz.ServiceWrite(segment, ctx), nil
		}
	case ResourceSession:
		switch lowerAccess {
		case "read":
			return authz.SessionRead(segment, ctx), nil
		case "write":
			return authz.SessionWrite(segment, ctx), nil
		}
	default:
		if processed, decision, err := EnforceEnterprise(authz, rsc, segment, access, ctx); processed {
			return decision, err
		}
		return Deny, fmt.Errorf("Invalid ACL resource requested: %q", rsc)
	}

	return Deny, fmt.Errorf("Invalid access level for %s resource: %s", rsc, access)
}
