package api

import (
	"blog-middleware/common"
	"blog-middleware/service"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)
import "blog-middleware/entity"

// Pull will pull from my github repo to local, update my blogs
func Pull(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //解析参数，默认是不会解析的
	w.Header().Set("content-type", "text/plain; charset=UTF-8")
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 - Bad Request!"))
		return
	}

	// pull from my repo
	localRepo := r.Form.Get("local-repo")
	if !common.DirExists(localRepo) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 - Bad Request!"))
		return
	}

	var commands []entity.Command
	commands = append(commands, entity.Command{Name: "cd", Args: []string{localRepo}})
	commands = append(commands, entity.Command{Name: "git", Args: []string{"pull"}})
	err = service.RunMultipleCommands(commands)

	if nil != err {
		log.Printf("err: %s", err)
	}

	common.LocalDir = localRepo[:strings.LastIndex(localRepo, "/")]
	err = service.ResolveLocalRepo(localRepo)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("200 - OK."))
}

// Init should be call only once
func Init(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //解析参数，默认是不会解析的
	w.Header().Set("content-type", "text/plain; charset=UTF-8")
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 - Bad Request!"))
		return
	}

	dir := r.Form.Get("dir")
	if !common.DirExists(dir) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 - Bad Request!"))
		return
	}
	common.LocalDir = dir
	// resolve dirs

	repo := r.Form.Get("repo")

	if "" == repo || "" == dir {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 - Bad Request!"))
		return
	}

	url, err := common.GetRepoName(repo)
	log.Printf("Repo Name: %s", url)
	common.RepoName = url

	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	var commands []entity.Command
	commands = append(commands, entity.Command{Name: "cd", Args: []string{dir}})
	commands = append(commands, entity.Command{Name: "git", Args: []string{"clone", repo}})
	err = service.RunMultipleCommands(commands)

	if nil != err {
		log.Printf("err: %s", err)
	}

	root := common.LocalDir + "/" + common.RepoName
	err = service.ResolveLocalRepo(root)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("200 - OK")))
}

func GetIP(w http.ResponseWriter, r *http.Request) {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(clientIP))
		return
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(ip))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Unknown IP"))
}
