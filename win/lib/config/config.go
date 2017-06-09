package config

import (
    "fmt"
    "path"
    "os"

    log "github.com/jbrodriguez/mlog"

    "gopkg.in/ini.v1"
)

var (
    Version = "0.0.1"

    ServerEndpoint = "https://audit.zhenai.com"

    // shell and rdp exec path
    XShellExec = ""
    SecureCRTExec = ""
    RDPExec = "mstsc"

    // config path
    confPath = path.Join(path.Join(os.Getenv("APPDATA"), "Jumpserver"), "config.ini")
)

func init() {
    serverEndpoint, err := getServerEndpoint()
    if err != nil {
        log.Error(err)
    } else {
        ServerEndpoint = serverEndpoint
    }

    initExecPaths()
}

func getServerEndpoint() (string, error) {
    conf, err := ini.LooseLoad(confPath)
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    serverSec, err := conf.GetSection("Default")
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    server, err := serverSec.GetKey("Server")
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    return server.String(), nil
}

func initExecPaths() {
    xshellExec, err := getExecPath("XShell")
    if err != nil {
        log.Warning("%s", err)
    } else {
        XShellExec = xshellExec
    }

    secureCRT, err := getExecPath("SecureCRT")
    if err != nil {
        log.Warning("%s", err)
    } else {
        SecureCRTExec = secureCRT
    }
}

func getExecPath(execType string) (string, error) {
    conf, err := ini.LooseLoad(confPath)
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    execPath, err := conf.GetSection(execType)
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    path, err := execPath.GetKey("Path")
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    return path.String(), nil
}
