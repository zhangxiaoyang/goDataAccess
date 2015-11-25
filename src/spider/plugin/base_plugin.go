package plugin

type PluginType int

const (
	PreSchedulerType  = iota
	PreDownloaderType = iota
	PreProcesserType  = iota
	PrePipelineType   = iota
)

type BasePlugin interface {
	Do(...interface{})
	GetPluginType() PluginType
}
