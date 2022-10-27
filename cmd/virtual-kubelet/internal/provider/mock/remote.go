package mock

import "github.com/melbahja/goph"

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

func ExecOnPi(cmd string) (output string, err error) {
	o, err := client.Run(cmd)
	return string(o), err
}

func ExecAsync(cmd string, fn func(out string, err error)) {
	go func() {
		fn(ExecOnPi(cmd))
	}()
}
