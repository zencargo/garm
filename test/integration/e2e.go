package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	client "github.com/cloudbase/garm/client"
	clientControllerInfo "github.com/cloudbase/garm/client/controller_info"
	clientCredentials "github.com/cloudbase/garm/client/credentials"
	clientFirstRun "github.com/cloudbase/garm/client/first_run"
	clientInstances "github.com/cloudbase/garm/client/instances"
	clientJobs "github.com/cloudbase/garm/client/jobs"
	clientLogin "github.com/cloudbase/garm/client/login"
	clientMetricsToken "github.com/cloudbase/garm/client/metrics_token"
	clientOrganizations "github.com/cloudbase/garm/client/organizations"
	clientPools "github.com/cloudbase/garm/client/pools"
	clientProviders "github.com/cloudbase/garm/client/providers"
	clientRepositories "github.com/cloudbase/garm/client/repositories"
	"github.com/cloudbase/garm/cmd/garm-cli/config"
	"github.com/cloudbase/garm/params"
	"github.com/go-openapi/runtime"
	openapiRuntimeClient "github.com/go-openapi/runtime/client"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

var (
	cli       *client.GarmAPI
	cfg       config.Config
	authToken runtime.ClientAuthInfoWriter

	credentialsName  = os.Getenv("CREDENTIALS_NAME")
	workflowFileName = os.Getenv("WORKFLOW_FILE_NAME")

	repoID            string
	repoPoolID        string
	repoInstanceName  string
	repoName          = os.Getenv("REPO_NAME")
	repoWebhookSecret = os.Getenv("REPO_WEBHOOK_SECRET")

	orgID            string
	orgPoolID        string
	orgInstanceName  string
	orgName          = os.Getenv("ORG_NAME")
	orgWebhookSecret = os.Getenv("ORG_WEBHOOK_SECRET")

	username = os.Getenv("GARM_USERNAME")
	password = os.Getenv("GARM_PASSWORD")
	fullName = os.Getenv("GARM_FULLNAME")
	email    = os.Getenv("GARM_EMAIL")
	name     = os.Getenv("GARM_NAME")
	baseURL  = os.Getenv("GARM_BASE_URL")
	ghtoken  = os.Getenv("GH_TOKEN")

	poolID string
)

// //////////////// //
// helper functions //
// ///////////////////
func handleError(err error) {
	if err != nil {
		panic(fmt.Sprintf("error encountered: %s", err))
	}
}

func getGithubClient() *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghtoken})
	tc := oauth2.NewClient(context.Background(), ts)

	return github.NewClient(tc)
}

func printResponse(resp interface{}) {
	b, err := json.MarshalIndent(resp, "", "  ")
	handleError(err)
	log.Println(string(b))
}

// ///////////
// Garm Init /
// ///////////
func firstRun(apiCli *client.GarmAPI, newUser params.NewUserParams) (params.User, error) {
	firstRunResponse, err := apiCli.FirstRun.FirstRun(
		clientFirstRun.NewFirstRunParams().WithBody(newUser),
		authToken)
	if err != nil {
		return params.User{}, err
	}
	return firstRunResponse.Payload, nil
}

func login(apiCli *client.GarmAPI, params params.PasswordLoginParams) (string, error) {
	loginResponse, err := apiCli.Login.Login(
		clientLogin.NewLoginParams().WithBody(params),
		authToken)
	if err != nil {
		return "", err
	}
	return loginResponse.Payload.Token, nil
}

// ////////////////////////////
// Credentials and Providers //
// ////////////////////////////
func listCredentials(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (params.Credentials, error) {
	listCredentialsResponse, err := apiCli.Credentials.ListCredentials(
		clientCredentials.NewListCredentialsParams(),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listCredentialsResponse.Payload, nil
}

func listProviders(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (params.Providers, error) {
	listProvidersResponse, err := apiCli.Providers.ListProviders(
		clientProviders.NewListProvidersParams(),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listProvidersResponse.Payload, nil
}

// ////////////////////////
// // Controller info ////
// ////////////////////////
func getControllerInfo(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (params.ControllerInfo, error) {
	controllerInfoResponse, err := apiCli.ControllerInfo.ControllerInfo(
		clientControllerInfo.NewControllerInfoParams(),
		apiAuthToken)
	if err != nil {
		return params.ControllerInfo{}, err
	}
	return controllerInfoResponse.Payload, nil
}

// ////////
// Jobs //
// ////////
func listJobs(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (params.Jobs, error) {
	listJobsResponse, err := apiCli.Jobs.ListJobs(
		clientJobs.NewListJobsParams(),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listJobsResponse.Payload, nil
}

func waitLabelledJob(label string, timeout time.Duration) *params.Job {
	var timeWaited time.Duration = 0
	var jobs params.Jobs
	var err error

	log.Printf(">>> Waiting for job with label %s", label)
	for timeWaited < timeout {
		jobs, err = listJobs(cli, authToken)
		handleError(err)
		for _, job := range jobs {
			for _, jobLabel := range job.Labels {
				if jobLabel == label {
					return &job
				}
			}
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	printResponse(jobs)
	panic(fmt.Sprintf("Failed to wait job with label %s", label))
}

func waitJobStatus(id int64, status params.JobStatus, timeout time.Duration) *params.Job {
	var timeWaited time.Duration = 0
	var job *params.Job

	log.Printf(">>> Waiting for job %d to reach status %s", id, status)
	for timeWaited < timeout {
		jobs, err := listJobs(cli, authToken)
		handleError(err)

		job = nil
		for _, j := range jobs {
			if j.ID == id {
				job = &j
				break
			}
		}

		if job == nil {
			if status == params.JobStatusCompleted {
				// The job is not found in the list. We can safely assume
				// that it is completed
				return job
			}
			// if the job is not found, and expected status is not "completed",
			// we need to error out.
			panic(fmt.Sprintf("Job %d not found, expected to be found in status %s", id, status))
		} else if job.Status == string(status) {
			return job
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	printResponse(job)
	panic(fmt.Sprintf("timeout waiting for job %d to reach status %s", id, status))
}

func waitInstanceStatus(name string, status commonParams.InstanceStatus, runnerStatus params.RunnerStatus, timeout time.Duration) *params.Instance {
	var timeWaited time.Duration = 0
	var instance *params.Instance

	log.Printf(">>> Waiting for instance %s status to reach status %s and runner status %s", name, status, runnerStatus)
	for timeWaited < timeout {
		instance, err := getInstance(cli, authToken, name)
		handleError(err)
		log.Printf(">>> Instance %s status: %s", name, instance.Status)
		if instance.Status == status && instance.RunnerStatus == runnerStatus {
			return instance
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	printResponse(instance)
	panic(fmt.Sprintf("timeout waiting for instance %s status to reach status %s and runner status %s", name, status, runnerStatus))
}

func waitInstanceToBeRemoved(name string, timeout time.Duration) {
	var timeWaited time.Duration = 0
	var instance *params.Instance

	log.Printf(">>> Waiting for instance %s to be removed", name)
	for timeWaited < timeout {
		instances, err := listInstances(cli, authToken)
		handleError(err)

		instance = nil
		for _, i := range instances {
			if i.Name == name {
				instance = &i
				break
			}
		}
		if instance == nil {
			// The instance is not found in the list. We can safely assume
			// that it is removed
			return
		}

		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	printResponse(instance)
	panic(fmt.Sprintf("Instance %s was not removed within the timeout", name))
}

func waitPoolRunningIdleInstances(poolID string, timeout time.Duration) {
	var timeWaited time.Duration = 0
	var instances params.Instances
	var err error

	pool, err := getPool(cli, authToken, poolID)
	handleError(err)

	log.Printf(">>> Waiting for pool %s to have all instances as idle running", poolID)
	for timeWaited < timeout {
		instances, err = listInstances(cli, authToken)
		handleError(err)

		poolInstances := make(params.Instances, 0)
		runningIdleCount := 0

		for _, instance := range instances {
			if instance.PoolID == poolID {
				poolInstances = append(poolInstances, instance)
			}
			if instance.Status == commonParams.InstanceRunning && instance.RunnerStatus == params.RunnerIdle {
				runningIdleCount++
			}
		}
		log.Printf(">>> Pool instances")
		printResponse(poolInstances)
		log.Printf(">>> Running idle count: %d", runningIdleCount)
		log.Printf(">>> Pool min idle runners: %d", pool.MinIdleRunners)
		log.Printf(">>> Pool ID: %s", pool.ID)
		if runningIdleCount == int(pool.MinIdleRunners) && runningIdleCount == len(poolInstances) {
			instance := poolInstances[0]
			// update global variables with instance names
			if pool.RepoID != "" {
				// repo pool
				repoInstanceName = instance.Name
			}
			if pool.OrgID != "" {
				// org pool
				orgInstanceName = instance.Name
			}
			return
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	printResponse(instances)
	panic(fmt.Sprintf("timeout waiting for pool %s to have all idle instances running", poolID))
}

// //////////////////
// / Metrics Token //
// //////////////////
func getMetricsToken(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (string, error) {
	getMetricsTokenResponse, err := apiCli.MetricsToken.GetMetricsToken(
		clientMetricsToken.NewGetMetricsTokenParams(),
		apiAuthToken)
	if err != nil {
		return "", err
	}
	return getMetricsTokenResponse.Payload.Token, nil
}

// ///////////////
// Repositories //
// ///////////////
func createRepo(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoParams params.CreateRepoParams) (*params.Repository, error) {
	createRepoResponse, err := apiCli.Repositories.CreateRepo(
		clientRepositories.NewCreateRepoParams().WithBody(repoParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &createRepoResponse.Payload, nil
}

func listRepos(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (params.Repositories, error) {
	listReposResponse, err := apiCli.Repositories.ListRepos(
		clientRepositories.NewListReposParams(),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listReposResponse.Payload, nil
}

func updateRepo(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID string, repoParams params.UpdateEntityParams) (*params.Repository, error) {
	updateRepoResponse, err := apiCli.Repositories.UpdateRepo(
		clientRepositories.NewUpdateRepoParams().WithRepoID(repoID).WithBody(repoParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &updateRepoResponse.Payload, nil
}

func getRepo(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID string) (*params.Repository, error) {
	getRepoResponse, err := apiCli.Repositories.GetRepo(
		clientRepositories.NewGetRepoParams().WithRepoID(repoID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &getRepoResponse.Payload, nil
}

func installRepoWebhook(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID string, webhookParams params.InstallWebhookParams) (*params.HookInfo, error) {
	installRepoWebhookResponse, err := apiCli.Repositories.InstallRepoWebhook(
		clientRepositories.NewInstallRepoWebhookParams().WithRepoID(repoID).WithBody(webhookParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &installRepoWebhookResponse.Payload, nil
}

func getRepoWebhook(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID string) (*params.HookInfo, error) {
	getRepoWebhookResponse, err := apiCli.Repositories.GetRepoWebhookInfo(
		clientRepositories.NewGetRepoWebhookInfoParams().WithRepoID(repoID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &getRepoWebhookResponse.Payload, nil
}

func createRepoPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID string, poolParams params.CreatePoolParams) (*params.Pool, error) {
	createRepoPoolResponse, err := apiCli.Repositories.CreateRepoPool(
		clientRepositories.NewCreateRepoPoolParams().WithRepoID(repoID).WithBody(poolParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &createRepoPoolResponse.Payload, nil
}

func listRepoPools(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID string) (params.Pools, error) {
	listRepoPoolsResponse, err := apiCli.Repositories.ListRepoPools(
		clientRepositories.NewListRepoPoolsParams().WithRepoID(repoID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listRepoPoolsResponse.Payload, nil
}

func getRepoPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID, poolID string) (*params.Pool, error) {
	getRepoPoolResponse, err := apiCli.Repositories.GetRepoPool(
		clientRepositories.NewGetRepoPoolParams().WithRepoID(repoID).WithPoolID(poolID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &getRepoPoolResponse.Payload, nil
}

func updateRepoPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID, poolID string, poolParams params.UpdatePoolParams) (*params.Pool, error) {
	updateRepoPoolResponse, err := apiCli.Repositories.UpdateRepoPool(
		clientRepositories.NewUpdateRepoPoolParams().WithRepoID(repoID).WithPoolID(poolID).WithBody(poolParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &updateRepoPoolResponse.Payload, nil
}

func listRepoInstances(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID string) (params.Instances, error) {
	listRepoInstancesResponse, err := apiCli.Repositories.ListRepoInstances(
		clientRepositories.NewListRepoInstancesParams().WithRepoID(repoID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listRepoInstancesResponse.Payload, nil
}

func deleteRepo(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID string) error {
	return apiCli.Repositories.DeleteRepo(
		clientRepositories.NewDeleteRepoParams().WithRepoID(repoID),
		apiAuthToken)
}

func deleteRepoPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, repoID, poolID string) error {
	return apiCli.Repositories.DeleteRepoPool(
		clientRepositories.NewDeleteRepoPoolParams().WithRepoID(repoID).WithPoolID(poolID),
		apiAuthToken)
}

// ////////////////
// Organizations //
// ////////////////
func createOrg(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgParams params.CreateOrgParams) (*params.Organization, error) {
	createOrgResponse, err := apiCli.Organizations.CreateOrg(
		clientOrganizations.NewCreateOrgParams().WithBody(orgParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &createOrgResponse.Payload, nil
}

func listOrgs(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (params.Organizations, error) {
	listOrgsResponse, err := apiCli.Organizations.ListOrgs(
		clientOrganizations.NewListOrgsParams(),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listOrgsResponse.Payload, nil
}

func updateOrg(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID string, orgParams params.UpdateEntityParams) (*params.Organization, error) {
	updateOrgResponse, err := apiCli.Organizations.UpdateOrg(
		clientOrganizations.NewUpdateOrgParams().WithOrgID(orgID).WithBody(orgParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &updateOrgResponse.Payload, nil
}

func getOrg(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID string) (*params.Organization, error) {
	getOrgResponse, err := apiCli.Organizations.GetOrg(
		clientOrganizations.NewGetOrgParams().WithOrgID(orgID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &getOrgResponse.Payload, nil
}

func installOrgWebhook(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID string, webhookParams params.InstallWebhookParams) (*params.HookInfo, error) {
	installOrgWebhookResponse, err := apiCli.Organizations.InstallOrgWebhook(
		clientOrganizations.NewInstallOrgWebhookParams().WithOrgID(orgID).WithBody(webhookParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &installOrgWebhookResponse.Payload, nil
}

func getOrgWebhook(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID string) (*params.HookInfo, error) {
	getOrgWebhookResponse, err := apiCli.Organizations.GetOrgWebhookInfo(
		clientOrganizations.NewGetOrgWebhookInfoParams().WithOrgID(orgID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &getOrgWebhookResponse.Payload, nil
}

func createOrgPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID string, poolParams params.CreatePoolParams) (*params.Pool, error) {
	createOrgPoolResponse, err := apiCli.Organizations.CreateOrgPool(
		clientOrganizations.NewCreateOrgPoolParams().WithOrgID(orgID).WithBody(poolParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &createOrgPoolResponse.Payload, nil
}

func listOrgPools(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID string) (params.Pools, error) {
	listOrgPoolsResponse, err := apiCli.Organizations.ListOrgPools(
		clientOrganizations.NewListOrgPoolsParams().WithOrgID(orgID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listOrgPoolsResponse.Payload, nil
}

func getOrgPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID, poolID string) (*params.Pool, error) {
	getOrgPoolResponse, err := apiCli.Organizations.GetOrgPool(
		clientOrganizations.NewGetOrgPoolParams().WithOrgID(orgID).WithPoolID(poolID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &getOrgPoolResponse.Payload, nil
}

func updateOrgPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID, poolID string, poolParams params.UpdatePoolParams) (*params.Pool, error) {
	updateOrgPoolResponse, err := apiCli.Organizations.UpdateOrgPool(
		clientOrganizations.NewUpdateOrgPoolParams().WithOrgID(orgID).WithPoolID(poolID).WithBody(poolParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &updateOrgPoolResponse.Payload, nil
}

func listOrgInstances(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID string) (params.Instances, error) {
	listOrgInstancesResponse, err := apiCli.Organizations.ListOrgInstances(
		clientOrganizations.NewListOrgInstancesParams().WithOrgID(orgID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listOrgInstancesResponse.Payload, nil
}

func deleteOrg(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID string) error {
	return apiCli.Organizations.DeleteOrg(
		clientOrganizations.NewDeleteOrgParams().WithOrgID(orgID),
		apiAuthToken)
}

func deleteOrgPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, orgID, poolID string) error {
	return apiCli.Organizations.DeleteOrgPool(
		clientOrganizations.NewDeleteOrgPoolParams().WithOrgID(orgID).WithPoolID(poolID),
		apiAuthToken)
}

// ////////////
// Instances //
// ////////////
func listInstances(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (params.Instances, error) {
	listInstancesResponse, err := apiCli.Instances.ListInstances(
		clientInstances.NewListInstancesParams(),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listInstancesResponse.Payload, nil
}

func getInstance(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, instanceID string) (*params.Instance, error) {
	getInstancesResponse, err := apiCli.Instances.GetInstance(
		clientInstances.NewGetInstanceParams().WithInstanceName(instanceID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &getInstancesResponse.Payload, nil
}

func deleteInstance(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, instanceID string) error {
	return apiCli.Instances.DeleteInstance(
		clientInstances.NewDeleteInstanceParams().WithInstanceName(instanceID),
		apiAuthToken)
}

// ////////
// Pools //
// ////////
func listPools(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter) (params.Pools, error) {
	listPoolsResponse, err := apiCli.Pools.ListPools(
		clientPools.NewListPoolsParams(),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listPoolsResponse.Payload, nil
}

func getPool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, poolID string) (*params.Pool, error) {
	getPoolResponse, err := apiCli.Pools.GetPool(
		clientPools.NewGetPoolParams().WithPoolID(poolID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &getPoolResponse.Payload, nil
}

func updatePool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, poolID string, poolParams params.UpdatePoolParams) (*params.Pool, error) {
	updatePoolResponse, err := apiCli.Pools.UpdatePool(
		clientPools.NewUpdatePoolParams().WithPoolID(poolID).WithBody(poolParams),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return &updatePoolResponse.Payload, nil
}

func listPoolInstances(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, poolID string) (params.Instances, error) {
	listPoolInstancesResponse, err := apiCli.Instances.ListPoolInstances(
		clientInstances.NewListPoolInstancesParams().WithPoolID(poolID),
		apiAuthToken)
	if err != nil {
		return nil, err
	}
	return listPoolInstancesResponse.Payload, nil
}

func deletePool(apiCli *client.GarmAPI, apiAuthToken runtime.ClientAuthInfoWriter, poolID string) error {
	return apiCli.Pools.DeletePool(
		clientPools.NewDeletePoolParams().WithPoolID(poolID),
		apiAuthToken)
}

// /////////////////
// Main functions //
// /////////////////
//
// /////////////
// Garm Init //
// /////////////
func Login() {
	log.Println(">>> Login")
	loginParams := params.PasswordLoginParams{
		Username: username,
		Password: password,
	}
	token, err := login(cli, loginParams)
	handleError(err)
	printResponse(token)
	authToken = openapiRuntimeClient.BearerToken(token)
	cfg.Managers = []config.Manager{
		{
			Name:    name,
			BaseURL: baseURL,
			Token:   token,
		},
	}
	cfg.ActiveManager = name
	err = cfg.SaveConfig()
	handleError(err)
}

func FirstRun() {
	existingCfg, err := config.LoadConfig()
	handleError(err)
	if existingCfg != nil {
		if existingCfg.HasManager(name) {
			log.Println(">>> Already initialized")
			return
		}
	}

	log.Println(">>> First run")
	newUser := params.NewUserParams{
		Username: username,
		Password: password,
		FullName: fullName,
		Email:    email,
	}
	user, err := firstRun(cli, newUser)
	handleError(err)
	printResponse(user)
}

// ///////////
// Cleanup //
// ///////////
func GracefulCleanup() {
	DisableRepoPool()
	DisableOrgPool()

	DeleteInstance(repoInstanceName)
	DeleteInstance(orgInstanceName)

	WaitRepoPoolNoInstances(6 * time.Minute)
	WaitOrgPoolNoInstances(6 * time.Minute)

	DeleteRepoPool()
	DeleteOrgPool()
	DeletePool()

	DeleteRepo()
	DeleteOrg()
}

func GhOrgRunnersCleanup() {
	if orgPoolID == "" {
		log.Println(">>> No organization pool ID provided, skipping organization runners cleanup")
		return
	}

	log.Println(">>> Org Github runners cleanup")
	client := getGithubClient()
	ghOrgRunners, _, err := client.Actions.ListOrganizationRunners(context.Background(), orgName, nil)
	handleError(err)

	// Remove organization runners
	poolLabel := fmt.Sprintf("runner-pool-id:%s", orgPoolID)
	for _, orgRunner := range ghOrgRunners.Runners {
		for _, label := range orgRunner.Labels {
			if label.GetName() == poolLabel {
				log.Printf(">>> Removing organization runner %s", orgRunner.GetName())
				_, err := client.Actions.RemoveOrganizationRunner(context.Background(), orgName, orgRunner.GetID())
				handleError(err)
			}
		}
	}
}

func GhRepoRunnersCleanup() {
	if repoPoolID == "" {
		log.Println(">>> No repository pool ID provided, skipping repository runners cleanup")
		return
	}

	log.Println(">>> Repo Github runners cleanup")
	client := getGithubClient()
	ghRepoRunners, _, err := client.Actions.ListRunners(context.Background(), orgName, repoName, nil)
	handleError(err)

	// Remove repository runners
	poolLabel := fmt.Sprintf("runner-pool-id:%s", repoPoolID)
	for _, repoRunner := range ghRepoRunners.Runners {
		for _, label := range repoRunner.Labels {
			if label.GetName() == poolLabel {
				log.Printf(">>> Removing repository runner %s", repoRunner.GetName())
				_, err := client.Actions.RemoveRunner(context.Background(), orgName, repoName, repoRunner.GetID())
				handleError(err)
			}
		}
	}
}

// ////////////////////////////
// Credentials and Providers //
// ////////////////////////////
func ListCredentials() {
	log.Println(">>> List credentials")
	credentials, err := listCredentials(cli, authToken)
	handleError(err)
	printResponse(credentials)
}

func ListProviders() {
	log.Println(">>> List providers")
	providers, err := listProviders(cli, authToken)
	handleError(err)
	printResponse(providers)
}

// ////////////////////////
// // Controller info ////
// ////////////////////////
func GetControllerInfo() {
	log.Println(">>> Get controller info")
	controllerInfo, err := getControllerInfo(cli, authToken)
	handleError(err)
	printResponse(controllerInfo)
}

// //////////////////
// / Metrics Token //
// //////////////////
func GetMetricsToken() {
	log.Println(">>> Get metrics token")
	token, err := getMetricsToken(cli, authToken)
	handleError(err)
	printResponse(token)
}

// ///////////////
// Repositories //
// ///////////////
func CreateRepo() {
	repos, err := listRepos(cli, authToken)
	handleError(err)
	if len(repos) > 0 {
		log.Println(">>> Repo already exists, skipping create")
		repoID = repos[0].ID
		return
	}
	log.Println(">>> Create repo")
	createParams := params.CreateRepoParams{
		Owner:           orgName,
		Name:            repoName,
		CredentialsName: credentialsName,
		WebhookSecret:   repoWebhookSecret,
	}
	repo, err := createRepo(cli, authToken, createParams)
	handleError(err)
	printResponse(repo)
	repoID = repo.ID
}

func ListRepos() {
	log.Println(">>> List repos")
	repos, err := listRepos(cli, authToken)
	handleError(err)
	printResponse(repos)
}

func UpdateRepo() {
	log.Println(">>> Update repo")
	updateParams := params.UpdateEntityParams{
		CredentialsName: fmt.Sprintf("%s-clone", credentialsName),
	}
	repo, err := updateRepo(cli, authToken, repoID, updateParams)
	handleError(err)
	printResponse(repo)
}

func GetRepo() {
	log.Println(">>> Get repo")
	repo, err := getRepo(cli, authToken, repoID)
	handleError(err)
	printResponse(repo)
}

func CreateRepoPool() {
	pools, err := listRepoPools(cli, authToken, repoID)
	handleError(err)
	if len(pools) > 0 {
		log.Println(">>> Repo pool already exists, skipping create")
		repoPoolID = pools[0].ID
		return
	}
	log.Println(">>> Create repo pool")
	poolParams := params.CreatePoolParams{
		MaxRunners:     2,
		MinIdleRunners: 0,
		Flavor:         "default",
		Image:          "ubuntu:22.04",
		OSType:         commonParams.Linux,
		OSArch:         commonParams.Amd64,
		ProviderName:   "lxd_local",
		Tags:           []string{"repo-runner"},
		Enabled:        true,
	}
	repo, err := createRepoPool(cli, authToken, repoID, poolParams)
	handleError(err)
	printResponse(repo)
	repoPoolID = repo.ID
}

func ListRepoPools() {
	log.Println(">>> List repo pools")
	pools, err := listRepoPools(cli, authToken, repoID)
	handleError(err)
	printResponse(pools)
}

func GetRepoPool() {
	log.Println(">>> Get repo pool")
	pool, err := getRepoPool(cli, authToken, repoID, repoPoolID)
	handleError(err)
	printResponse(pool)
}

func UpdateRepoPool() {
	log.Println(">>> Update repo pool")
	var maxRunners uint = 5
	var idleRunners uint = 1
	poolParams := params.UpdatePoolParams{
		MinIdleRunners: &idleRunners,
		MaxRunners:     &maxRunners,
	}
	pool, err := updateRepoPool(cli, authToken, repoID, repoPoolID, poolParams)
	handleError(err)
	printResponse(pool)
}

func InstallRepoWebhook() {
	log.Println(">>> Install repo webhook")
	webhookParams := params.InstallWebhookParams{
		WebhookEndpointType: params.WebhookEndpointDirect,
	}
	webhookInfo, err := installRepoWebhook(cli, authToken, repoID, webhookParams)
	handleError(err)
	printResponse(webhookInfo)
}

func GetRepoWebhook() {
	log.Println(">>> Get repo webhook")
	webhookInfo, err := getRepoWebhook(cli, authToken, repoID)
	handleError(err)
	printResponse(webhookInfo)
}

func DisableRepoPool() {
	if repoID == "" {
		log.Println(">>> No repo ID provided, skipping disable repo pool")
		return
	}
	if repoPoolID == "" {
		log.Println(">>> No repo pool ID provided, skipping disable repo pool")
		return
	}

	enabled := false
	_, err := updateRepoPool(cli, authToken, repoID, repoPoolID, params.UpdatePoolParams{Enabled: &enabled})
	handleError(err)
	log.Printf("repo pool %s disabled", repoPoolID)
}

func WaitRepoPoolNoInstances(timeout time.Duration) {
	if repoID == "" {
		log.Println(">>> No repo ID provided, skipping repo pool wait no instances")
		return
	}
	if repoPoolID == "" {
		log.Println(">>> No repo pool ID provided, skipping repo pool wait no instances")
		return
	}

	var timeWaited time.Duration = 0
	var pool *params.Pool
	var err error

	for timeWaited < timeout {
		log.Println(">>> Wait until repo pool has no instances")
		pool, err = getRepoPool(cli, authToken, repoID, repoPoolID)
		handleError(err)
		if len(pool.Instances) == 0 {
			return
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	printResponse(pool)
	for _, instance := range pool.Instances {
		printResponse(instance)
	}

	panic(fmt.Sprintf("Failed to wait for repo pool %s to have no instances", repoPoolID))
}

func WaitRepoInstance(timeout time.Duration) {
	var timeWaited time.Duration = 0
	var instance params.Instance

	for timeWaited < timeout {
		instances, err := listRepoInstances(cli, authToken, repoID)
		handleError(err)
		if len(instances) > 0 {
			instance = instances[0]
			log.Printf("instance %s status: %s", instance.Name, instance.Status)
			if instance.Status == commonParams.InstanceRunning && instance.RunnerStatus == params.RunnerIdle {
				repoInstanceName = instance.Name
				log.Printf("Repo instance %s is in running state", repoInstanceName)
				return
			}
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	repo, err := getRepo(cli, authToken, repoID)
	handleError(err)
	printResponse(repo)

	instanceDetails, err := getInstance(cli, authToken, instance.Name)
	handleError(err)
	printResponse(instanceDetails)

	panic("Failed to wait for repo instance to be ready")
}

func ListRepoInstances() {
	log.Println(">>> List repo instances")
	instances, err := listRepoInstances(cli, authToken, repoID)
	handleError(err)
	printResponse(instances)
}

func DeleteRepo() {
	if repoID == "" {
		log.Println(">>> No repo ID provided, skipping delete repo")
		return
	}

	log.Println(">>> Delete repo")
	err := deleteRepo(cli, authToken, repoID)
	handleError(err)
	log.Printf("repo %s deleted", repoID)
}

func DeleteRepoPool() {
	if repoID == "" {
		log.Println(">>> No repo ID provided, skipping delete repo pool")
		return
	}
	if repoPoolID == "" {
		log.Println(">>> No repo pool ID provided, skipping delete repo pool")
		return
	}

	log.Println(">>> Delete repo pool")
	err := deleteRepoPool(cli, authToken, repoID, repoPoolID)
	handleError(err)
	log.Printf("repo pool %s deleted", repoPoolID)
}

// ////////////////
// Organizations //
// ////////////////
func CreateOrg() {
	orgs, err := listOrgs(cli, authToken)
	handleError(err)
	if len(orgs) > 0 {
		log.Println(">>> Org already exists, skipping create")
		orgID = orgs[0].ID
		return
	}
	log.Println(">>> Create org")
	orgParams := params.CreateOrgParams{
		Name:            orgName,
		CredentialsName: credentialsName,
		WebhookSecret:   orgWebhookSecret,
	}
	org, err := createOrg(cli, authToken, orgParams)
	handleError(err)
	printResponse(org)
	orgID = org.ID
}

func ListOrgs() {
	log.Println(">>> List orgs")
	orgs, err := listOrgs(cli, authToken)
	handleError(err)
	printResponse(orgs)
}

func UpdateOrg() {
	log.Println(">>> Update org")
	updateParams := params.UpdateEntityParams{
		CredentialsName: fmt.Sprintf("%s-clone", credentialsName),
	}
	org, err := updateOrg(cli, authToken, orgID, updateParams)
	handleError(err)
	printResponse(org)
}

func GetOrg() {
	log.Println(">>> Get org")
	org, err := getOrg(cli, authToken, orgID)
	handleError(err)
	printResponse(org)
}

func InstallOrgWebhook() {
	log.Println(">>> Install org webhook")
	webhookParams := params.InstallWebhookParams{
		WebhookEndpointType: params.WebhookEndpointDirect,
	}
	webhookInfo, err := installOrgWebhook(cli, authToken, orgID, webhookParams)
	handleError(err)
	printResponse(webhookInfo)
}

func GetOrgWebhook() {
	log.Println(">>> Get org webhook")
	webhookInfo, err := getOrgWebhook(cli, authToken, orgID)
	handleError(err)
	printResponse(webhookInfo)
}

func CreateOrgPool() {
	pools, err := listOrgPools(cli, authToken, orgID)
	handleError(err)
	if len(pools) > 0 {
		log.Println(">>> Org pool already exists, skipping create")
		orgPoolID = pools[0].ID
		return
	}
	log.Println(">>> Create org pool")
	poolParams := params.CreatePoolParams{
		MaxRunners:     2,
		MinIdleRunners: 0,
		Flavor:         "default",
		Image:          "ubuntu:22.04",
		OSType:         commonParams.Linux,
		OSArch:         commonParams.Amd64,
		ProviderName:   "lxd_local",
		Tags:           []string{"org-runner"},
		Enabled:        true,
	}
	org, err := createOrgPool(cli, authToken, orgID, poolParams)
	handleError(err)
	printResponse(org)
	orgPoolID = org.ID
}

func ListOrgPools() {
	log.Println(">>> List org pools")
	pools, err := listOrgPools(cli, authToken, orgID)
	handleError(err)
	printResponse(pools)
}

func GetOrgPool() {
	log.Println(">>> Get org pool")
	pool, err := getOrgPool(cli, authToken, orgID, orgPoolID)
	handleError(err)
	printResponse(pool)
}

func UpdateOrgPool() {
	log.Println(">>> Update org pool")
	var maxRunners uint = 5
	var idleRunners uint = 1
	poolParams := params.UpdatePoolParams{
		MinIdleRunners: &idleRunners,
		MaxRunners:     &maxRunners,
	}
	pool, err := updateOrgPool(cli, authToken, orgID, orgPoolID, poolParams)
	handleError(err)
	printResponse(pool)
}

// ///////
// Jobs //
// ///////
func TriggerWorkflow(labelName string) {
	log.Printf(">>> Trigger workflow with label %s", labelName)
	client := getGithubClient()

	eventReq := github.CreateWorkflowDispatchEventRequest{
		Ref: "main",
		Inputs: map[string]interface{}{
			"sleep_time":   "50",
			"runner_label": labelName,
		},
	}
	_, err := client.Actions.CreateWorkflowDispatchEventByFileName(context.Background(), orgName, repoName, workflowFileName, eventReq)
	handleError(err)
}

func ValidateJobLifecycle(label string) {
	log.Printf(">>> Validate GARM job lifecycle with label %s", label)

	// wait for job list to be updated
	job := waitLabelledJob(label, 4*time.Minute)

	// check expected job status
	job = waitJobStatus(job.ID, params.JobStatusQueued, 4*time.Minute)
	job = waitJobStatus(job.ID, params.JobStatusInProgress, 4*time.Minute)

	// check expected instance status
	instance := waitInstanceStatus(job.RunnerName, commonParams.InstanceRunning, params.RunnerActive, 5*time.Minute)

	// wait for job to be completed
	waitJobStatus(job.ID, params.JobStatusCompleted, 4*time.Minute)

	// wait for instance to be removed
	waitInstanceToBeRemoved(instance.Name, 5*time.Minute)

	// wait for GARM to rebuild the pool running idle instances
	waitPoolRunningIdleInstances(instance.PoolID, 6*time.Minute)
}

func DisableOrgPool() {
	if orgID == "" {
		log.Println(">>> No org ID provided, skipping disable org pool")
		return
	}
	if orgPoolID == "" {
		log.Println(">>> No org pool ID provided, skipping disable org pool")
		return
	}

	enabled := false
	_, err := updateOrgPool(cli, authToken, orgID, orgPoolID, params.UpdatePoolParams{Enabled: &enabled})
	handleError(err)
	log.Printf("org pool %s disabled", orgPoolID)
}

func WaitOrgPoolNoInstances(timeout time.Duration) {
	if orgID == "" {
		log.Println(">>> No org ID provided, skipping wait for org pool no instances")
		return
	}
	if orgPoolID == "" {
		log.Println(">>> No org pool ID provided, skipping wait for org pool no instances")
		return
	}

	var timeWaited time.Duration = 0
	var pool *params.Pool
	var err error

	for timeWaited < timeout {
		log.Println(">>> Wait until org pool has no instances")
		pool, err = getOrgPool(cli, authToken, orgID, orgPoolID)
		handleError(err)
		if len(pool.Instances) == 0 {
			return
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	printResponse(pool)
	for _, instance := range pool.Instances {
		printResponse(instance)
	}

	panic(fmt.Sprintf("Failed to wait for org pool %s to have no instances", orgPoolID))
}

func WaitOrgInstance(timeout time.Duration) {
	var timeWaited time.Duration = 0
	var instance params.Instance

	for timeWaited < timeout {
		instances, err := listOrgInstances(cli, authToken, orgID)
		handleError(err)
		if len(instances) > 0 {
			instance = instances[0]
			log.Printf("instance %s status: %s", instance.Name, instance.Status)
			if instance.Status == commonParams.InstanceRunning && instance.RunnerStatus == params.RunnerIdle {
				orgInstanceName = instance.Name
				log.Printf("Org instance %s is in running state", orgInstanceName)
				return
			}
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}
	org, err := getOrg(cli, authToken, orgID)
	handleError(err)
	printResponse(org)

	instanceDetails, err := getInstance(cli, authToken, instance.Name)
	handleError(err)
	printResponse(instanceDetails)

	panic("Failed to wait for org instance to be ready")
}

func ListOrgInstances() {
	log.Println(">>> List org instances")
	instances, err := listOrgInstances(cli, authToken, orgID)
	handleError(err)
	printResponse(instances)
}

func DeleteOrg() {
	if orgID == "" {
		log.Println(">>> No org ID provided, skipping delete org")
		return
	}

	log.Println(">>> Delete org")
	err := deleteOrg(cli, authToken, orgID)
	handleError(err)
	log.Printf("org %s deleted", orgID)
}

func DeleteOrgPool() {
	if orgID == "" {
		log.Println(">>> No org ID provided, skipping delete org pool")
		return
	}
	if orgPoolID == "" {
		log.Println(">>> No org pool ID provided, skipping delete org pool")
		return
	}

	log.Println(">>> Delete org pool")
	err := deleteOrgPool(cli, authToken, orgID, orgPoolID)
	handleError(err)
	log.Printf("org pool %s deleted", orgPoolID)
}

// ////////////
// Instances //
// ////////////
func ListInstances() {
	log.Println(">>> List instances")
	instances, err := listInstances(cli, authToken)
	handleError(err)
	printResponse(instances)
}

func GetInstance() {
	log.Println(">>> Get instance")
	instance, err := getInstance(cli, authToken, orgInstanceName)
	handleError(err)
	printResponse(instance)
}

func DeleteInstance(name string) {
	if name == "" {
		log.Println(">>> No instance name provided, skipping delete instance")
		return
	}

	err := deleteInstance(cli, authToken, name)
	for {
		log.Printf(">>> Waiting for instance %s to be deleted", name)
		instances, err := listInstances(cli, authToken)
		handleError(err)
		for _, instance := range instances {
			if instance.Name == name {
				time.Sleep(5 * time.Second)

				continue
			}
		}
		break
	}
	handleError(err)
	log.Printf("instance %s deleted", name)
}

// ////////
// Pools //
// ////////
func CreatePool() {
	pools, err := listPools(cli, authToken)
	handleError(err)
	for _, pool := range pools {
		if pool.Image == "ubuntu:20.04" {
			// this is the extra pool to be deleted, later, via [DELETE] pools dedicated API.
			poolID = pool.ID
			return
		}
	}
	log.Println(">>> Create pool")
	poolParams := params.CreatePoolParams{
		MaxRunners:     2,
		MinIdleRunners: 0,
		Flavor:         "default",
		Image:          "ubuntu:20.04",
		OSType:         commonParams.Linux,
		OSArch:         commonParams.Amd64,
		ProviderName:   "lxd_local",
		Tags:           []string{"ubuntu", "simple-runner"},
		Enabled:        true,
	}
	pool, err := createRepoPool(cli, authToken, repoID, poolParams)
	handleError(err)
	printResponse(pool)
	poolID = pool.ID
}

func ListPools() {
	log.Println(">>> List pools")
	pools, err := listPools(cli, authToken)
	handleError(err)
	printResponse(pools)
}

func UpdatePool() {
	log.Println(">>> Update pool")
	var maxRunners uint = 5
	var idleRunners uint = 0
	poolParams := params.UpdatePoolParams{
		MinIdleRunners: &idleRunners,
		MaxRunners:     &maxRunners,
	}
	pool, err := updatePool(cli, authToken, poolID, poolParams)
	handleError(err)
	printResponse(pool)
}

func GetPool() {
	log.Println(">>> Get pool")
	pool, err := getPool(cli, authToken, poolID)
	handleError(err)
	printResponse(pool)
}

func DeletePool() {
	if poolID == "" {
		log.Println(">>> No pool ID provided, skipping delete pool")
		return
	}

	log.Println(">>> Delete pool")
	err := deletePool(cli, authToken, poolID)
	handleError(err)
	log.Printf("pool %s deleted", poolID)
}

func ListPoolInstances() {
	log.Println(">>> List pool instances")
	instances, err := listPoolInstances(cli, authToken, repoPoolID)
	handleError(err)
	printResponse(instances)
}

func main() {
	/////////////
	// Cleanup //
	/////////////
	defer GhOrgRunnersCleanup()
	defer GhRepoRunnersCleanup()
	defer GracefulCleanup()

	//////////////////
	// initialize cli /
	//////////////////
	garmUrl, err := url.Parse(baseURL)
	handleError(err)
	apiPath, err := url.JoinPath(garmUrl.Path, client.DefaultBasePath)
	handleError(err)
	transportCfg := client.DefaultTransportConfig().
		WithHost(garmUrl.Host).
		WithBasePath(apiPath).
		WithSchemes([]string{garmUrl.Scheme})
	cli = client.NewHTTPClientWithConfig(nil, transportCfg)

	//////////////////
	// garm init //
	//////////////////
	FirstRun()
	Login()

	// ////////////////////////////
	// credentials and providers //
	// ////////////////////////////
	ListCredentials()
	ListProviders()

	// ///////////////////
	// controller info //
	// ///////////////////
	GetControllerInfo()

	////////////////////
	/// metrics token //
	////////////////////
	GetMetricsToken()

	//////////////////
	// repositories //
	//////////////////
	CreateRepo()
	ListRepos()
	UpdateRepo()
	GetRepo()

	//////////////////
	// webhooks //////
	//////////////////
	InstallRepoWebhook()
	GetRepoWebhook()

	CreateRepoPool()
	ListRepoPools()
	GetRepoPool()
	UpdateRepoPool()

	//////////////////
	// organizations //
	//////////////////
	CreateOrg()
	ListOrgs()
	UpdateOrg()
	GetOrg()

	//////////////////
	// webhooks //////
	//////////////////
	InstallOrgWebhook()
	GetOrgWebhook()

	CreateOrgPool()
	ListOrgPools()
	GetOrgPool()
	UpdateOrgPool()

	///////////////
	// instances //
	///////////////
	WaitRepoInstance(6 * time.Minute)
	ListRepoInstances()

	WaitOrgInstance(6 * time.Minute)
	ListOrgInstances()

	ListInstances()
	GetInstance()

	/////////
	// jobs //
	/////////
	TriggerWorkflow("org-runner")
	ValidateJobLifecycle("org-runner")

	TriggerWorkflow("repo-runner")
	ValidateJobLifecycle("repo-runner")

	///////////////
	// pools //
	///////////////
	CreatePool()
	ListPools()
	UpdatePool()
	GetPool()
	ListPoolInstances()
}
