// Package classification User API.
//
// The purpose of this service is to provide an application
// that is using plain go code to define an API
//
//      Host: localhost
//      Version: 0.0.1
//
// swagger:meta

package web

import (
	"context"
	"encoding/json"
	"net/http"

	"test/test_swagger/master/model"
)

func jsonResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

type ResponseError struct {
	Message string `json:"message"`
}
type APIResponse struct {
	Data  interface{}    `json:"data,omitempty"`
	Error *ResponseError `json:"error,omitempty"`
}

func WriteResponse(w http.ResponseWriter, code int, message string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	resp := ResponseMessage{Code: code, Message: message, Data: data}
	// b, err := json.Marshal(resp)
	// if err != nil {
	// 	logrus.Warnf("error when marshal response message, error:%v\n", err)
	// }
	return json.NewEncoder(w).Encode(resp)
}

//POST /api/v1/user
func (router *SwaRouter) CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	orm := router.orm
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return err
	}

	u := &model.User{}
	if err := orm.Where(user).FirstOrCreate(u).Error; err != nil {
		return err
	}

	return jsonResponse(w, APIResponse{})
}

// swagger:parameters getSingleUser
type GetUserParam struct {
	// an id of user info
	//
	// Required: true
	// in: path
	Id int `json:"id"`
}

type ResponseMessage struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// User Info
//
// swagger:response UserResponse
type UserWapper struct {
	// in: body
	Body ResponseMessage
}

// GET /api/v1/user
func (router *SwaRouter) GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	// swagger:route GET /users/{id} users getSingleUser
	//
	// get a user by userID
	//
	// This will show a user info
	//
	//     Responses:
	//       200: UserResponse
	orm := router.orm

	// get user from db
	user := &model.User{}
	if err := orm.Where("id = ?", vars["id"]).First(user).Error; err != nil {
		return err
	}

	return WriteResponse(w, 200, "success", user)
}
