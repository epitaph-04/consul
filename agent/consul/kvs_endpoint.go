package consul

import (
	"fmt"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/acl"
	"github.com/hashicorp/consul/agent/consul/state"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-memdb"
)

// KVS endpoint is used to manipulate the Key-Value store
type KVS struct {
	srv *Server
}

// preApply does all the verification of a KVS update that is performed BEFORE
// we submit as a Raft log entry. This includes enforcing the lock delay which
// must only be done on the leader.
func kvsPreApply(srv *Server, rule acl.Authorizer, op api.KVOp, dirEnt *structs.DirEntry) (bool, error) {
	// Verify the entry.
	if dirEnt.Key == "" && op != api.KVDeleteTree {
		return false, fmt.Errorf("Must provide key")
	}

	// Apply the ACL policy if any.
	if rule != nil {
		switch op {
		case api.KVDeleteTree:
			var authzContext acl.EnterpriseAuthorizerContext
			dirEnt.FillAuthzContext(&authzContext)

			if rule.KeyWritePrefix(dirEnt.Key, &authzContext) != acl.Allow {
				return false, acl.ErrPermissionDenied
			}

		case api.KVGet, api.KVGetTree:
			// Filtering for GETs is done on the output side.

		case api.KVCheckSession, api.KVCheckIndex:
			// These could reveal information based on the outcome
			// of the transaction, and they operate on individual
			// keys so we check them here.
			var authzContext acl.EnterpriseAuthorizerContext
			dirEnt.FillAuthzContext(&authzContext)

			if rule.KeyRead(dirEnt.Key, &authzContext) != acl.Allow {
				return false, acl.ErrPermissionDenied
			}

		default:
			var authzContext acl.EnterpriseAuthorizerContext
			dirEnt.FillAuthzContext(&authzContext)

			if rule.KeyWrite(dirEnt.Key, &authzContext) != acl.Allow {
				return false, acl.ErrPermissionDenied
			}
		}
	}

	// If this is a lock, we must check for a lock-delay. Since lock-delay
	// is based on wall-time, each peer would expire the lock-delay at a slightly
	// different time. This means the enforcement of lock-delay cannot be done
	// after the raft log is committed as it would lead to inconsistent FSMs.
	// Instead, the lock-delay must be enforced before commit. This means that
	// only the wall-time of the leader node is used, preventing any inconsistencies.
	if op == api.KVLock {
		state := srv.fsm.State()
		expires := state.KVSLockDelay(dirEnt.Key, &dirEnt.EnterpriseMeta)
		if expires.After(time.Now()) {
			srv.logger.Printf("[WARN] consul.kvs: Rejecting lock of %s due to lock-delay until %v",
				dirEnt.Key, expires)
			return false, nil
		}
	}

	return true, nil
}

// Apply is used to apply a KVS update request to the data store.
func (k *KVS) Apply(args *structs.KVSRequest, reply *bool) error {
	if done, err := k.srv.forward("KVS.Apply", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"kvs", "apply"}, time.Now())

	if err := k.srv.validateEnterpriseRequest(&args.DirEnt.EnterpriseMeta, true); err != nil {
		return err
	}

	// Perform the pre-apply checks.
	rule, err := k.srv.ResolveToken(args.Token)
	if err != nil {
		return err
	}
	ok, err := kvsPreApply(k.srv, rule, args.Op, &args.DirEnt)
	if err != nil {
		return err
	}
	if !ok {
		*reply = false
		return nil
	}

	// Apply the update.
	resp, err := k.srv.raftApply(structs.KVSRequestType, args)
	if err != nil {
		k.srv.logger.Printf("[ERR] consul.kvs: Apply failed: %v", err)
		return err
	}
	if respErr, ok := resp.(error); ok {
		return respErr
	}

	// Check if the return type is a bool.
	if respBool, ok := resp.(bool); ok {
		*reply = respBool
	}
	return nil
}

// Get is used to lookup a single key.
func (k *KVS) Get(args *structs.KeyRequest, reply *structs.IndexedDirEntries) error {
	if done, err := k.srv.forward("KVS.Get", args, args, reply); done {
		return err
	}
	if err := k.srv.validateEnterpriseRequest(&args.EnterpriseMeta, false); err != nil {
		return err
	}

	var entCtx acl.EnterpriseAuthorizerContext
	args.FillAuthzContext(&entCtx)

	rule, err := k.srv.ResolveToken(args.Token)
	if err != nil {
		return err
	}

	return k.srv.blockingQuery(
		&args.QueryOptions,
		&reply.QueryMeta,
		func(ws memdb.WatchSet, state *state.Store) error {
			index, ent, err := state.KVSGet(ws, args.Key, &args.EnterpriseMeta)
			if err != nil {
				return err
			}
			if rule != nil && rule.KeyRead(args.Key, &entCtx) != acl.Allow {
				return acl.ErrPermissionDenied
			}

			if ent == nil {
				// Must provide non-zero index to prevent blocking
				// Index 1 is impossible anyways (due to Raft internals)
				if index == 0 {
					reply.Index = 1
				} else {
					reply.Index = index
				}
				reply.Entries = nil
			} else {
				reply.Index = ent.ModifyIndex
				reply.Entries = structs.DirEntries{ent}
			}
			return nil
		})
}

// List is used to list all keys with a given prefix.
func (k *KVS) List(args *structs.KeyRequest, reply *structs.IndexedDirEntries) error {
	if done, err := k.srv.forward("KVS.List", args, args, reply); done {
		return err
	}
	if err := k.srv.validateEnterpriseRequest(&args.EnterpriseMeta, false); err != nil {
		return err
	}

	var entCtx acl.EnterpriseAuthorizerContext
	args.FillAuthzContext(&entCtx)

	rule, err := k.srv.ResolveToken(args.Token)
	if err != nil {
		return err
	}
	if rule != nil && k.srv.config.ACLEnableKeyListPolicy && rule.KeyList(args.Key, &entCtx) != acl.Allow {
		return acl.ErrPermissionDenied
	}

	return k.srv.blockingQuery(
		&args.QueryOptions,
		&reply.QueryMeta,
		func(ws memdb.WatchSet, state *state.Store) error {
			index, ent, err := state.KVSList(ws, args.Key, &args.EnterpriseMeta)
			if err != nil {
				return err
			}
			if rule != nil {
				ent = FilterDirEnt(rule, ent)
			}

			if len(ent) == 0 {
				// Must provide non-zero index to prevent blocking
				// Index 1 is impossible anyways (due to Raft internals)
				if index == 0 {
					reply.Index = 1
				} else {
					reply.Index = index
				}
				reply.Entries = nil
			} else {
				reply.Index = index
				reply.Entries = ent
			}
			return nil
		})
}

// ListKeys is used to list all keys with a given prefix to a separator.
// An optional separator may be specified, which can be used to slice off a part
// of the response so that only a subset of the prefix is returned. In this
// mode, the keys which are omitted are still counted in the returned index.
func (k *KVS) ListKeys(args *structs.KeyListRequest, reply *structs.IndexedKeyList) error {
	if done, err := k.srv.forward("KVS.ListKeys", args, args, reply); done {
		return err
	}
	if err := k.srv.validateEnterpriseRequest(&args.EnterpriseMeta, false); err != nil {
		return err
	}

	var entCtx acl.EnterpriseAuthorizerContext
	args.FillAuthzContext(&entCtx)

	rule, err := k.srv.ResolveToken(args.Token)
	if err != nil {
		return err
	}
	if rule != nil && k.srv.config.ACLEnableKeyListPolicy && rule.KeyList(args.Prefix, &entCtx) != acl.Allow {
		return acl.ErrPermissionDenied
	}

	return k.srv.blockingQuery(
		&args.QueryOptions,
		&reply.QueryMeta,
		func(ws memdb.WatchSet, state *state.Store) error {
			index, entries, err := state.KVSList(ws, args.Prefix, &args.EnterpriseMeta)
			if err != nil {
				return err
			}

			// Must provide non-zero index to prevent blocking
			// Index 1 is impossible anyways (due to Raft internals)
			if index == 0 {
				reply.Index = 1
			} else {
				reply.Index = index
			}

			if rule != nil {
				entries = FilterDirEnt(rule, entries)
			}

			// Collect the keys from the filtered entries
			prefixLen := len(args.Prefix)
			sepLen := len(args.Seperator)

			var keys []string
			seen := make(map[string]bool)

			for _, e := range entries {
				// Always accumulate if no separator provided
				if sepLen == 0 {
					keys = append(keys, e.Key)
					continue
				}

				// Parse and de-duplicate the returned keys based on the
				// key separator, if provided.
				after := e.Key[prefixLen:]
				sepIdx := strings.Index(after, args.Seperator)
				if sepIdx > -1 {
					key := e.Key[:prefixLen+sepIdx+sepLen]
					if ok := seen[key]; !ok {
						keys = append(keys, key)
						seen[key] = true
					}
				} else {
					keys = append(keys, e.Key)
				}
			}
			reply.Keys = keys
			return nil
		})
}
