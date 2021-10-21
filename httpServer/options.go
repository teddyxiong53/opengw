package httpServer

import (
	"embed"

	"github.com/gin-gonic/gin"
)
type FSType uint8
const(
	NormalType FSType = iota
	HTMLType
	FileType
)
type EngineOpts struct {
	engine      *gin.Engine
	embedFS     []EmbedFS
	middlewares gin.HandlersChain
	mode        string
}
type EmbedFS struct{
	Fs embed.FS
	FsPath string
	SubPath string
	Type    FSType
	Data    []byte
}

type Option interface {
	apply(*EngineOpts)
}

type funcOption struct {
	f func(*EngineOpts)
}

func (fo funcOption) apply(opt *EngineOpts) {
	fo.f(opt)
}

var DefaultEngineOption = &EngineOpts{}

func (e *EngineOpts) Apply(Opts ...Option) {
	for _, opt := range Opts {
		opt.apply(e)
	}
	if e.engine==nil{
		gin.SetMode(e.mode)
		e.engine = gin.Default()
	}
}

func WithEmbedFS(embedFss ...EmbedFS) Option {
	return funcOption{
		f: func(eo *EngineOpts) { eo.embedFS = embedFss },
	}
}

func WithMiddlewares(middleware gin.HandlersChain) Option {
	return funcOption{f: func(eo *EngineOpts) { eo.middlewares = append(eo.middlewares, middleware...) }}
}

func WithMode(mode string) Option {
	return funcOption{f: func(eo *EngineOpts) {
		eo.mode = mode

	}}
}

func WithEngine(engine *gin.Engine)Option{
	return funcOption{func(opts *EngineOpts) {
		if engine!=nil{
			opts.engine = engine
		}
	}}
}