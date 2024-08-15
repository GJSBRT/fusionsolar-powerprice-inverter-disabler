# Fusionsolar package
This package uses the REST API meant for the Fusionsolar dashboard. The API which you are meant to use does not seem to have all functionality, or Huawei's documentation is crap (like always :/).

## Authentication
Authenticating with the API has a few steps.

### 1. Get public key
Send a GET request to `https://eu5.fusionsolar.huawei.com/unisso/pubkey` to get the public key.

**Example response:**
```json
< HTTP/1.1 200 OK
< Server: product only
< Date: Thu, 15 Aug 2024 08:09:39 GMT
< Content-Type: application/json;charset=utf-8
< Transfer-Encoding: chunked
< Connection: keep-alive
< Vary: Accept-Encoding
< cache-control: no-cache,no-store,must-revalidate,max-age=0
< pragma: no-cache
< x-xss-protection: 1;mode=block
< expires: Thu, 01 Jan 1970 00:00:00 GMT
< x-frame-options: SAMEORIGIN
< strict-transport-security: max-age=31536000 ; includeSubDomains
< x-content-type-options: nosniff
< x-download-options: noopen
< content-security-policy: script-src 'self' 'unsafe-inline' 'unsafe-eval';connect-src 'self';form-action 'self'
< x-hiro-resp-from-server: 0
< upcase-conversion-headers: accessSession

{
	"version": "00000983",
	"pubKey": "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAq685CHyto2q8uiPZZPjT/RYFce0iLfubAisAbkoMqUIG9yYYXzIrSZmdvwXe/7EG4lyv9YZSC9Z81hmr1KlGzzZVJB7gk0OdGVjRZX4SbZ58dh3qOsbXZdRd1JSU35cUHs0J03TX40yT7G8U9c3bG9QA+jRuS7utplDukF+0lpnl6PgVpDLOqpeaRRcLvcN1eS9qm9tUya1bAY4SMdjvs5+gkOBopE3ygmGZ2eYdwOLo9Sy2i1QchSTi/jL20xsm/6OiLsP/RXa8NmcV/DZT0C4DL+nHbKacfMgPJukvms0Zu8Ova2CjLmf64S36MyNc945ggNlxZLxAiKO8COQ/fjQpjmKVK6L7MVSdFraM7crOwsHB6q/NDxA+ZLeiC8DtNq0u6R80+iiScdhw/WCZe4vDxPlKtSVTkPOaNlF03ReCx3VvLXwLcSz7ZyD6YBaXF8QpJhKw+fcoyJ9evrWVxjIQAEpbhacf8S4zaY6KXwjrr1DxQ/EGymPzbZjmaI/fAgMBAAE=\n-----END PUBLIC KEY-----",
	"timeStamp": 1723709379363,
	"enableEncrypt": true
}
```

### 2. Encrypt your password with it.
Encrypt your password using the library of your chosing. The public key is in pkcs1 format.

The nonce is the first 16 characters of the encrypted password encoded in hex.

### 3. Send a login request. 
Next send a POST request to `https://eu5.fusionsolar.huawei.com/unisso/v3/validateUser.action?timeStamp=<TIMESTAMP FROM STEP 1>&nonce=<NONCE>`

Use the timestamp from the response of step 1 and encrypted password and nonce from step 2.

**Example request:**
```json
> POST /unisso/v3/validateUser.action?timeStamp=1723709379363&nonce=<REDACTED> HTTP/1.1
> Host: eu5.fusionsolar.huawei.com
> Cookie: UNISSO_V_C_S=<REDACTED>
> Content-Type: application/json
> Content-Length: 600

{
    "organizationName": "",
    "username": "<REDACTED>",
    "password": "<REDACTED>"
}
```

**Example response:**
```json
< HTTP/1.1 200 OK
< Server: product only
< Date: Thu, 15 Aug 2024 08:09:47 GMT
< Content-Type: application/json;charset=utf-8
< Transfer-Encoding: chunked
< Connection: keep-alive
< Vary: Accept-Encoding
< cache-control: no-cache,no-store,must-revalidate,max-age=0
< pragma: no-cache
< x-xss-protection: 1;mode=block
< expires: Thu, 01 Jan 1970 00:00:00 GMT
< x-frame-options: SAMEORIGIN
< strict-transport-security: max-age=31536000 ; includeSubDomains
< x-content-type-options: nosniff
< x-download-options: noopen
< content-security-policy: script-src 'self' 'unsafe-inline' 'unsafe-eval';connect-src 'self';form-action 'self'
< x-hiro-resp-from-server: 0
< upcase-conversion-headers: accessSession

{
	"errorCode": "470",
	"errorMsg": null,
	"redirectURL": null,
	"respMultiRegionName": [
		"-5",
		"/rest/dp/web/v1/auth/on-sso-credential-ready?ticket=<REDACTED>&regionName=region001"
	],
	"serviceUrl": null,
	"verifyCodeCreate": false,
	"twoFactorStatus": null
}
```

### 4. Complete authentication
To complete the authentication process send a GET request to `https://eu5.fusionsolar.huawei.com/rest/dp/web/v1/auth/on-sso-credential-ready?ticket=<TICKET HERE>&regionName=region001`

This path should also be in the response of the request of step 3.

> [!IMPORTANT]
> Make sure to include the following header: `Priority: u=0, i`.

**Example response:**
```json
< HTTP/1.1 200 OK
< Server: product only
< Date: Thu, 15 Aug 2024 08:13:00 GMT
< Content-Type: application/json;charset=UTF-8
< Content-Length: 2
< Connection: keep-alive
< expires: Thu, 01 Jan 1970 00:00:00 GMT
< pragma: no-cache
< cache-control: no-cache, no-store, max-age=0
< x-trace-id: <REDACTED>
< x-span-id: <REDACTED>
< x-parent-id: <REDACTED>
< x-sysprops-sampling: <REDACTED>
< x-autotask-sampling: 1
< x-sampling: true
< x-trace-enable: false
< x-hiro-resp-from-server: 0
< upcase-conversion-headers: accessSession

{}
```

### 5. Send keepalive
Lastly to make non GET request we need some kind of token to include in the headers of requests.

Make a GET request to `https://uni001eu5.fusionsolar.huawei.com/rest/dpcloud/auth/v1/keep-alive`

The payload in the response of this request needs to be added as a `roarand` header in POST requests.

**Example response:**
```json
< HTTP/1.1 200 OK
< Server: CloudWAF
< Date: Thu, 15 Aug 2024 08:19:32 GMT
< Content-Type: application/json
< Transfer-Encoding: chunked
< Connection: keep-alive
< lubanops-gtrace-id: <REDACTED>
< lubanops-nenv-id: 668
< x-envoy-upstream-service-time: 21

{
	"code": 0,
	"message": null,
	"payload": "c-mk<REDACTED>"
}
```
