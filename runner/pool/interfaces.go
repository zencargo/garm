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

package pool

import (
	"context"

	"github.com/google/go-github/v57/github"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/params"
)

type poolHelper interface {
	GetGithubRunners() ([]*github.Runner, error)
	GetGithubRegistrationToken() (string, error)
	GetRunnerInfoFromWorkflow(job params.WorkflowJob) (params.RunnerInfo, error)
	RemoveGithubRunner(runnerID int64) (*github.Response, error)
	FetchTools() ([]commonParams.RunnerApplicationDownload, error)

	InstallHook(ctx context.Context, req *github.Hook) (params.HookInfo, error)
	UninstallHook(ctx context.Context, url string) error
	GetHookInfo(ctx context.Context) (params.HookInfo, error)

	GetJITConfig(ctx context.Context, instanceName string, pool params.Pool, labels []string) (map[string]string, *github.Runner, error)

	FetchDbInstances() ([]params.Instance, error)
	ListPools() ([]params.Pool, error)
	GithubURL() string
	JwtToken() string
	String() string
	GetPoolByID(poolID string) (params.Pool, error)
	ValidateOwner(job params.WorkflowJob) error
	UpdateState(param params.UpdatePoolStateParams) error
	WebhookSecret() string
	ID() string
	PoolBalancerType() params.PoolBalancerType
}
