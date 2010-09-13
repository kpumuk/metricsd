package web

import (
    "container/vector"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "strings"
    "github.com/hoisie/web.go"
    "github.com/hoisie/mustache.go"
    "./config"
)

func Start() {
    web.Get("/", summary)
    web.Get("/metric/(.*)/(.*)/(.*)", metric_graph)
    web.Get("/metric/(.*)/(.*)", host_metric)
    web.Get("/metric/(.*)", metric)
    web.Get("/graph/(.*)/(.*)/(.*)", graph)
    web.Get("/host/(.*)", host)
    web.Run(config.GlobalConfig.Listen)
}

func summary() string {
    return mustache.RenderFile(template("summary"), map[string] interface{}{
        "metrics": getRrdFiles("all", "", "-yesno.rrd"),
    })
}

func metric(metric string) string {
    return mustache.RenderFile(template("metric"), map[string] interface{} {
        "metric": humanize(metric),
        "hosts": getSources(metric),
    })
}

func host_metric(metric, source string) string {
    return mustache.RenderFile(template("host_metric"), map[string] interface{}{
        "source": source,
        "metric": metric,
        "metricTitle": humanize(metric),
        "metrics": getRrdFiles("all", metric, ".rrd"),
    })
}

func host(source string) string {
    return mustache.RenderFile(template("host"), map[string] interface{}{
        "source": source,
        "metrics": getRrdFiles(source, "", "-yesno.rrd"),
    })
}

func metric_graph(metric, source, writer string) string {
    return mustache.RenderFile(template("metric_graph"), map[string] interface{}{
        "source": source,
        "metric": metric,
        "metricTitle": humanize(metric),
        "writer": writer,
        "writerTitle": humanize(writer),
    })
}

func getRrdFiles(source, metric, suffix string) (files vector.Vector) {
    type elem struct {
        Name  string
        Title string
        Writer string
    }

    dir, err := ioutil.ReadDir(path.Join(config.GlobalConfig.DataDir, source))
    if err != nil { return }
    for _, fi := range dir {
        if fi.IsDirectory() { continue }

        if strings.HasSuffix(fi.Name, suffix) {
            split  := strings.LastIndex(fi.Name, "-")
            name   := fi.Name[0:split]
            if metric != "" && name != metric { continue }
            writer := fi.Name[split + 1:len(fi.Name) - len(".rrd")]
            title  := humanizeGraphName(name, writer)
            files.Push(&elem{name, title, writer})
        }
    }
    return
}

func getSources(metric string) (sources vector.Vector) {
    type elem struct {
        Source string
        Files vector.Vector
    }

    dir, err := ioutil.ReadDir(path.Join(config.GlobalConfig.DataDir))
    if err != nil { return }
    for _, fi := range dir {
        if !fi.IsDirectory() { continue }
        sources.Push(&elem{fi.Name, getRrdFiles(fi.Name, metric, "-yesno.rrd")})
    }
    return
}

func humanizeGraphName(name, writer string) string {
    return fmt.Sprintf("%s :: %s", humanize(name), humanize(writer))
}

func humanize(s string) string {
    return strings.Title(strings.Replace(strings.Replace(s, "_", " ", -1), "-", " — ", -1))
}

func graph(ctx *web.Context, source, metric, writer string) {
    ctx.SetHeader("Content-Type", "image/png", true)

    var from int
    if len(ctx.Request.Params["rra"]) > 0 {
        switch ctx.Request.Params["rra"][0] {
        case "hourly":  from = -14400
        case "daily":   from = -86400
        case "weekly":  from = -604800
        case "monthly": from = -2678400
        case "yearly":  from = -33053184
        default:        from = -86400
        }
    } else {
        from = -86400
    }
    var width  string = "620"
    var height string = "240"
    if w, err := ctx.Request.Params["width"];  err { width  = w[0] }
    if h, err := ctx.Request.Params["height"]; err { height = h[0] }

    rrd_file := fmt.Sprintf("%s/%s/%s-%s.rrd", config.GlobalConfig.DataDir, source, metric, writer)
    args := mustache.RenderFile(fmt.Sprintf("templates/writers/%s.mustache", writer), map[string] interface{} {
        "title":    humanizeGraphName(metric, writer),
        "source":   source,
        "rrd_file": rrd_file,
        "width":    width,
        "height":   height,
        "from":     from,
    })
    r, w, err := os.Pipe()
    if err != nil {
        config.GlobalConfig.Logger.Error("Pipe: %s", err)
        return
    }

    // config.GlobalConfig.Logger.Debug("started, %s", strings.Split(args, "\n", -1))
    pid, err := os.ForkExec("/usr/bin/rrdtool", strings.Split(args, "\n", -1), os.Environ(), "", []*os.File{ nil, w, w })
    w.Close()
    bytes, _ := ioutil.ReadAll(r)
    r.Close()
    ctx.Write(bytes)
    // config.GlobalConfig.Logger.Debug("done")

    wait, err := os.Wait(pid, 0)
    if err != nil {
        config.GlobalConfig.Logger.Error("wait: %s\n", err)
        return
    }
    if !wait.Exited() || wait.ExitStatus() != 0 {
        config.GlobalConfig.Logger.Error("date: %v\n", wait)
        return
    }
    return
}

func template(name string) string {
    return path.Join(config.GlobalConfig.DataDir, fmt.Sprintf("../templates/%s.mustache", name))
}
