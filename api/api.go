package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lzjlxebr/blog-middleware/common"
	"github.com/lzjlxebr/blog-middleware/dao"
	"github.com/lzjlxebr/blog-middleware/service"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)
import "github.com/lzjlxebr/blog-middleware/entity"

// Pull will pull from my github repo to local, update my blogs
func Pull(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //解析参数，默认是不会解析的
	w.Header().Set("content-type", "text/plain; charset=UTF-8")
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 - Bad Request!"))
		return
	}

	err = os.Mkdir(common.LocalDir, os.ModePerm)
	common.ErrorBus(err)

	// pull from my repo to local repo /tmp/edgeless-notes/
	var commands []entity.Command
	commands = append(commands, entity.Command{Name: "cd", Args: []string{common.LocalDir}})
	commands = append(commands, entity.Command{Name: "git", Args: []string{"pull"}})
	err = service.RunMultipleCommands(commands)

	if nil != err {
		log.Printf("err: %s", err)
	}

	service.CloneRepoToLocal()

	// then resolve local repo and put them to aws s3
	err = service.ResolveLocalRepo(path.Join(common.LocalDir, "/", common.RepoName))

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("200 - OK."))
}

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	category := vars["category"]
	id := vars["id"]
	if category == "" {
		categories := dao.CategoryFindAll()

		_ = json.NewEncoder(w).Encode(categories)
	} else if id == "" {
		i, err := strconv.ParseUint(category, 10, 64)
		common.ErrorBus(err)
		blogs := dao.BlogFindByCategory(i)
		_ = json.NewEncoder(w).Encode(blogs)
	} else {
		i, err := strconv.ParseUint(id, 10, 64)
		common.ErrorBus(err)
		blog := dao.BlogFindById(i)
		_ = json.NewEncoder(w).Encode(blog)
	}
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
