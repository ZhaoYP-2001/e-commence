package main

import (
	"context"
	"e_commerce/auth/service"
	"e_commerce/kitex_gen/auth"
	"fmt"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	resp = new(auth.DeliveryResp)
	userId := req.UserId
	prvKey, err := service.NewPrvKey(userId)
	if err != nil {
		return nil, fmt.Errorf("deliver token error: %w", err)
	}
	accessToken, refreshToken, err := service.GenerateToken(prvKey)
	if err != nil {
		return resp, fmt.Errorf("deliver token error: %w", err)
	}

	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	return resp, nil

}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	resp = new(auth.VerifyResp)
	prvKey, err := service.LoadPrvKey(req.UserId)
	if err != nil {
		return resp, fmt.Errorf("verify token error: %w", err)
	}
	err = service.VerifyAccessToken(&prvKey.PublicKey, req.AccessToken)
	if err != nil {
		return resp, fmt.Errorf("verify token error: %w", err)
	}

	resp.Res = true
	return resp, nil
}

// RefreshTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RefreshTokenByRPC(ctx context.Context, req *auth.RefreshTokenReq) (resp *auth.RefreshResp, err error) {
	resp = new(auth.RefreshResp)
	userId := req.UserId
	prvKey, err := service.LoadPrvKey(userId)
	if err != nil {
		return resp, fmt.Errorf("refresh token error: %w", err)
	}

	err = service.VerifyRefreshToken(&prvKey.PublicKey, req.RefreshToken)
	if err != nil {
		return resp, fmt.Errorf("refresh token error: %w", err)
	}

	accessToken, refreshToken, err := service.GenerateToken(prvKey)
	if err != nil {
		return resp, fmt.Errorf("deliver token error: %w", err)
	}

	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	return resp, nil
}
