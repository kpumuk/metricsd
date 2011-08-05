package web

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"metricsd/config"
	"github.com/hoisie/web.go"
	"github.com/hoisie/mustache.go"
)

/***** Web routines ***********************************************************/

func Start() {
	web.Config.StaticDir = path.Join(config.RootDir, "public")
	web.Get("/", summary)
	web.Get("/metric/(.*)/(.*)/(.*)", metric_graph)
	web.Get("/metric/(.*)/(.*)", host_metric)
	web.Get("/metric/(.*)", metric)
	web.Get("/graph/(.*)/(.*)/(.*)\\.png", graph)
	web.Get("/graph/(.*)/(.*)/(.*)", graph)
	web.Get("/host/(.*)", host)
	web.Run(config.Listen)
}

func summary() string {
	return mustache.RenderFile(template("summary"), map[string]interface{}{
		"metrics": browser.ListCountGraphsGrouped(),
	})
}

func metric(metric string) string {
	return mustache.RenderFile(template("metric"), map[string]interface{}{
		"metric": metric,
		"hosts":  browser.ListSources(metric),
	})
}

func host_metric(metric, source string) string {
	return mustache.RenderFile(template("host_metric"), map[string]interface{}{
		"source":  source,
		"metric":  metric,
		"metrics": browser.List("all", metric, ".rrd"),
	})
}

func host(source string) string {
	return mustache.RenderFile(template("host"), map[string]interface{}{
		"source":  source,
		"metrics": browser.List(source, "", "-count.rrd"),
	})
}

func metric_graph(metric, source, writer string) string {
	return mustache.RenderFile(template("metric_graph"), map[string]interface{}{
		"source": source,
		"metric": metric,
		"writer": writer,
	})
}

func graph(ctx *web.Context, source, metric, writer string) {
	ctx.SetHeader("Content-Type", "image/png", true)

	params := struct {
		Rra           string
		Width, Height int
		Dark          bool
		Start, End    int
	}{"daily", 620, 240, false, 0, 0}
	ctx.Request.UnmarshalParams(&params)

	if params.Start != 0 || params.End != 0 {
		params.Rra = "custom"
	}

	var start, end int
	switch params.Rra {
	case "hourly":
		start = -14400
		end = -300
	case "daily":
		start = -86400
		end = -300
	case "weekly":
		start = -604800
		end = -1800
	case "monthly":
		start = -2678400
		end = -7200
	case "yearly":
		start = -33053184
		end = -86400
	default:
		start = -86400
		end = -300
		params.Rra = "custom"
	}

	if params.Start == 0 {
		params.Start = start
	}

	if params.End == 0 {
		params.End = end
	}

	rrd_file := fmt.Sprintf("%s/%s/%s-%s.rrd", config.DataDir, source, metric, writer)
	args := mustache.RenderFile(template("writers/"+writer), map[string]interface{}{
		"source":   source,
		"metric":   metric,
		"writer":   writer,
		"rrd_file": rrd_file,
		"start":    params.Start,
		"end":      params.End,
		"width":    params.Width,
		"height":   params.Height,
		"rra":      params.Rra,
		"interval": config.SliceInterval,
		"dark":     params.Dark,
	})
	r, w, err := os.Pipe()
	if err != nil {
		config.Logger.Error("Pipe: %s", err)
		return
	}

	// config.Logger.Debug("started, %s", strings.Split(args, "\n", -1))
	attr := &os.ProcAttr{Dir: "", Env: os.Environ(), Files: []*os.File{nil, w, w}}
	process, err := os.StartProcess("/usr/bin/rrdtool", strings.Split(args, "\n"), attr)
	defer process.Release()
	w.Close()
	io.Copy(ctx, r)
	r.Close()

	wait, err := process.Wait(0)
	if err != nil {
		config.Logger.Error("wait: %s\n", err)
		return
	}
	if !wait.Exited() || wait.ExitStatus() != 0 {
		config.Logger.Error("date: %v\n", wait)
		return
	}
	return
}

/***** Helper functions *******************************************************/

func template(name string) string {
	return path.Join(config.RootDir, fmt.Sprintf("templates/%s.mustache", name))
}
