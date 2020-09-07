package service

import (
	"blog-middleware/common"
	"blog-middleware/entity"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Multiple command line run logic:
// Command(name, args for name, "command1; command2; command3; ...")
func RunMultipleCommands(commands []entity.Command) error {
	var dir string
	for _, command := range commands {
		path, err := exec.LookPath(command.Name)
		if err != nil {
			log.Fatalf("installing %s is in your future", command.Name)
			return err
		}
		cmd := exec.Command(path, command.Args...)
		if strings.Contains(command.Name, "cd") {
			dir = command.Args[0]
		} else {
			cmd.Dir = dir
			if err := cmd.Run(); err != nil {
				return err
			}
			var out bytes.Buffer
			cmd.Stdout = &out
			log.Printf("command: %s, %q\r\n %s", command.Name, command.Args, out.String())
		}
	}
	return nil
}

// Multiple command line run logic:
// entity.Command(name, args for name, "command1; command2; command3; ...")
func runMultipleCommandsInContext(ctx context.Context, commands []entity.Command) *error {
	var dir string
	for _, command := range commands {
		path, err := exec.LookPath(command.Name)
		if err != nil {
			log.Fatalf("installing %s is in your future", command.Name)
			return &err
		}
		cmd := exec.CommandContext(ctx, path, command.Args...)
		if strings.Contains(command.Name, "cd") {
			dir = command.Args[0]
		} else {
			cmd.Dir = dir
			if err := cmd.Run(); err != nil {
				return &err
			}
			var out bytes.Buffer
			cmd.Stdout = &out
			log.Printf("command: %s, %q\r\n %s", command.Name, command.Args, out.String())
		}
	}
	return nil
}

func runInContext(ctx context.Context, name string, arg ...string) *error {
	path, err := exec.LookPath(name)
	if err != nil {
		log.Fatalf("installing %s is in your future", name)
		return &err
	}
	log.Printf("%s is available at %s\n", name, path)
	cmd := exec.CommandContext(ctx, name, arg...)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return &err
	}
	log.Printf("%s", out.String())
	return nil
}

func run(name string, arg ...string) *error {
	path, err := exec.LookPath(name)
	if err != nil {
		log.Fatalf("installing %s is in your future", name)
		return &err
	}
	log.Printf("%s is available at %s\n", name, path)

	cmd := exec.Command(path, arg...)
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return &err
	}
	log.Printf("%s", out.String())
	return nil
}

func ResolveLocalRepo(root string) error {
	// Read the local repo dir hierarchy
	log.Printf("common.LocalDir(): %s \n", common.LocalDir)

	var blacklist []string
	m := make(map[string][]string)
	// write the resolve result to local JSON file
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		isInBlackList := false
		for _, s := range blacklist {
			if strings.Contains(path, s) {
				isInBlackList = true
				break
			}
		}
		if isInBlackList {
			return nil
		}

		if nil != info {
			if info.IsDir() {
				if !strings.HasPrefix(info.Name(), ".") {
					if path != filepath.FromSlash(root) {
						m[path] = []string{}
					}
				} else {
					blacklist = append(blacklist, info.Name())
				}
			} else {
				for s, i := range m {
					if strings.Contains(path, s) {
						m[s] = append(i, path)
					}
				}
			}
		}

		return nil
	})

	jsonString, err := json.Marshal(m)
	if nil != err {
		return err
	}

	config := filepath.FromSlash(fmt.Sprintf("%s%s", common.LocalDir, "/config.json"))
	err = ioutil.WriteFile(config, jsonString, 0644)
	if nil != err {
		return err
	}

	return nil
}
