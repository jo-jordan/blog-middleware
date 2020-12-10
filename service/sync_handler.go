package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/lzjlxebr/blog-middleware/common"
	"github.com/lzjlxebr/blog-middleware/entity"
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

func ResolveLocalRepo(root string) error {
	// Read the local repo dir hierarchy
	log.Printf("common.LocalDir(): %s \n", common.LocalDir)
	node, err := snowflake.NewNode(1)
	common.ErrorBus(err)

	var blacklist []string
	blogsInTree := make(map[string]entity.Blog)
	var blogsInList []entity.Blog
	// write the resolve result to local JSON file
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		// Black list check
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

		// Info may be nil
		if nil != info {
			if info.IsDir() {
				if !strings.HasPrefix(info.Name(), ".") {
					if path != filepath.FromSlash(root) {
						// Add category
						blogsInTree[info.Name()] = entity.Blog{
							ID:       uint64(node.Generate().Int64()),
							ParentId: 0,
							Type:     0,
							Name:     info.Name(),
						}
					}
				} else {
					blacklist = append(blacklist, info.Name())
				}
			} else {
				for s, cate := range blogsInTree {
					if strings.Contains(path, s) {
						// Add blog to category
						// https://edgeless.me/notes/about-me/resume.md
						url := common.RootUrl + cate.Name + "/" + info.Name()
						cate.Children = append(cate.Children, entity.Blog{
							ID:       uint64(node.Generate().Int64()),
							ParentId: cate.ID,
							Type:     1,
							Name:     url,
						})
						blogsInTree[s] = cate
					}
				}
			}
		}

		return nil
	})

	for _, category := range blogsInTree {
		blogsInList = append(blogsInList, category)
	}
	jsonString, err := json.Marshal(blogsInList)
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
