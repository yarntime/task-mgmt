package controller

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/yarntime/task-mgmt/pkg/mysql"
	v1 "github.com/yarntime/task-mgmt/pkg/types"
	"net/http"
)

type Controller struct {
	dbWorker      *mysql.Worker
	jobController *JobController
	customConfig  v1.CustomConfig
	appConfig     v1.ApplicationConfig
}

func NewController(c *v1.Config) *Controller {
	return &Controller{
		dbWorker:      mysql.NewDBWorker(c),
		jobController: NewJobController(c),
		customConfig:  c.CustomCfg,
		appConfig:     c.AppCfg,
	}
}

func (c *Controller) Create(w http.ResponseWriter, req *http.Request) {
	glog.V(3).Info("scan the monitor objects, create cron jobs")

	c.jobController.DeleteCronJob(c.customConfig)
	monitorObjects := c.dbWorker.List()
	for _, monitorObject := range monitorObjects {
		glog.V(4).Infof("Creating cronjob for monitor object: %v\n", monitorObject)
		for _, appConf := range c.appConfig.App {
			if monitorObject.MonitorTypes&appConf.Id != 0 {
				objParams := c.dbWorker.GetParams(monitorObject.Metric, appConf.Id)
				err := c.jobController.CreateCronJob(monitorObject, c.customConfig, appConf, objParams)
				if err != nil {
					glog.Errorf("Failed to create cronjob: %s", err.Error())
				}
			}
		}
	}

	res := &v1.ApiResponse{
		Code:    200,
		Message: "Successful to scan monitor objects.",
	}

	r, _ := json.Marshal(res)
	w.Write(r)
}

func (c *Controller) Delete(w http.ResponseWriter, req *http.Request) {
	c.jobController.DeleteCronJob(c.customConfig)
	res := &v1.ApiResponse{
		Code:    200,
		Message: "Successful to delete all cron jobs.",
	}

	r, _ := json.Marshal(res)
	w.Write(r)
}
