package logger

import (
    "os"
    "path"
    "runtime"

    "github.com/jbrodriguez/mlog"
)

func init() {
    if runtime.GOOS == "windows" {
        fileName := path.Join(os.Getenv("TEMP"), "jms.log")
        mlog.StartEx(mlog.LevelInfo, fileName, 5*1024*1024, 5)
    } else {
        // for unit test
        mlog.Start(mlog.LevelInfo, "")
    }
}
