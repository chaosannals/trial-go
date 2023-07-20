// Code generated by goctl. DO NOT EDIT.
package types

type Request struct {
	Name string `path:"name,options=you|me"`
}

type Response struct {
	Message string `json:"message"`
}

type UserInfo struct {
	Name string `path:"name,options=aaaa|bbbb"`
}

type UserGetRequest struct {
	UserId int `json:"userId"`
}

type UserGetResponse struct {
	Message string   `json:"message"`
	Data    UserInfo `json:"data"`
}

type UserAddRequest struct {
	Data UserInfo `json:"data"`
}

type UserAddResponse struct {
	UserId  int    `json:"userId"`
	Message string `json:"message"`
}
