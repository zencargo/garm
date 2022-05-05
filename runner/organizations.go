// Copyright 2022 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package runner

import (
	"context"
	"log"
	"strings"

	"garm/auth"
	runnerErrors "garm/errors"
	"garm/params"
	"garm/runner/common"
	"garm/runner/pool"

	"github.com/pkg/errors"
)

func (r *Runner) CreateOrganization(ctx context.Context, param params.CreateOrgParams) (org params.Organization, err error) {
	if !auth.IsAdmin(ctx) {
		return org, runnerErrors.ErrUnauthorized
	}

	if err := param.Validate(); err != nil {
		return params.Organization{}, errors.Wrap(err, "validating params")
	}

	creds, ok := r.credentials[param.CredentialsName]
	if !ok {
		return params.Organization{}, runnerErrors.NewBadRequestError("credentials %s not defined", param.CredentialsName)
	}

	_, err = r.store.GetOrganization(ctx, param.Name)
	if err != nil {
		if !errors.Is(err, runnerErrors.ErrNotFound) {
			return params.Organization{}, errors.Wrap(err, "fetching repo")
		}
	} else {
		return params.Organization{}, runnerErrors.NewConflictError("organization %s already exists", param.Name)
	}

	org, err = r.store.CreateOrganization(ctx, param.Name, creds.Name, param.WebhookSecret)
	if err != nil {
		return params.Organization{}, errors.Wrap(err, "creating organization")
	}

	defer func() {
		if err != nil {
			r.store.DeleteOrganization(ctx, org.ID)
		}
	}()

	poolMgr, err := r.loadOrgPoolManager(org)
	if err := poolMgr.Start(); err != nil {
		return params.Organization{}, errors.Wrap(err, "starting pool manager")
	}
	r.organizations[org.ID] = poolMgr
	return org, nil
}

func (r *Runner) ListOrganizations(ctx context.Context) ([]params.Organization, error) {
	if !auth.IsAdmin(ctx) {
		return nil, runnerErrors.ErrUnauthorized
	}

	orgs, err := r.store.ListOrganizations(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing organizations")
	}

	return orgs, nil
}

func (r *Runner) GetOrganizationByID(ctx context.Context, orgID string) (params.Organization, error) {
	if !auth.IsAdmin(ctx) {
		return params.Organization{}, runnerErrors.ErrUnauthorized
	}

	org, err := r.store.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return params.Organization{}, errors.Wrap(err, "fetching repository")
	}
	return org, nil
}

func (r *Runner) DeleteOrganization(ctx context.Context, orgID string) error {
	if !auth.IsAdmin(ctx) {
		return runnerErrors.ErrUnauthorized
	}

	org, err := r.store.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return errors.Wrap(err, "fetching repo")
	}

	poolMgr, ok := r.organizations[org.ID]
	if ok {
		if err := poolMgr.Stop(); err != nil {
			log.Printf("failed to stop pool for repo %s", org.ID)
		}
		delete(r.organizations, orgID)
	}

	pools, err := r.store.ListOrgPools(ctx, orgID)
	if err != nil {
		return errors.Wrap(err, "fetching repo pools")
	}

	if len(pools) > 0 {
		poolIds := []string{}
		for _, pool := range pools {
			poolIds = append(poolIds, pool.ID)
		}

		return runnerErrors.NewBadRequestError("repo has pools defined (%s)", strings.Join(poolIds, ", "))
	}

	if err := r.store.DeleteOrganization(ctx, orgID); err != nil {
		return errors.Wrapf(err, "removing organization %s", orgID)
	}
	return nil
}

func (r *Runner) UpdateOrganization(ctx context.Context, orgID string, param params.UpdateRepositoryParams) (params.Organization, error) {
	if !auth.IsAdmin(ctx) {
		return params.Organization{}, runnerErrors.ErrUnauthorized
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	org, err := r.store.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return params.Organization{}, errors.Wrap(err, "fetching org")
	}

	if param.CredentialsName != "" {
		// Check that credentials are set before saving to db
		if _, ok := r.credentials[param.CredentialsName]; !ok {
			return params.Organization{}, runnerErrors.NewBadRequestError("invalid credentials (%s) for org %s", param.CredentialsName, org.Name)
		}
	}

	org, err = r.store.UpdateOrganization(ctx, orgID, param)
	if err != nil {
		return params.Organization{}, errors.Wrap(err, "updating org")
	}

	poolMgr, ok := r.organizations[org.ID]
	if ok {
		internalCfg, err := r.getInternalConfig(org.CredentialsName)
		if err != nil {
			return params.Organization{}, errors.Wrap(err, "fetching internal config")
		}
		newState := params.UpdatePoolStateParams{
			WebhookSecret: org.WebhookSecret,
			Internal:      internalCfg,
		}
		org.Internal = internalCfg
		// stop the pool mgr
		if err := poolMgr.RefreshState(newState); err != nil {
			return params.Organization{}, errors.Wrap(err, "updating pool manager")
		}
	} else {
		poolMgr, err := r.loadOrgPoolManager(org)
		if err != nil {
			return params.Organization{}, errors.Wrap(err, "loading pool manager")
		}
		r.organizations[org.ID] = poolMgr
	}

	return org, nil
}

func (r *Runner) CreateOrgPool(ctx context.Context, orgID string, param params.CreatePoolParams) (params.Pool, error) {
	if !auth.IsAdmin(ctx) {
		return params.Pool{}, runnerErrors.ErrUnauthorized
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.organizations[orgID]
	if !ok {
		return params.Pool{}, runnerErrors.ErrNotFound
	}

	createPoolParams, err := r.appendTagsToCreatePoolParams(param)
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "fetching pool params")
	}

	pool, err := r.store.CreateOrganizationPool(ctx, orgID, createPoolParams)
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "creating pool")
	}

	return pool, nil
}

func (r *Runner) GetOrgPoolByID(ctx context.Context, orgID, poolID string) (params.Pool, error) {
	if !auth.IsAdmin(ctx) {
		return params.Pool{}, runnerErrors.ErrUnauthorized
	}

	pool, err := r.store.GetOrganizationPool(ctx, orgID, poolID)
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "fetching pool")
	}
	return pool, nil
}

func (r *Runner) DeleteOrgPool(ctx context.Context, orgID, poolID string) error {
	if !auth.IsAdmin(ctx) {
		return runnerErrors.ErrUnauthorized
	}

	// TODO: dedup instance count verification
	pool, err := r.store.GetOrganizationPool(ctx, orgID, poolID)
	if err != nil {
		return errors.Wrap(err, "fetching pool")
	}

	instances, err := r.store.ListPoolInstances(ctx, pool.ID)
	if err != nil {
		return errors.Wrap(err, "fetching instances")
	}

	// TODO: implement a count function
	if len(instances) > 0 {
		runnerIDs := []string{}
		for _, run := range instances {
			runnerIDs = append(runnerIDs, run.ID)
		}
		return runnerErrors.NewBadRequestError("pool has runners: %s", strings.Join(runnerIDs, ", "))
	}

	if err := r.store.DeleteOrganizationPool(ctx, orgID, poolID); err != nil {
		return errors.Wrap(err, "deleting pool")
	}
	return nil
}

func (r *Runner) ListOrgPools(ctx context.Context, orgID string) ([]params.Pool, error) {
	if !auth.IsAdmin(ctx) {
		return []params.Pool{}, runnerErrors.ErrUnauthorized
	}

	pools, err := r.store.ListOrgPools(ctx, orgID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching pools")
	}
	return pools, nil
}

func (r *Runner) UpdateOrgPool(ctx context.Context, orgID, poolID string, param params.UpdatePoolParams) (params.Pool, error) {
	if !auth.IsAdmin(ctx) {
		return params.Pool{}, runnerErrors.ErrUnauthorized
	}

	pool, err := r.store.GetOrganizationPool(ctx, orgID, poolID)
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "fetching pool")
	}

	maxRunners := pool.MaxRunners
	minIdleRunners := pool.MinIdleRunners

	if param.MaxRunners != nil {
		maxRunners = *param.MaxRunners
	}
	if param.MinIdleRunners != nil {
		minIdleRunners = *param.MinIdleRunners
	}

	if minIdleRunners > maxRunners {
		return params.Pool{}, runnerErrors.NewBadRequestError("min_idle_runners cannot be larger than max_runners")
	}

	newPool, err := r.store.UpdateOrganizationPool(ctx, orgID, poolID, param)
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "updating pool")
	}
	return newPool, nil
}

func (r *Runner) ListOrgInstances(ctx context.Context, orgID string) ([]params.Instance, error) {
	if !auth.IsAdmin(ctx) {
		return nil, runnerErrors.ErrUnauthorized
	}

	instances, err := r.store.ListOrgInstances(ctx, orgID)
	if err != nil {
		return []params.Instance{}, errors.Wrap(err, "fetching instances")
	}
	return instances, nil
}

func (r *Runner) loadOrgPoolManager(org params.Organization) (common.PoolManager, error) {
	cfg, err := r.getInternalConfig(org.CredentialsName)
	if err != nil {
		return nil, errors.Wrap(err, "fetching internal config")
	}
	org.Internal = cfg
	poolManager, err := pool.NewOrganizationPoolManager(r.ctx, org, r.providers, r.store)
	if err != nil {
		return nil, errors.Wrap(err, "creating pool manager")
	}
	return poolManager, nil
}

func (r *Runner) findOrgPoolManager(name string) (common.PoolManager, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	org, err := r.store.GetOrganization(r.ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "fetching repo")
	}

	if orgPoolMgr, ok := r.organizations[org.ID]; ok {
		return orgPoolMgr, nil
	}
	return nil, errors.Wrapf(runnerErrors.ErrNotFound, "organization %s not configured", name)
}
