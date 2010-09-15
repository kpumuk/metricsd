package web

import (
    "container/vector"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "sort"
    "strings"
    "github.com/hoisie/web.go"
    "github.com/hoisie/mustache.go"
    "./config"
)

/***** Web routines ***********************************************************/

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
        "metrics": browser.ListYesNoGraphsGrouped(),
    })
}

func metric(metric string) string {
    return mustache.RenderFile(template("metric"), map[string] interface{} {
        "metric": metric,
        "hosts": browser.ListSources(metric),
    })
}

func host_metric(metric, source string) string {
    return mustache.RenderFile(template("host_metric"), map[string] interface{}{
        "source": source,
        "metric": metric,
        "metrics": browser.List("all", metric, ".rrd"),
    })
}

func host(source string) string {
    return mustache.RenderFile(template("host"), map[string] interface{}{
        "source": source,
        "metrics": browser.List(source, "", "-yesno.rrd"),
    })
}

func metric_graph(metric, source, writer string) string {
    return mustache.RenderFile(template("metric_graph"), map[string] interface{}{
        "source": source,
        "metric": metric,
        "writer": writer,
    })
}

func graph(ctx *web.Context, source, metric, writer string) {
    ctx.SetHeader("Content-Type", "image/png", true)

    var from int
    var rra = "daily"
    if len(ctx.Request.Params["rra"]) > 0 {
        rra = ctx.Request.Params["rra"][0]
    }
    switch rra {
    case "hourly":  from = -14400
    case "daily":   from = -86400
    case "weekly":  from = -604800
    case "monthly": from = -2678400
    case "yearly":  from = -33053184
    default:        from = -86400
    }

    var width  string = "620"
    var height string = "240"
    if w, err := ctx.Request.Params["width"];  err { width  = w[0] }
    if h, err := ctx.Request.Params["height"]; err { height = h[0] }

    rrd_file := fmt.Sprintf("%s/%s/%s-%s.rrd", config.GlobalConfig.DataDir, source, metric, writer)
    args := mustache.RenderFile(template("writers/" + writer), map[string] interface{} {
        "source":   source,
        "metric":   metric,
        "writer":   writer,
        "rrd_file": rrd_file,
        "width":    width,
        "height":   height,
        "from":     from,
        "rra":      rra,
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

/***** Helper functions *******************************************************/

type graphItem struct {
    Name     string
    Writer   string
    Group    string
    Title    string
}

func (graph *graphItem) Less(graphToCompare interface{}) bool {
    g := graphToCompare.(*graphItem)
    return  graph.Name < g.Name ||
            (graph.Name == g.Name && graph.Writer < g.Writer)
}

type graphItemGroup struct {
    Group    string
    HasGroup bool
    Graphs   *vector.Vector
}

func (group *graphItemGroup) Less(groupToCompare interface{}) bool {
    g := groupToCompare.(*graphItemGroup)
    return len(g.Group) == 0 || (len(group.Group) > 0 && group.Group < g.Group)
}

type graphItemSource struct {
    Source string
    Graphs *vector.Vector
}

func (source *graphItemSource) Less(sourceToCompare interface{}) bool {
    s := sourceToCompare.(*graphItemSource)
    return source.Source == "all" || (s.Source != "all" && source.Source < s.Source)
}

type Browser struct {}
var browser = &Browser{}

func (browser *Browser) ListYesNoGraphsGrouped() (groups *vector.Vector) {
    groups = new(vector.Vector)
    for _, elem := range *browser.List("all", "", "-yesno.rrd") {
        file := elem.(*graphItem)
        found := false
        for _, elem := range *groups {
            group := elem.(*graphItemGroup)
            if file.Group == group.Group {
                group.Graphs.Push(file)
                found = true
            }
        }
        if !found {
            group := &graphItemGroup{file.Group, len(file.Group) > 0, new(vector.Vector)}
            group.Graphs.Push(file)
            groups.Push(group)
        }
    }
    sort.Sort(groups)
    for _, elem := range *groups {
        group := elem.(*graphItemGroup)
        sort.Sort(group.Graphs)
    }

    return
}

func (browser *Browser) ListSources(metric string) (sources *vector.Vector) {
    sources = new(vector.Vector)
    dir, err := ioutil.ReadDir(path.Join(config.GlobalConfig.DataDir))
    if err != nil { return }
    for _, fi := range dir {
        if !fi.IsDirectory() { continue }
        if graphs := browser.List(fi.Name, metric, ".rrd"); graphs.Len() > 0 {
            sources.Push(&graphItemSource{fi.Name, graphs})
        }
    }
    return
}

func (*Browser) List(source, metric, suffix string) (files *vector.Vector) {
    files = new(vector.Vector)
    dir, err := ioutil.ReadDir(path.Join(config.GlobalConfig.DataDir, source))
    if err != nil { return }

    for _, fi := range dir {
        if fi.IsDirectory() { continue }

        if strings.HasSuffix(fi.Name, suffix) {
            var name, writer, group, title string

            split := strings.LastIndex(fi.Name, "-")
            name = fi.Name[0:split]
            if len(metric) > 0 && name != metric { continue }
            writer = fi.Name[split + 1:len(fi.Name) - len(".rrd")]

            split = strings.Index(name, "$")
            if split >= 0 {
                group = name[0:split]
                title = name[split+1:]
            } else {
                group = ""
                title = name
            }
            files.Push(&graphItem{name, writer, group, title})
        }
    }
    return
}

func template(name string) string {
    return path.Join(config.GlobalConfig.DataDir, fmt.Sprintf("../templates/%s.mustache", name))
}
