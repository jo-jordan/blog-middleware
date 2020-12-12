package service

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lzjlxebr/blog-middleware/common"
	"github.com/lzjlxebr/blog-middleware/common/hashcode"
	"github.com/lzjlxebr/blog-middleware/dao"
	"github.com/lzjlxebr/blog-middleware/entity"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Multiple command line run logic:
// Command(name, args for name, "command1; command2; command3; ...")
func RunMultipleCommands(commands []entity.Command) error {
	var dir string
	for _, command := range commands {
		lookPath, err := exec.LookPath(command.Name)
		if err != nil {
			log.Fatalf("installing %s is in your future", command.Name)
			return err
		}
		cmd := exec.Command(lookPath, command.Args...)
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

func CloneRepoToLocal() {
	var commands []entity.Command
	commands = append(commands, entity.Command{Name: "cd", Args: []string{common.LocalDir}})
	commands = append(commands, entity.Command{Name: "git", Args: []string{"clone", common.RepoURL}})
	err := RunMultipleCommands(commands)

	common.ErrorBus(err)
}

func ResolveLocalRepo(root string) error {
	// Read the local repo dir hierarchy
	log.Printf("common.LocalDir(): %s \n", common.LocalDir)

	var blacklist []string
	blogsInTree := make(map[string]entity.Blog)
	var blogsInList []entity.Blog

	// write the resolve result to local JSON file
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

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
						curBlog := entity.Blog{
							ID:       uint64(hashcode.String(info.Name())),
							ParentId: 0,
							Type:     0,
							Path:     info.Name(),
							Name:     info.Name(),
						}
						blogsInTree[info.Name()] = curBlog
						blogsInList = append(blogsInList, curBlog)
					}
				} else {
					blacklist = append(blacklist, info.Name())
				}
			} else {
				for s, cate := range blogsInTree {
					if strings.Contains(path, s) {
						// Add blog to category
						// https://edgeless.me/notes/about-me/resume.md
						pp := common.RootUrl + cate.Name + "/" + info.Name()
						curBlog := entity.Blog{
							ID:       uint64(hashcode.String(pp)),
							ParentId: cate.ID,
							Type:     1,
							Path:     pp,
							Name:     info.Name(),
						}
						blogsInTree[s] = cate
						blogsInList = append(blogsInList, curBlog)
					}
				}
			}
		}

		return nil
	})

	// Persistence to MySQL
	dao.BlogSaveAll(blogsInList)

	// Sync to aws s3
	svc := CreateS3Client(common.BucketRegion)
	blogKeys := ListObject(svc, common.BucketName)

	var blogDeleteObjects []s3manager.BatchDeleteObject
	for _, key := range blogKeys {
		blogDeleteObjects = append(blogDeleteObjects, s3manager.BatchDeleteObject{
			Object: &s3.DeleteObjectInput{
				Key:    aws.String(key),
				Bucket: aws.String(common.BucketName),
			},
		})
	}
	// Delete all objects from bucket
	BatchDeleteObject(common.BucketRegion, blogDeleteObjects)

	var blogObjects []s3manager.BatchUploadObject
	for _, blog := range blogsInList {
		if blog.Type == 1 {
			file, err := os.Open(path.Join(common.LocalDir, blog.Path))
			common.ErrorBus(err)
			blogObjects = append(blogObjects,
				s3manager.BatchUploadObject{
					Object: &s3manager.UploadInput{
						Key:    aws.String(blog.Path),
						Bucket: aws.String(common.BucketName),
						Body:   file,
					},
					After: func() error {
						return file.Close()
					},
				},
			)
		}
	}
	// Put
	BatchUploadObject(common.BucketRegion, blogObjects)

	return nil
}
