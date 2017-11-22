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
	version := strings.TrimPrefix(build.BranchName, "release-")

	if build.State == tc.BuildStateQueued {
		res += color(version, Cyan, false)
	} else if build.State == tc.BuildStatusRunning {
		p := int(build.PercentageComplete / 100 * float32(len(version)))
		res += color(version[:p], "48;5;22m", false)
		res += color(version[p:], Yellow, false)
	} else if build.Status != tc.BuildStatusSuccess {
		res += color(version, Red, false)
	} else {
		res += color(version, Green, false)
	}

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
