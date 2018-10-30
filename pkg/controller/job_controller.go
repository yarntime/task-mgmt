package controller

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/yarntime/task-mgmt/pkg/types"
	"strconv"
)

const (
	cronNamePrefix = "skyform-ai-job-"
)

type JobController struct {
	config       *types.Config
	currentCount int32
}

func NewJobController(c *types.Config) *JobController {
	return &JobController{
		config:       c,
		currentCount: 0,
	}
}

func componentParams(obj *types.MonitorObject, objParams []string) []string {
	objParams = append(objParams, []string{
		fmt.Sprintf("--host=%s", obj.Host),
		fmt.Sprintf("--instance_name=%s", obj.InstanceName),
		fmt.Sprintf("--kpi=%s", obj.Metric),
		fmt.Sprintf("--index=%s", obj.ESIndex),
		fmt.Sprintf("--doc_type=%s", obj.ESType),
		fmt.Sprintf("--freq=%s", obj.SampleRate),
	}...)
	return objParams
}

func (jc *JobController) CreateCronJob(obj *types.MonitorObject, customConf types.CustomConfig, appConf types.Application, objParams []string) error {
	jobCommand := "python main.py "
	params := componentParams(obj, objParams)
	for i := 0; i < len(params); i++ {
		jobCommand += params[i] + " "
	}

	containerOpts := []string{}
	if appConf.CpuLimit > float32(0) {
		containerOpts = append(containerOpts, fmt.Sprintf("--cpu-shares=%d", int(appConf.CpuLimit*1024)))
	}
	if appConf.MemoryLimit != "" {
		containerOpts = append(containerOpts, fmt.Sprintf("--memory=%s", appConf.MemoryLimit))
	}


	// Create job specification
	jobSpec := types.JobSpec{
		Command:          jobCommand,
		MaxRunTime:       customConf.Global.JobMaxRunTime,
		ContainerImage:   appConf.Image,
		ContainerOptions: containerOpts,
	}
	// Create cron job specification
	cronJobSpec := types.CronSpec{
		Schedule:          appConf.Cron,
		OverlapPolicy:     customConf.Global.ConcurrencyPolicy,
		MaxConcurrentRuns: customConf.Global.MaxConcurrencyRuns,
		MaxFails:          customConf.Global.MaxFails,
		HistoryLimit:      customConf.Global.HistoryLimit,
		JobSpec:           &jobSpec,
	}

	cronJobSpec.Name = cronNamePrefix + strconv.Itoa(appConf.Id) + "-" + strconv.Itoa(obj.ID)
	cronJobSpecBytes, err := json.Marshal(&cronJobSpec)
	if err != nil {
		return err
	}
	glog.V(6).Info("Submit cron " + cronJobSpec.Name)
	if err := CreateCron(cronJobSpecBytes); err != nil {
		return err
	}

	return nil
}

func (jc *JobController) DeleteCronJob(customConf types.CustomConfig) error {
	_, err :=
		runCommand([]string{"cronjob", "remove", "all"}, nil)
	if err != nil {
		return err
	}
	return nil
}
