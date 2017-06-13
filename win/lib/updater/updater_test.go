package updater

import (
    . "testing"
)

func TestCheckUpdate(t *T) {
    latestVersion, _, err := checkUpdate()
    if err != nil {
        t.Errorf("failed in testing checkUpdate - %s", err)
    }

    if len(latestVersion) == 0 {
        t.Errorf("failed in testing checkUpdate - version is empty")
    }
}

func TestCheckVersion(t *T) {
    var v1, v2 string
    var ret bool
    var err error

    v1 = "1.10.2"
    v2 = "1.9.0"
    ret, err = checkVersion(v1, v2)
    if err != nil {
        t.Errorf("%s", err)
    }

    if ret != true {
        t.Error("failed in testing checkVersion")
    }

    ret, err = checkVersion(v2, v1)
    if err != nil {
        t.Errorf("%s", err)
    }

    if ret != false {
        t.Error("failed in testing checkVersion")
    }

    v1 = "1.2.0-2-g36fe0aa"
    v2 = "1.2.1"
    ret, err = checkVersion(v1, v2)
    if err != nil {
        t.Errorf("%s", err)
    }

    if ret != true {
        t.Error("failed in testing checkVersion")
    }

    ret, err = checkVersion(v2, v1)
    if err != nil {
        t.Errorf("%s", err)
    }

    if ret != false {
        t.Error("failed in testing checkVersion")
    }
}
