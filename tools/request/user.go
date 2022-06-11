package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forum/models"
	"forum/tools/authorization"
	"io/ioutil"
	"net/http"
	"os"
)

func GetAllUser(SID string) ([]models.User, error) {
	//init
	All_user := []models.User{}
	url := os.Getenv("url_api") + "users"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//request
	authorization.SetAuthorizationBearer(SID, req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonReqBody map[string]interface{}
	json.Unmarshal(reqBody, &jsonReqBody)
	if _, ok := jsonReqBody["err"]; ok {
		if _, ok := jsonReqBody["msg"]; ok {
			return nil, errors.New(jsonReqBody["msg"].(string))
		}
		return nil, errors.New(jsonReqBody["err"].(string))
	}
	err = json.Unmarshal(reqBody, &All_user)
	if err != nil {
		return nil, err
	}
	return All_user, nil
}

func GetUserByName(name string) (models.User, error) {
	user := models.User{}
	url := os.Getenv("url_api") + "user/by-username/" + name
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return user, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.User{}, err
	}
	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.User{}, err
	}
	var jsonReqBody map[string]interface{}
	json.Unmarshal(reqBody, &jsonReqBody)
	if _, ok := jsonReqBody["err"]; ok {
		if _, ok := jsonReqBody["msg"]; ok {
			return models.User{}, errors.New(jsonReqBody["msg"].(string))
		}
		return models.User{}, errors.New(jsonReqBody["err"].(string))
	}
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetMe(SID string) (models.User, error) {
	user := models.User{}
	url := os.Getenv("url_api") + "user"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return user, err
	}

	authorization.SetAuthorizationBearer(SID, req)
	resp, err := client.Do(req)
	if err != nil {
		return models.User{}, err
	}
	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.User{}, err
	}
	var jsonReqBody map[string]interface{}
	json.Unmarshal(reqBody, &jsonReqBody)
	if _, ok := jsonReqBody["err"]; ok {
		if _, ok := jsonReqBody["msg"]; ok {
			return models.User{}, errors.New(jsonReqBody["msg"].(string))
		}
		return models.User{}, errors.New(jsonReqBody["err"].(string))
	}
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func PutUser(user models.User, SID string) error {
	url := os.Getenv("url_api") + "user"
	client := &http.Client{}
	modifiedUser := make(map[string]interface{})
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = json.Unmarshal(userBytes, &modifiedUser)
	if err != nil {
		return err
	}
	if modifiedUser["profilepicture"].(string) != "" {
		modifiedUser["profilepicture"] = fmt.Sprintf("%02x", user.ProfilePicture)
	} else {
		modifiedUser["profilepicture"] = ""
	}

	data, err := json.Marshal(modifiedUser)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("content-Type", "application/json")
	authorization.SetAuthorizationBearer(SID, req)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var jsonReqBody map[string]interface{}
	json.Unmarshal(reqBody, &jsonReqBody)
	if _, ok := jsonReqBody["err"]; ok {
		if _, ok := jsonReqBody["msg"]; ok {
			return errors.New(jsonReqBody["msg"].(string))
		}
		return errors.New(jsonReqBody["err"].(string))
	}
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		return err
	}
	return nil
}

func GetUserUsername(UUID string) string {
	client := &http.Client{}
	url := os.Getenv("url_api") + "username"

	params := make(map[string]string)
	params["UUID"] = UUID
	data, err := json.Marshal(params)
	if err != nil {
		return ""
	}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return ""
	}
	req.Header.Set("content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var jsonReqBody map[string]string
	json.Unmarshal(reqBody, &jsonReqBody)
	return jsonReqBody["username"]
}

// func GetUserById(sess *session.SessionStore) (models.User, error) {
// 	user := models.User{}
// 	url := os.Getenv("url_api") + "users"
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !tools.IsConnected(sess) {
// 		return nil, errors.New("You are not connected")
// 	}

// }
