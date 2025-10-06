package grpc

type tokenInfo struct{}
type serviceTokenInfo struct{}

var AuthToken = tokenInfo{}
var AuthServiceToken = serviceTokenInfo{}
