package main

import (
	"context"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aestek/tc/server"
	"github.com/aestek/tc/tc"
	"google.golang.org/grpc"
)

func projectName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return path.Base(dir)
}

func tcEnvInfos(client server.TCServiceClient, pn, env string) string {
	build, err := client.LastBuild(context.Background(), &server.ProjectEnv{
		Project: pn,
		Env:     env,
	})
	if err != nil {
		return ""
	}

	if build == nil {
		return ""
	}

	res := color("["+string(env[0]), Black, false)
	c := Green
	if build.State == tc.BuildStateQueued {
		c = Cyan
	}
	if build.State == tc.BuildStatusRunning {
		c = Yellow
	}
	if build.Status != tc.BuildStatusSuccess {
		c = Red
	}

	res += color(strings.TrimPrefix(build.BranchName, "release-"), c, false)
	res += color("]", Black, false)
	return res
}

func tcInfos() string {
	conn, err := grpc.Dial("127.0.0.1:6363", grpc.WithInsecure())
	if err != nil {
		return ""
	}

	client := server.NewTCServiceClient(conn)

	pn := projectName()

	return tcEnvInfos(client, pn, "stag") + tcEnvInfos(client, pn, "prod")
}
