package parameter

import (
    "fmt"
    "os"
    "strings"
    "encoding/base64"
    "encoding/json"
)

type InputParams struct {
    Token string            `json:"token"`
    AssetID int             `json:"asset_id"`
    ShellType string        `json:"shell_type"`
}

func CheckParams() (*InputParams, error) {
    arg_count := len(os.Args)
    if arg_count < 2 {
        return nil, fmt.Errorf("error - input params are not enough")
    }

    encodeStr := undecorateParamsStr(os.Args[1])
    params := &InputParams{}
    err := decodeParams(encodeStr, params)
    if err != nil {
        return nil, err
    }

    return params, nil
}

func decodeParams(encodeStr string, params *InputParams) (error) {
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

func undecorateParamsStr(params string) (string) {
    params = strings.TrimPrefix(params, "jumpserver://")
    params = strings.TrimSuffix(params, "/")
    return params
}
