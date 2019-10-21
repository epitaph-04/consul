package agent

func init() {
	registerEndpoint("/v1/acl/bootstrap", []string{"PUT"}, (*HTTPServer).ACLBootstrap)
	registerEndpoint("/v1/acl/create", []string{"PUT"}, (*HTTPServer).ACLCreate)
	registerEndpoint("/v1/acl/update", []string{"PUT"}, (*HTTPServer).ACLUpdate)
	registerEndpoint("/v1/acl/destroy/", []string{"PUT"}, (*HTTPServer).ACLDestroy)
	registerEndpoint("/v1/acl/info/", []string{"GET"}, (*HTTPServer).ACLGet)
	registerEndpoint("/v1/acl/clone/", []string{"PUT"}, (*HTTPServer).ACLClone)
	registerEndpoint("/v1/acl/list", []string{"GET"}, (*HTTPServer).ACLList)
	registerEndpoint("/v1/acl/login", []string{"POST"}, (*HTTPServer).ACLLogin)
	registerEndpoint("/v1/acl/logout", []string{"POST"}, (*HTTPServer).ACLLogout)
	registerEndpoint("/v1/acl/replication", []string{"GET"}, (*HTTPServer).ACLReplicationStatus)
	registerEndpoint("/v1/acl/policies", []string{"GET"}, (*HTTPServer).ACLPolicyList)
	registerEndpoint("/v1/acl/policy", []string{"PUT"}, (*HTTPServer).ACLPolicyCreate)
	registerEndpoint("/v1/acl/policy/", []string{"GET", "PUT", "DELETE"}, (*HTTPServer).ACLPolicyCRUD)
	registerEndpoint("/v1/acl/roles", []string{"GET"}, (*HTTPServer).ACLRoleList)
	registerEndpoint("/v1/acl/role", []string{"PUT"}, (*HTTPServer).ACLRoleCreate)
	registerEndpoint("/v1/acl/role/name/", []string{"GET"}, (*HTTPServer).ACLRoleReadByName)
	registerEndpoint("/v1/acl/role/", []string{"GET", "PUT", "DELETE"}, (*HTTPServer).ACLRoleCRUD)
	registerEndpoint("/v1/acl/binding-rules", []string{"GET"}, (*HTTPServer).ACLBindingRuleList)
	registerEndpoint("/v1/acl/binding-rule", []string{"PUT"}, (*HTTPServer).ACLBindingRuleCreate)
	registerEndpoint("/v1/acl/binding-rule/", []string{"GET", "PUT", "DELETE"}, (*HTTPServer).ACLBindingRuleCRUD)
	registerEndpoint("/v1/acl/auth-methods", []string{"GET"}, (*HTTPServer).ACLAuthMethodList)
	registerEndpoint("/v1/acl/auth-method", []string{"PUT"}, (*HTTPServer).ACLAuthMethodCreate)
	registerEndpoint("/v1/acl/auth-method/", []string{"GET", "PUT", "DELETE"}, (*HTTPServer).ACLAuthMethodCRUD)
	registerEndpoint("/v1/acl/rules/translate", []string{"POST"}, (*HTTPServer).ACLRulesTranslate)
	registerEndpoint("/v1/acl/rules/translate/", []string{"GET"}, (*HTTPServer).ACLRulesTranslateLegacyToken)
	registerEndpoint("/v1/acl/tokens", []string{"GET"}, (*HTTPServer).ACLTokenList)
	registerEndpoint("/v1/acl/token", []string{"PUT"}, (*HTTPServer).ACLTokenCreate)
	registerEndpoint("/v1/acl/token/self", []string{"GET"}, (*HTTPServer).ACLTokenSelf)
	registerEndpoint("/v1/acl/token/", []string{"GET", "PUT", "DELETE"}, (*HTTPServer).ACLTokenCRUD)
	registerEndpoint("/v1/agent/token/", []string{"PUT"}, (*HTTPServer).AgentToken)
	registerEndpoint("/v1/agent/self", []string{"GET"}, (*HTTPServer).AgentSelf)
	registerEndpoint("/v1/agent/host", []string{"GET"}, (*HTTPServer).AgentHost)
	registerEndpoint("/v1/agent/maintenance", []string{"PUT"}, (*HTTPServer).AgentNodeMaintenance)
	registerEndpoint("/v1/agent/reload", []string{"PUT"}, (*HTTPServer).AgentReload)
	registerEndpoint("/v1/agent/monitor", []string{"GET"}, (*HTTPServer).AgentMonitor)
	registerEndpoint("/v1/agent/metrics", []string{"GET"}, (*HTTPServer).AgentMetrics)
	registerEndpoint("/v1/agent/services", []string{"GET"}, (*HTTPServer).AgentServices)
	registerEndpoint("/v1/agent/service/", []string{"GET"}, (*HTTPServer).AgentService)
	registerEndpoint("/v1/agent/checks", []string{"GET"}, (*HTTPServer).AgentChecks)
	registerEndpoint("/v1/agent/members", []string{"GET"}, (*HTTPServer).AgentMembers)
	registerEndpoint("/v1/agent/join/", []string{"PUT"}, (*HTTPServer).AgentJoin)
	registerEndpoint("/v1/agent/leave", []string{"PUT"}, (*HTTPServer).AgentLeave)
	registerEndpoint("/v1/agent/force-leave/", []string{"PUT"}, (*HTTPServer).AgentForceLeave)
	registerEndpoint("/v1/agent/health/service/id/", []string{"GET"}, (*HTTPServer).AgentHealthServiceByID)
	registerEndpoint("/v1/agent/health/service/name/", []string{"GET"}, (*HTTPServer).AgentHealthServiceByName)
	registerEndpoint("/v1/agent/check/register", []string{"PUT"}, (*HTTPServer).AgentRegisterCheck)
	registerEndpoint("/v1/agent/check/deregister/", []string{"PUT"}, (*HTTPServer).AgentDeregisterCheck)
	registerEndpoint("/v1/agent/check/pass/", []string{"PUT"}, (*HTTPServer).AgentCheckPass)
	registerEndpoint("/v1/agent/check/warn/", []string{"PUT"}, (*HTTPServer).AgentCheckWarn)
	registerEndpoint("/v1/agent/check/fail/", []string{"PUT"}, (*HTTPServer).AgentCheckFail)
	registerEndpoint("/v1/agent/check/update/", []string{"PUT"}, (*HTTPServer).AgentCheckUpdate)
	registerEndpoint("/v1/agent/connect/authorize", []string{"POST"}, (*HTTPServer).AgentConnectAuthorize)
	registerEndpoint("/v1/agent/connect/ca/roots", []string{"GET"}, (*HTTPServer).AgentConnectCARoots)
	registerEndpoint("/v1/agent/connect/ca/leaf/", []string{"GET"}, (*HTTPServer).AgentConnectCALeafCert)
	registerEndpoint("/v1/agent/service/register", []string{"PUT"}, (*HTTPServer).AgentRegisterService)
	registerEndpoint("/v1/agent/service/deregister/", []string{"PUT"}, (*HTTPServer).AgentDeregisterService)
	registerEndpoint("/v1/agent/service/maintenance/", []string{"PUT"}, (*HTTPServer).AgentServiceMaintenance)
	registerEndpoint("/v1/catalog/register", []string{"PUT"}, (*HTTPServer).CatalogRegister)
	registerEndpoint("/v1/catalog/connect/", []string{"GET"}, (*HTTPServer).CatalogConnectServiceNodes)
	registerEndpoint("/v1/catalog/deregister", []string{"PUT"}, (*HTTPServer).CatalogDeregister)
	registerEndpoint("/v1/catalog/datacenters", []string{"GET"}, (*HTTPServer).CatalogDatacenters)
	registerEndpoint("/v1/catalog/nodes", []string{"GET"}, (*HTTPServer).CatalogNodes)
	registerEndpoint("/v1/catalog/services", []string{"GET"}, (*HTTPServer).CatalogServices)
	registerEndpoint("/v1/catalog/service/", []string{"GET"}, (*HTTPServer).CatalogServiceNodes)
	registerEndpoint("/v1/catalog/node/", []string{"GET"}, (*HTTPServer).CatalogNodeServices)
	registerEndpoint("/v1/config/", []string{"GET", "DELETE"}, (*HTTPServer).Config)
	registerEndpoint("/v1/config", []string{"PUT"}, (*HTTPServer).ConfigApply)
	registerEndpoint("/v1/connect/ca/configuration", []string{"GET", "PUT"}, (*HTTPServer).ConnectCAConfiguration)
	registerEndpoint("/v1/connect/ca/roots", []string{"GET"}, (*HTTPServer).ConnectCARoots)
	registerEndpoint("/v1/connect/intentions", []string{"GET", "POST"}, (*HTTPServer).IntentionEndpoint)
	registerEndpoint("/v1/connect/intentions/match", []string{"GET"}, (*HTTPServer).IntentionMatch)
	registerEndpoint("/v1/connect/intentions/check", []string{"GET"}, (*HTTPServer).IntentionCheck)
	registerEndpoint("/v1/connect/intentions/", []string{"GET", "PUT", "DELETE"}, (*HTTPServer).IntentionSpecific)
	registerEndpoint("/v1/coordinate/datacenters", []string{"GET"}, (*HTTPServer).CoordinateDatacenters)
	registerEndpoint("/v1/coordinate/nodes", []string{"GET"}, (*HTTPServer).CoordinateNodes)
	registerEndpoint("/v1/coordinate/node/", []string{"GET"}, (*HTTPServer).CoordinateNode)
	registerEndpoint("/v1/coordinate/update", []string{"PUT"}, (*HTTPServer).CoordinateUpdate)
	registerEndpoint("/v1/datacenter-configs", []string{"GET"}, (*HTTPServer).DatacenterConfigList)
	registerEndpoint("/v1/datacenter-config/", []string{"GET", "PUT", "DELETE"}, (*HTTPServer).DatacenterConfigCRUD)
	registerEndpoint("/v1/datacenter-config", []string{"PUT"}, (*HTTPServer).DatacenterConfigCreate)
	registerEndpoint("/v1/discovery-chain/", []string{"GET", "POST"}, (*HTTPServer).DiscoveryChainRead)
	registerEndpoint("/v1/event/fire/", []string{"PUT"}, (*HTTPServer).EventFire)
	registerEndpoint("/v1/event/list", []string{"GET"}, (*HTTPServer).EventList)
	registerEndpoint("/v1/health/node/", []string{"GET"}, (*HTTPServer).HealthNodeChecks)
	registerEndpoint("/v1/health/checks/", []string{"GET"}, (*HTTPServer).HealthServiceChecks)
	registerEndpoint("/v1/health/state/", []string{"GET"}, (*HTTPServer).HealthChecksInState)
	registerEndpoint("/v1/health/service/", []string{"GET"}, (*HTTPServer).HealthServiceNodes)
	registerEndpoint("/v1/health/connect/", []string{"GET"}, (*HTTPServer).HealthConnectServiceNodes)
	registerEndpoint("/v1/internal/ui/nodes", []string{"GET"}, (*HTTPServer).UINodes)
	registerEndpoint("/v1/internal/ui/node/", []string{"GET"}, (*HTTPServer).UINodeInfo)
	registerEndpoint("/v1/internal/ui/services", []string{"GET"}, (*HTTPServer).UIServices)
	registerEndpoint("/v1/kv/", []string{"GET", "PUT", "DELETE"}, (*HTTPServer).KVSEndpoint)
	registerEndpoint("/v1/operator/raft/configuration", []string{"GET"}, (*HTTPServer).OperatorRaftConfiguration)
	registerEndpoint("/v1/operator/raft/peer", []string{"DELETE"}, (*HTTPServer).OperatorRaftPeer)
	registerEndpoint("/v1/operator/keyring", []string{"GET", "POST", "PUT", "DELETE"}, (*HTTPServer).OperatorKeyringEndpoint)
	registerEndpoint("/v1/operator/autopilot/configuration", []string{"GET", "PUT"}, (*HTTPServer).OperatorAutopilotConfiguration)
	registerEndpoint("/v1/operator/autopilot/health", []string{"GET"}, (*HTTPServer).OperatorServerHealth)
	registerEndpoint("/v1/query", []string{"GET", "POST"}, (*HTTPServer).PreparedQueryGeneral)
	// specific prepared query endpoints have more complex rules for allowed methods, so
	// the prefix is registered with no methods.
	registerEndpoint("/v1/query/", []string{}, (*HTTPServer).PreparedQuerySpecific)
	registerEndpoint("/v1/session/create", []string{"PUT"}, (*HTTPServer).SessionCreate)
	registerEndpoint("/v1/session/destroy/", []string{"PUT"}, (*HTTPServer).SessionDestroy)
	registerEndpoint("/v1/session/renew/", []string{"PUT"}, (*HTTPServer).SessionRenew)
	registerEndpoint("/v1/session/info/", []string{"GET"}, (*HTTPServer).SessionGet)
	registerEndpoint("/v1/session/node/", []string{"GET"}, (*HTTPServer).SessionsForNode)
	registerEndpoint("/v1/session/list", []string{"GET"}, (*HTTPServer).SessionList)
	registerEndpoint("/v1/status/leader", []string{"GET"}, (*HTTPServer).StatusLeader)
	registerEndpoint("/v1/status/peers", []string{"GET"}, (*HTTPServer).StatusPeers)
	registerEndpoint("/v1/snapshot", []string{"GET", "PUT"}, (*HTTPServer).Snapshot)
	registerEndpoint("/v1/txn", []string{"PUT"}, (*HTTPServer).Txn)
}
