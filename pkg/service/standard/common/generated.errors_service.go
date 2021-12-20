// Code generated by errors generator. DO NOT EDIT.

package common

func init() {
	GoFiles = append(GoFiles, errors_service)
}

var errors_service = bytesReader{
	name: "errors_service",
	data: []byte(`[
  {
    "Code": "InternalServiceError",
    "HTTPCode": 500,
    "Message": "Service has some internal Error. Pls Contact With Admin.",
    "Comment": "系统开发兜底的错误提示"
  },
  {
    "Code": "InvalidChargeType",
    "HTTPCode": 400,
    "Message": "ChargeType is not valid.",
    "Comment": "不支持该计费类型，请重新选择计费方式。"
  },
  {
    "Code": "ResourceNotFound",
    "HTTPCode": 404,
    "Message": "The specified resource {{ResourceName}} cannot be found.",
    "Comment": "指定的资源找不到"
  },
  {
    "Code": "DuplicatedResource",
    "HTTPCode": 409,
    "Message": "Resource {{ResourceName}} already exists.",
    "Comment": "指定的资源已经存在"
  },
  {
    "Code": "ServiceFlowLimitExceeded",
    "HTTPCode": 429,
    "Message": "Request was rejected because the request speed of this openAPI is beyond the current flow control limit.",
    "Comment": "请求过于频繁，超出了服务本身的基本限速"
  },
  {
    "Code": "InternalServiceTimeout",
    "HTTPCode": 504,
    "Message": "Internal Service is timeout. Pls Contact With Admin.",
    "Comment": "内部服务执行超时"
  }
]`),
}