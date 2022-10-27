package mock

import (
	"fmt"
	"github.com/melbahja/goph"
	v1 "k8s.io/api/core/v1"
	"strings"
)

var (
	client *goph.Client
)

func init() {
	auth, err := goph.UseAgent()
	if err != nil {
		panic(err)
	}
	client, err = goph.New("pi", "192.168.1.46", auth)
	if err != nil {
		panic(err)
	}
}

func ExecuteCmd(cmd string) (output string, err error) {
	o, err := client.Run(cmd)
	return string(o), err
}

func ExecAsync(cmd string, fn func(out string, err error)) {
	go func() {
		fn(ExecuteCmd(cmd))
	}()
}

func GenerateCmd(pod *v1.Pod) string {
	// ignore init containers
	var (
		worker     = pod.Spec.Containers[0]
		env        []string
		workDir    = "/foo/model"
		dataset    string
		step       string
		replaceEnv = map[string]string{
			"jialei-starwhale-controller:8082": "jialei.pre.intra.starwhale.ai",
			"jialei-minio:9000":                "jialei-minio.pre.intra.starwhale.ai:80",
		}
	)

	for _, e := range worker.Env {
		value := e.Value
		for src, dst := range replaceEnv {
			value = strings.ReplaceAll(value, src, dst)
		}
		env = append(env, fmt.Sprintf("%s='%s'", e.Name, value))
		if e.Name == "SW_DATASET_URI" {
			dataset = value
		}
		if e.Name == "SW_TASK_STEP" {
			step = value
		}
	}

	swCmd := fmt.Sprintf("swcli model eval ./src --step=%s --task-index=0 --dataset %s", step, dataset)

	cmd := fmt.Sprintf("cd %s; export %s; %s", workDir, strings.Join(env, " "), swCmd)
	return fmt.Sprintf("docker exec 20 bash -c \"%s\"", cmd)
}
