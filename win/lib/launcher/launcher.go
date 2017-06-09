package launcher

import (
    "fmt"
    "runtime"
    "os/exec"

    log "github.com/jbrodriguez/mlog"

    "../config"
)

func LaunchShell(shellType string, sessionFile string) (error) {
    switch shellType {
        case "XShell":
            return launchXShell(sessionFile)
        case "SecureCRT":
            return launchSecureCRT(sessionFile)
        case "RDP":
            return launchRDP(sessionFile)
        default:
            return fmt.Errorf("%s shell unknow", shellType)
    }
}

func launchXShell(sessionFile string) error {
    if runtime.GOOS == "windows" {
        log.Info("exec: %s %s", config.XShellExec, sessionFile)
        c := exec.Command(config.XShellExec, sessionFile)
        if err := c.Run(); err != nil {
            return fmt.Errorf("%s", err)
        }
    }
    return nil
}

func launchSecureCRT(sessionFile string) error {
    return nil
}

func launchRDP(sessionFile string) error {
    if runtime.GOOS == "windows" {
        log.Info("exec: %s %s", config.RDPExec, sessionFile)
        c := exec.Command(config.RDPExec, sessionFile)
        if err := c.Run(); err != nil {
            return fmt.Errorf("%s", err)
        }
        return nil
    } else {
        return fmt.Errorf("RDP connection way must be run on windows OS")
    }
}
