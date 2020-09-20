package audio

import (
	"github.com/go-logr/logr"
	"github.com/tinyzimmer/kvdi/pkg/audio/gst"
)

// PlaybackPipelineOpts are options passed to the playback pipeline.
type PlaybackPipelineOpts struct {
	PulseServer, DeviceName, SourceFormat string
	SourceRate, SourceChannels            int
}

// NewPlaybackPipeline returns a new Pipeline for audio playback.
func NewPlaybackPipeline(logger logr.Logger, opts *PlaybackPipelineOpts) (*gst.Pipeline, error) {
	cfg := &gst.PipelineConfig{
		Plugins: []*gst.Plugin{
			{
				Name: "pulsesrc",
				Data: map[string]interface{}{
					"server": opts.PulseServer,
					"device": opts.DeviceName,
				},
				SinkCaps: gst.NewRawCaps(opts.SourceFormat, opts.SourceRate, opts.SourceChannels),
			},
			{
				Name: "cutter",
			},
			{
				Name: "opusenc",
			},
			{
				Name: "webmmux",
			},
			{
				InternalSink: true,
			},
		},
	}
	return gst.NewPipelineFromConfig(logger, cfg)
}

// NewSinkPipeline returns a pipeline that dumps audio data to a null device as fast as possible.
// This is useful for flushing the contents of a mic buffer when there are no applications listening
// to it.
func NewSinkPipeline(logger logr.Logger, opts *PlaybackPipelineOpts) (*gst.Pipeline, error) {
	cfg := &gst.PipelineConfig{
		Plugins: []*gst.Plugin{
			{
				Name: "pulsesrc",
				Data: map[string]interface{}{
					"server": opts.PulseServer,
					"device": opts.DeviceName,
				},
				SinkCaps: gst.NewRawCaps(opts.SourceFormat, opts.SourceRate, opts.SourceChannels),
			},
			{
				Name: "fakesink",
				Data: map[string]interface{}{
					"sync": false,
				},
			},
		},
	}
	return gst.NewPipelineFromConfig(logger, cfg)
}
