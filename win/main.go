package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path"
    "runtime"
    "strings"
    "strconv"
    "time"

    "github.com/widuu/goini"
)

var (
    ConfPath string = path.Join(path.Join(os.Getenv("APPDATA"), "Jumpserver"), "config.ini")
)

type InputParams struct {
    Token string            `json:"token"`
    AssetID int             `json:"asset_id"`
    ShellType string        `json:"shell_type"`
}

func init() {
    fileName := path.Join(os.Getenv("TEMP"), "jms.log")
    f, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
		log.Fatalln(err)
    }

    log.SetOutput(f)
    log.Println("Start up")
}

func GetServerEndpoint() (string, error) {
    conf := goini.SetConfig(ConfPath)
    server := conf.GetValue("Default", "Server")
    if len(server) == 0 {
        return "", fmt.Errorf("Error while reading config - server not found in config file")
    }

    return server, nil
}

func GetShellExePath(shellType string) (string, error) {
    conf := goini.SetConfig(ConfPath)
    path := conf.GetValue(shellType, "Path")
    if len(path) == 0 {
        return "", fmt.Errorf("Error while reading config - %s path not found in config file", shellType)
    }

    return path, nil
}

type SessionJson struct {
    Result string           `json:"result"`
    AssetName string        `json:"asset_name"`
    Err bool                `json:"err"`
}


func getJson(token string, assetID int, shellType string, target *SessionJson) error {
    server, err := GetServerEndpoint()
    if err != nil {
        return err
    }

    url := server + "/api/shell/v1/session-file/" + strconv.Itoa(assetID) + "/"
    client := &http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return fmt.Errorf("Error while requesting %s - %s", url, err)
    }

    req.Header.Set("Authorization", "Bearer " + token)

    q := req.URL.Query()
    q.Add("shell_type", shellType)
    req.URL.RawQuery = q.Encode()

    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("Error while downloading %s - %s", url, err)
    }
    defer resp.Body.Close()

    if ! strings.HasPrefix(resp.Status, "2") {
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return err
        }
        return fmt.Errorf(string(b))
    }

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    err = json.Unmarshal(b, &target)
    return err
}

func DownloadSessionFile(token string, assetID int, shellType string, sessionDir string) (string, error) {
    respJson := SessionJson{}
    err := getJson(token, assetID, shellType, &respJson)
    if err != nil {
        return "", fmt.Errorf("Error while get asset(id: %d) session file - %s", assetID, err)
    }

    if respJson.Err {
        return "", fmt.Errorf(respJson.Result)
    }

    assetName := respJson.AssetName
    if len(assetName) == 0 {
        assetName = path.Base(sessionDir)
    }

    sessionFile := path.Join(sessionDir, assetName) + ".xsh"
    err = ioutil.WriteFile(sessionFile, []byte(respJson.Result), 0644)
    if err != nil {
        return "", fmt.Errorf("Error while writing %s - %s", sessionFile, err)
    }

    return sessionFile, nil
}

func LaunchShell(shellType string, sessionFile string) (error) {
    shellExePath, err := GetShellExePath(shellType)
    if err != nil {
        return err
    }

    if shellType == "XShell" {
        log.Println("ready cmd: ", shellExePath, sessionFile)
        if runtime.GOOS == "windows" {
            log.Println("exec: ", shellExePath, sessionFile)
            c := exec.Command(shellExePath, sessionFile)
            if err := c.Run(); err != nil {
                return fmt.Errorf("", err)
            }
        }
    } else {
        return fmt.Errorf("%s is an unknow shell")
    }

    return nil
}


func DecodeParams(encodeStr string, params *InputParams) (error) {
    decodeBytes, err := base64.StdEncoding.DecodeString(encodeStr)
    if err != nil {
        return fmt.Errorf("error while decoding params - %s", err)
    }

    err = json.Unmarshal(decodeBytes, params)
    if err != nil {
        return fmt.Errorf("error while unmarshal params json string - %s", err)
    }

    return nil
}

func UndecorateParamsStr(params string) (string) {
    params = strings.TrimPrefix(params, "jumpserver://")
    params = strings.TrimSuffix(params, "/")
    return params
}

func work() (error) {
    arg_count := len(os.Args)
    if arg_count < 2 {
        return fmt.Errorf("Params are not enough.")
    }

    log.Println("Undecorate params string.")
    encodeStr := UndecorateParamsStr(os.Args[1])
    params := InputParams{}
    log.Println("Decoding params.")
    err := DecodeParams(encodeStr, &params)
	if err != nil {
        return err
	}

    token := params.Token
    id := params.AssetID
    shellType := params.ShellType

    sessionDir, err := ioutil.TempDir("", "jms-")
	if err != nil {
        return err
	}

    //defer os.RemoveAll(sessionDir)

    log.Println("Downloading session file.")
    sessionFile, err := DownloadSessionFile(token, id, shellType, sessionDir)
	if err != nil {
        return err
	}

    log.Println("Launching shell.")
    err = LaunchShell(shellType, sessionFile)
	if err != nil {
		return err
	}

    return nil
}


func main() {
    err := work()
    if err != nil {
		log.Println(err)
        time.Sleep(time.Second * 2)
    }
}
