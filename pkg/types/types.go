package types

type Application struct {
	Application string   `json:"application"`
	Id          int      `json:"id"`
	Image       string   `json:image`
	Cmd         []string `json:"cmd"`
	Cron        string   `json:"cron"`
	CpuLimit    float32  `json:"cpuLimit"`
	MemoryLimit string   `json:"memoryLimit"`
	Params      []string `json:"params"`
}

type ApplicationConfig struct {
	App []Application `json:"app"`
}

type GlobalConfig struct {
	MysqlHost          string   `json:"mysql_host"`
	MysqlUser          string   `json:"mysql_user"`
	MysqlPwd           string   `json:"mysql_pwd"`
	MysqlDB            string   `json:"mysql_db"`
	Params             []string `json:"params"`
	Namespace          string
	ConcurrencyPolicy  string
	ImagePullPolicy    string
	JobMaxRunTime      string
	MaxConcurrencyRuns int32
	MaxFails           int32
	HistoryLimit       int32
}

type ColumnMap struct {
	Column    string `json:"column"`
	Parameter string `json:"parameter"`
}
type DBParameters struct {
	TableName string      `json:"table_name"`
	Columns   []ColumnMap `json:"column_map"`
}
type CustomConfig struct {
	Global GlobalConfig   `json:"global"`
	DBMap  []DBParameters `json:"DB_map"`
}

type Config struct {
	AppCfg    ApplicationConfig
	CustomCfg CustomConfig
	Host      string
}

type MonitorObject struct {
	ID           int
	Host         string
	InstanceName string
	Metric       string
	MonitorTypes int
	ESIndex      string
	ESType       string
	SampleRate   string
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CronInfo defines the information of an AIP cron job
type CronInfo struct {
	CronSpec *CronSpec  `json:"cronjobspec,omitempty"`
	Finished []*JobInfo `json:"finished,omitempty"`
	Run      []*JobInfo `json:"run,omitempty"`
}

// JobInfo defines the information of an AIP job
type JobInfo struct {
	ID         int64  `json:"id"`
	Status     string `json:"status,omitempty"`
	SubmitTime string `json:"submitTime,omitempty"`
	StartTime  string `json:"startTime,omitempty"`
	EndTime    string `json:"endTime,omitempty"`
	ExitCode   int32  `json:"exitCode,omitempty"`
}

// CronSpec defines a cronjob specification
type CronSpec struct {
	Name     string
	Schedule string
	JobSpec  *JobSpec
	// The following fields are optional
	StartDeadline     int32  `json:",omitempty"`
	OverlapPolicy     string `json:",omitempty"`
	MaxConcurrentRuns int32  `json:",omitempty"`
	MaxFails          int32  `json:",omitempty"`
	HistoryLimit      int32  `json:",omitempty"`
}

// JobSpec specifies the specification of a job.
type JobSpec struct {
	Command          string
	JobType          string    `json:",omitempty"`
	JobName          string    `json:",omitempty"`
	Envs             []string  `json:",omitempty"`
	Queue            string    `json:",omitempty"`
	Hosts            []string  `json:",omitempty"`
	MinNumSlots      int       `json:",omitempty"`
	MaxNumSlots      int       `json:",omitempty"`
	Resource         *Resource `json:",omitempty"`
	BeginTime        string    `json:",omitempty"`
	Deadline         string    `json:",omitempty"`
	MaxRunTime       string    `json:",omitempty"`
	InFile           string    `json:",omitempty"`
	OutFile          *FileOp   `json:",omitempty"`
	ErrFile          *FileOp   `json:",omitempty"`
	Project          string    `json:",omitempty"`
	UserGroup        string    `json:",omitempty"`
	JobGroup         string    `json:",omitempty"`
	JobDescription   string    `json:",omitempty"`
	Cwd              string    `json:",omitempty"`
	Application      string    `json:",omitempty"`
	Interactive      *bool     `json:",omitempty"`
	ShellMode        bool      `json:",omitempty"`
	LoginShell       string    `json:",omitempty"`
	ContainerImage   string    `json:",omitempty"`
	ContainerOptions []string  `json:",omitempty"`
	Debug            bool      `json:",omitempty"`
}

// Resource specifies the resource requirement or a job.
type Resource struct {
	Colocate   *bool
	MaxPerHost *int
	Select     string
	Need       string
	Sort       string
}

// FileOp specifies the path of the file to which the standard output
// and the standard error of a job should be appended.
//
// If Overwrite is set to true, the file is overwritten instead of being
// appended.
type FileOp struct {
	Name      string
	Overwrite bool
}
