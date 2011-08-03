package web

import (
	"io/ioutil"
	"path"
	"sort"
	"strings"
	"metricsd/config"
)

type graphItem struct {
	Name   string
	Writer string
	Group  string
	Title  string
}

func (graph *graphItem) Less(graphToCompare interface{}) bool {
	g := graphToCompare.(*graphItem)
	return graph.Name < g.Name ||
		(graph.Name == g.Name && graph.Writer < g.Writer)
}

type graphItemsList []*graphItem

// Swap exchanges the elements at indexes i and j.
func (l graphItemsList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Swap exchanges the elements at indexes i and j.
func (l graphItemsList) Len() int {
	return len(l)
}

// Swap exchanges the elements at indexes i and j.
func (l graphItemsList) Less(i, j int) bool {
	return l[i].Less(l[j])
}


type graphItemGroup struct {
	Group    string
	HasGroup bool
	Graphs   graphItemsList
}

func (group *graphItemGroup) Less(groupToCompare *graphItemGroup) bool {
	return len(groupToCompare.Group) == 0 || (len(group.Group) > 0 && group.Group < groupToCompare.Group)
}

type graphItemSource struct {
	Source string
	Graphs graphItemsList
}

func (source *graphItemSource) Less(sourceToCompare interface{}) bool {
	s := sourceToCompare.(*graphItemSource)
	return source.Source == "all" || (s.Source != "all" && source.Source < s.Source)
}

type graphItemGroupsList []*graphItemGroup

// Swap exchanges the elements at indexes i and j.
func (l graphItemGroupsList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Swap exchanges the elements at indexes i and j.
func (l graphItemGroupsList) Len() int {
	return len(l)
}

// Swap exchanges the elements at indexes i and j.
func (l graphItemGroupsList) Less(i, j int) bool {
	return l[i].Less(l[j])
}


type Browser struct{}

var browser = &Browser{}

func (browser *Browser) ListCountGraphsGrouped() (groups graphItemGroupsList) {
	groups = make(graphItemGroupsList, 0, 10)
	for _, file := range browser.List("all", "", "-count.rrd") {
		found := false
		for _, group := range groups {
			if file.Group == group.Group {
				group.Graphs = append(group.Graphs, file)
				found = true
			}
		}
		if !found {
			group := &graphItemGroup{file.Group, len(file.Group) > 0, make(graphItemsList, 0, 10)}
			group.Graphs = append(group.Graphs, file)
			groups = append(groups, group)
		}
	}
	sort.Sort(groups)
	for _, group := range groups {
		sort.Sort(group.Graphs)
	}

	return
}

func (browser *Browser) ListSources(metric string) (sources []*graphItemSource) {
	sources = make([]*graphItemSource, 0, 10)
	dir, err := ioutil.ReadDir(path.Join(config.DataDir))
	if err != nil {
		return
	}
	for _, fi := range dir {
		if !fi.IsDirectory() {
			continue
		}
		if graphs := browser.List(fi.Name, metric, ".rrd"); len(graphs) > 0 {
			sources = append(sources, &graphItemSource{fi.Name, graphs})
		}
	}
	return
}

func (*Browser) List(source, metric, suffix string) (files graphItemsList) {
	files = make(graphItemsList, 0, 10)
	dir, err := ioutil.ReadDir(path.Join(config.DataDir, source))
	if err != nil {
		return
	}

	for _, fi := range dir {
		if fi.IsDirectory() {
			continue
		}

		if strings.HasSuffix(fi.Name, suffix) {
			var name, writer, group, title string

			split := strings.LastIndex(fi.Name, "-")
			name = fi.Name[:split]
			if len(metric) > 0 && name != metric {
				continue
			}
			writer = fi.Name[split+1 : len(fi.Name)-len(".rrd")]

			split = strings.Index(name, "$")
			if split < 0 {
				split = strings.Index(name, ".")
			}
			if split >= 0 {
				group = name[:split]
				title = name[split+1:]
			} else {
				group = ""
				title = name
			}
			files = append(files, &graphItem{name, writer, group, title})
		}
	}
	return
}
