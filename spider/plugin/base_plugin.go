package plugin

type PluginType int

const (
	AfterSchedulerType   = iota
	AfterDownloaderType  = iota
	AfterProcesserType   = iota
	AfterPipelineType    = iota
	BeforeSchedulerType  = iota
	BeforeDownloaderType = iota
	BeforeProcesserType  = iota
	BeforePipelineType   = iota
)

type BasePlugin interface {
	Do(PluginType, ...interface{})
}
