package encrypter

import (
    "fmt"
    "syscall"
    "unsafe"
    "encoding/binary"
    "unicode/utf16"
)

const (
    CRYPTPROTECT_UI_FORBIDDEN = 0x1
)

var (
    dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
    dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")

    procEncryptData = dllcrypt32.NewProc("CryptProtectData")
    procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
    procLocalFree   = dllkernel32.NewProc("LocalFree")
)

type DATA_BLOB struct {
    cbData uint32
    pbData *byte
}

func newBlob(d []byte) *DATA_BLOB {
    if len(d) == 0 {
        return &DATA_BLOB{}
    }

    return &DATA_BLOB{
        pbData: &d[0],
        cbData: uint32(len(d)),
    }
}

func (b *DATA_BLOB) ToByteArray() []byte {
    d := make([]byte, b.cbData)
    copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
    return d
}

func encrypt(data []byte) ([]byte, error) {
    var outblob DATA_BLOB
    r, _, err := procEncryptData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
    if r == 0 {
        return nil, err
    }
    defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
    return outblob.ToByteArray(), nil
}

func convertToUTF16LittleEndianBytes(s string) []byte {
    u := utf16.Encode([]rune(s))
    b := make([]byte, 2*len(u))
    for index, value := range u {
        binary.LittleEndian.PutUint16(b[index*2:], value)
    }
    return b
}

func EncryptPWD(pwd string) (string, error) {
    s := convertToUTF16LittleEndianBytes(pwd)
    enc, err := encrypt(s)
    if err != nil {
        return "", fmt.Errorf("error while encrypt rpd pwd - %s", err)
    }

    return fmt.Sprintf("%x", enc), nil
}
