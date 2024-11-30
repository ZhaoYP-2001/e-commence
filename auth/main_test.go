package main

import (
	"context"
	"e_commerce/auth/config"
	"e_commerce/kitex_gen/auth"
	"e_commerce/kitex_gen/auth/authservice"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthServer(t *testing.T) {

	authClient, err := authservice.NewClient("auth", client.WithHostPorts(config.ServerAddr))
	if err != nil {
		panic(err)
	}

	user_id := "test_user"
	ctx := context.Background()
	deliverResp, err := authClient.DeliverTokenByRPC(ctx, &auth.DeliverTokenReq{UserId: user_id})
	if err != nil {
		t.Error(err)
	}

	verifyResp, err := authClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{UserId: user_id, AccessToken: deliverResp.AccessToken})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, verifyResp.Res)

	verifyResp, err = authClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{UserId: user_id, AccessToken: deliverResp.RefreshToken})
	t.Log(fmt.Sprintf("verify access-token using refresh-token error: %s", err))
	assert.Equal(t, (*auth.VerifyResp)(nil), verifyResp)

	verifyResp, err = authClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{UserId: user_id, AccessToken: "invalid token"})
	t.Log(fmt.Sprintf("verify access-token using invalid token error: %s", err))
	assert.Equal(t, (*auth.VerifyResp)(nil), verifyResp)

	refreshResp, err := authClient.RefreshTokenByRPC(ctx, &auth.RefreshTokenReq{UserId: user_id, RefreshToken: deliverResp.RefreshToken})
	assert.Equal(t, nil, err)

	verifyResp, err = authClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{UserId: user_id, AccessToken: refreshResp.AccessToken})
	assert.Equal(t, true, verifyResp.Res)

}
