package writers

import (
    "container/vector"
    "exec"
    "fmt"
    "log"
    "os"
    "sort"
    "strings"
)


type Writer interface {
    Rollup(time int64, key string, samples *vector.IntVector)
}

func getRrdFile(data string, t string, key string) string {
    return fmt.Sprintf("%s/%s-%s.rrd", data, t, key)
}

func runRrd(data string, argv []string) {
    log.Stdout(strings.Join(argv, " "))
    p, error := exec.Run("/usr/bin/rrdtool", argv, nil, data, exec.PassThrough, exec.PassThrough, exec.PassThrough)
    if error != nil {

    } else {
        if error = p.Close(); error != nil {

        }
    }
}

/******************************************************************************/

type Quartiles struct {
    Data string
}

type Quartile struct {
    time int64
    lo, q1, q2, q3, hi, total int
}

func (quartiles *Quartiles) Rollup(time int64, key string, samples *vector.IntVector) {
    if samples.Len() < 2 { return }
    sort.Sort(samples)
    lo := samples.At(0)
    hi := samples.At(samples.Len() - 1)
    number := samples.Len()
    lo_c := number / 2
    hi_c := number - lo_c
    data := new(Quartile)
    if lo_c > 0 && hi_c > 0 {
        lo_samples := samples.Slice(0, lo_c)
        hi_samples := samples.Slice(lo_c, hi_c)
        lo_sum := 0
        hi_sum := 0
        lo_samples.Do(func(elem interface {}) { lo_sum += elem.(int) })
        hi_samples.Do(func(elem interface {}) { hi_sum += elem.(int) })
        q1 := lo_sum / lo_c
        q2 := (lo_sum + hi_sum) / (lo_c + hi_c)
        q3 := hi_sum / hi_c

        data.time = time
        data.lo = lo
        data.q1 = q1
        data.q2 = q2
        data.q3 = q3
        data.hi = hi
        data.total = number
    }
    quartiles.save(time, key, data)
}

func (self *Quartiles) save(t int64, key string, data *Quartile) {
    file := getRrdFile(self.Data, "quartiles", key)
    log.Stdoutf("File: %s", file)
    if _, err := os.Stat(file); err != nil {
        argv := []string{
            "/usr/bin/rrdtool",
            "create", file,
            "--step", "10",
            "--start", fmt.Sprintf("%d", data.time - 1),
            "DS:q1:GAUGE:600:0:U",
            "DS:q2:GAUGE:600:0:U",
            "DS:q3:GAUGE:600:0:U",
            "DS:hi:GAUGE:600:0:U",
            "DS:lo:GAUGE:600:0:U",
            "DS:total:GAUGE:600:0:U",
            "RRA:AVERAGE:0.5:1:25920",      // 72 hours at 1 sample per 10 secs
            "RRA:AVERAGE:0.5:60:4320",      // 1 month at 1 sample per 10 mins
            "RRA:AVERAGE:0.5:2880:5475",    // 5 years at 1 sample per 8 hours
            "RRA:MIN:0.5:1:25920",          // 72 hours at 1 sample per 10 secs
            "RRA:MIN:0.5:60:4320",          // 1 month at 1 sample per 10 mins
            "RRA:MIN:0.5:2880:5475",        // 5 years at 1 sample per 8 hours
            "RRA:MAX:0.5:1:25920",          // 72 hours at 1 sample per 10 secs
            "RRA:MAX:0.5:60:4320",          // 1 month at 1 sample per 10 mins
            "RRA:MAX:0.5:2880:5475",        // 5 years at 1 sample per 8 hours
        }
        runRrd(self.Data, argv)
    }
    argv := []string{
        "/usr/bin/rrdtool",
        "update", file,
        fmt.Sprintf("%d:%d:%d:%d:%d:%d:%d", data.time, data.q1, data.q2, data.q3, data.lo, data.hi, data.total),
    }
    runRrd(self.Data, argv)
}

/******************************************************************************/

type YesOrNo struct {
    Data string
}

func (self *YesOrNo) Rollup(time int64, key string, samples *vector.IntVector) {
	var ok, fail uint
	samples.Do(func(elem interface{}) {
	    value := elem.(int)
	    if value > 0 {
	        ok += 1
        } else {
            fail += 1
        }
	})
	self.save(time, key, ok, fail)
}

func (self *YesOrNo) save(t int64, key string, ok uint, fail uint) {
    file := getRrdFile(self.Data, "yesno", key)
    log.Stdoutf("File: %s", file)
    if _, err := os.Stat(file); err != nil {
        argv := []string{
            "/usr/bin/rrdtool",
            "create", file,
            "--step", "10",
            "--start", fmt.Sprintf("%d", t - 1),
            "DS:ok:GAUGE:600:0:U",
            "DS:fail:GAUGE:600:0:U",
            "RRA:AVERAGE:0.5:1:25920",      // 72 hours at 1 sample per 10 secs
            "RRA:AVERAGE:0.5:60:4320",      // 1 month at 1 sample per 10 mins
            "RRA:AVERAGE:0.5:2880:5475",    // 5 years at 1 sample per 8 hours
        }
        runRrd(self.Data, argv)
    }
    argv := []string{
        "/usr/bin/rrdtool",
        "update", file,
        fmt.Sprintf("%d:%d:%d", t, ok, fail),
    }
    runRrd(self.Data, argv)
}
