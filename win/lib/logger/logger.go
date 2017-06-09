package logger

import (
    "os"
    "path"

    "github.com/jbrodriguez/mlog"
)

func init() {
    fileName := path.Join(os.Getenv("TEMP"), "jms.log")
    mlog.StartEx(mlog.LevelInfo, fileName, 2*1024, 5)
    mlog.Info("Start up")
}
