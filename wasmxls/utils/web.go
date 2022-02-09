package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"syscall/js"
)

func Download(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GoBufferToJsUint8Array(b *bytes.Buffer) js.Value {
	arrayConstructor := js.Global().Get("Uint8Array")
	result := arrayConstructor.New(b.Len())
	js.CopyBytesToJS(result, b.Bytes()[0:b.Len()])
	return result
}

func GoBufferToJsBlob(b *bytes.Buffer) js.Value {
	u8array := GoBufferToJsUint8Array(b)
	blobConstructor := js.Global().Get("Blob")
	return blobConstructor.New([]interface{}{u8array})
}

func GoBufferToJsObjectURL(b *bytes.Buffer) js.Value {
	blob := GoBufferToJsBlob(b)
	jsurl := js.Global().Get("URL")
	return jsurl.Call("createObjectURL", blob)
}

// 用于测试的时候提供一个下载按钮
func MakeADownloadToBody(href js.Value, name string) {
	document := js.Global().Get("document")
	a := document.Call("createElement", "a")
	a.Set("href", href)
	a.Set("innerHTML", "下载")
	a.Set("download", name)
	document.Get("body").Call("appendChild", a)
}
