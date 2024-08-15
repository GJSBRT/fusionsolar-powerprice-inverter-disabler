package fusionsolar

import (
	"fmt"
	"bytes"
	"errors"
	"strings"
	"net/http"
	"io/ioutil"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"encoding/hex"
	"crypto/sha512"
	"encoding/json"
)

var (
	ErrFailedLogin 			= fmt.Errorf("Failed to login. Got a non 470 error code")
	ErrRoarandNotSet 		= fmt.Errorf("Cannot make request. 'roarand' is not set. Please send keepalive!")
	ErrFailedToGetRegionUrl = fmt.Errorf("Failed to get region url from validate user response")
)

func (fs *Fusionsolar) getPublicKey() (*PublicKeyResponse, error) {
	request, err := http.NewRequest("GET", "https://eu5.fusionsolar.huawei.com/unisso/pubkey", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "")

	response, err := fs.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Returned a non-200 status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var publicKeyResponse PublicKeyResponse
	err = json.Unmarshal([]byte(body), &publicKeyResponse)
	if err != nil {
		return nil, err
	}

	return &publicKeyResponse, nil
}

func (fs *Fusionsolar) encryptPassword(publicKey string) (string, string, error) {
	block, _ := pem.Decode([]byte(publicKey))

    pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", "", err
	}

	encryptedPassword, err := rsa.EncryptOAEP(sha512.New384(), rand.Reader, pub, []byte(fs.configuration.Password), nil)
	if err != nil {
		return "", "", err
	}

	nonce := hex.EncodeToString([]byte(encryptedPassword[:16]))

	return string(encryptedPassword), nonce, nil
}

func (fs *Fusionsolar) sendLoginRequest(timestamp int64, encryptedPassword string, nonce string) (*ValidateUserResponse, error) {
	url := fmt.Sprintf("https://eu5.fusionsolar.huawei.com/unisso/v3/validateUser.action?timeStamp=%d&nonce=%s", timestamp, nonce)

	validateUserRequest := ValidateUserRequest{
		OrganizationName: 	"",
		Username:			fs.configuration.Username,
		Password:			encryptedPassword,
	}

	jsonBody, err := json.Marshal(validateUserRequest)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "")

	response, err := fs.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Returned a non-200 status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var validateUserResponse ValidateUserResponse
	err = json.Unmarshal([]byte(body), &validateUserResponse)
	if err != nil {
		return nil, err
	}

	// 470 seems to be a success code
	if validateUserResponse.ErrorCode != "470" {
		return nil, ErrFailedLogin
	}

	return &validateUserResponse, nil
}

func (fs *Fusionsolar) completeAuthentication(path string) (error) {
	url := fmt.Sprintf("https://eu5.fusionsolar.huawei.com/%s", path)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Priority", "u=0, i")
	request.Header.Set("User-Agent", "")

	response, err := fs.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Returned a non-200 status code: %d", response.StatusCode))
	}

	return nil
}

func (fs *Fusionsolar) getKeepAlive() (*KeepAliveResponse, error) {
	request, err := http.NewRequest("GET", "https://uni001eu5.fusionsolar.huawei.com/rest/dpcloud/auth/v1/keep-alive", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "")

	response, err := fs.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Returned a non-200 status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var keepAliveResponse KeepAliveResponse
	err = json.Unmarshal([]byte(body), &keepAliveResponse)
	if err != nil {
		return nil, err
	}

	return &keepAliveResponse, nil
}

func (fs *Fusionsolar) logout() (error) {
	if fs.roarand == "" {
		return ErrRoarandNotSet
	}

	request, err := http.NewRequest("DELETE", "https://uni001eu5.fusionsolar.huawei.com/rest/dp/uidm/auth/v1/logout", nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "")
	request.Header.Set("roarand", fs.roarand)

	response, err := fs.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Returned a non-200 status code: %d", response.StatusCode))
	}

	return nil
}

func (fs *Fusionsolar) authenticate() (error) {
	publicKeyResponse, err := fs.getPublicKey()
	if err != nil {
		return err
	}

	encryptedPassword, nonce, err := fs.encryptPassword(publicKeyResponse.PubKey)
	if err != nil {
		return err
	}

	validateUserResponse, err := fs.sendLoginRequest(publicKeyResponse.Timestamp, encryptedPassword, nonce)
	if err != nil {
		return err
	}

	var path string
	for _, entry := range validateUserResponse.RespMultiRegionName {
		if strings.Contains(entry, "/rest/") {
			path = entry
		}
	}

	if path == "" {
		return ErrFailedToGetRegionUrl
	}

	err = fs.completeAuthentication(path)
	if err != nil {
		return err
	}

	keepAliveResponse, err := fs.getKeepAlive()
	if err != nil {
		return err
	}

	fs.roarand = keepAliveResponse.Payload

	return nil
}
