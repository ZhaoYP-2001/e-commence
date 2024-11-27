package main

import (
	"context"
	"crypto/rsa"
	auth "e_commerce/kitex_gen/auth"
	utils "e_commerce/utils"
	"fmt"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct {
	rsaKey *rsa.PrivateKey
}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	resp = new(auth.DeliveryResp)
	token, err := utils.GenerateToken(s.rsaKey, req.UserId)
	if err != nil {
		return resp, fmt.Errorf("deliver token error: %s", err)
	}
	resp.Token = token
	return
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	resp = new(auth.VerifyResp)
	_, err = utils.VerifyToken(&s.rsaKey.PublicKey, req.Token)
	if err != nil {
		return resp, fmt.Errorf("verify token error: %s", err)
	}
	resp.Res = true
	return resp, nil
}

// GetPubkeyByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) GetPubkeyByRPC(ctx context.Context, req *auth.Empty) (resp *auth.PubkeyResp, err error) {
	resp = new(auth.PubkeyResp)
	pubKey, err := utils.GetRSAPubKeyStr(s.rsaKey)
	if err != nil {
		return resp, fmt.Errorf("get public key error: %s", err)
	}
	resp.PubKey = pubKey
	return resp, nil
}
