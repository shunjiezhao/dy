package util

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"time"

	"log"
	"testing"
)

const (
	image         = "mysql:latest"
	containerPort = "3306/tcp"
)

var mysql string

const defaultMongoURI = "dy:123456@tcp(localhost:3306)/dy?charset=utf8mb4&parseTime=True&loc=Local"

var MysqlURL string

func MysqlStartInDocker(m *testing.M) (string, func()) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	cli.ContainerRemove(ctx, "MysqlTest", types.ContainerRemoveOptions{Force: true})
	//, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   false,
		ExposedPorts: map[nat.Port]struct{}{
			containerPort: {},
		},
		Env: []string{"MYSQL_DATABASE=dy", "MYSQL_USER=dy", "MYSQL_PASSWORD=123456", "MYSQL_RANDOM_ROOT_PASSWORD=\"yes\""},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "E:\\DY2023\\pkg\\configs\\sql",
				Target: "/docker-entrypoint-initdb.d"},
		},
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				// port = 0 会分配空闲端口
				{HostIP: "localhost", HostPort: "0"},
			},
		},
	}, nil, nil, "MysqlTest")
	if err != nil {
		log.Fatalln("check docker server is start?")
	}
	containerID := resp.ID

	clean := func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %s", err)
		}
		fmt.Println("remove container")
		cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	}

	err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	// 获取容器信息
	inspect, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}
	// 获取分配端口号
	Portbindings := inspect.NetworkSettings.Ports[containerPort][0]
	fmt.Printf("listening address :%v\n", Portbindings)
	fmt.Printf("MysqlURL Port : %s", Portbindings.HostPort)
	// mongoDB 连接url
	MysqlURL = fmt.Sprintf("dy:123456@tcp(localhost:%s)/dy?charset=utf8mb4&parseTime=True&loc=Local",
		Portbindings.HostPort)

	time.Sleep(time.Second * 14) // 等待 mysql 初始话
	return MysqlURL, clean
}
