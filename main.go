package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/labstack/echo"
)

type hookBody struct {
	ObjectKind   string `json:"object_kind"`
	Before       string `json:"before"`
	After        string `json:"after"`
	Ref          string `json:"ref"`
	CheckoutSha  string `json:"checkout_sha"`
	UserID       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	UserUsername string `json:"user_username"`
	UserEmail    string `json:"user_email"`
	UserAvatar   string `json:"user_avatar"`
	ProjectID    int    `json:"project_id"`
	Project      struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Description       string      `json:"description"`
		WebURL            string      `json:"web_url"`
		AvatarURL         interface{} `json:"avatar_url"`
		GitSSHURL         string      `json:"git_ssh_url"`
		GitHTTPURL        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		Homepage          string      `json:"homepage"`
		URL               string      `json:"url"`
		SSHURL            string      `json:"ssh_url"`
		HTTPURL           string      `json:"http_url"`
	} `json:"project"`
	Repository struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		Description     string `json:"description"`
		Homepage        string `json:"homepage"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		VisibilityLevel int    `json:"visibility_level"`
	} `json:"repository"`
	Commits []struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []string      `json:"added"`
		Modified []string      `json:"modified"`
		Removed  []interface{} `json:"removed"`
	} `json:"commits"`
	TotalCommitsCount int `json:"total_commits_count"`
}

func main() {
	file, err := os.OpenFile("info.log", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)
	e := echo.New()

	e.POST("/hooks", hooks)
	e.Logger.Fatal(e.Start(":6969"))
}

//Create ... Create Handler
func hooks(c echo.Context) error {
	h := new(hookBody)
	if err := c.Bind(h); err != nil {

		return c.JSON(http.StatusBadRequest, "")
	}

	log.Println("Project : ", h.Project.Name, "user :", h.UserName, " push ", getCommitMessages(h.Commits))
	dockerLogout()
	dockerLogin()
	if _, err := os.Stat("bsw-web"); os.IsNotExist(err) {
		gitClone()
	} else {
		gitPull()
	}

	dockerBuild()
	// dockerPush()
	dockerStop()
	dockerRun()

	return c.JSON(http.StatusOK, h)

}

func getCommitMessages(commits []struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	URL       string    `json:"url"`
	Author    struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	Added    []string      `json:"added"`
	Modified []string      `json:"modified"`
	Removed  []interface{} `json:"removed"`
}) string {
	result := ""
	for _, commit := range commits {
		result += commit.Message + ","
	}
	return result
}

func gitClone() {
	log.Println("start gitClone")
	cmdName := "git"
	cmdArgs := []string{"clone", "git@gitlab.com:bsw.com/bsw-web.git"}
	cmd(cmdName, cmdArgs)
	log.Println("gitClone done")
}

func gitPull() {
	log.Println("start gitpull")
	cmdName := "cd"
	cmdArgs := []string{"bsw-web", "&&", "git", "pull", "git@gitlab.com:bsw.com/bsw-web.git"}
	cmd(cmdName, cmdArgs)
	log.Println("gitpull done")
}

func dockerBuild() {
	log.Println("start dockerBuild")
	cmdName := "docker"
	cmdArgs := []string{"build", "-t", "registry.odds.team/hosp/bsw-web:dev", "bsw-web/."}
	cmd(cmdName, cmdArgs)
	log.Println("dockerBuild done")
}

func dockerPush() {
	log.Println("start dockerPush")
	cmdName := "docker"
	cmdArgs := []string{"push", "registry.odds.team/hosp/bsw-web:dev"}
	cmd(cmdName, cmdArgs)
	log.Println("dockerPush done")
}

func dockerLogout() {
	log.Println("start dockerLogout")
	cmdName := "docker"
	cmdArgs := []string{"logout", "registry.gitlab.com"}
	cmd(cmdName, cmdArgs)
	log.Println("dockerLogout done")
}

func dockerRun() {
	log.Println("start dockerRun")
	cmdName := "docker"
	cmdArgs := []string{"run", "--rm", "-d", "--name", "bsw-web", "-p", "4200:80", "registry.odds.team/hosp/bsw-web:dev"}
	cmd(cmdName, cmdArgs)
	log.Println("dockerRun done")

}

func dockerStop() {
	log.Println("start dockerStop")
	cmdName := "docker"
	cmdArgs := []string{"stop", "bsw-web"}
	cmd(cmdName, cmdArgs)
	log.Println("dockerStop done")

}

func dockerLogin() {
	log.Println("start dockerLogin")
	cmdName := "docker"
	cmdArgs := []string{"login", "-u", "gdunghi@gmail.com", "-p", "vbomi0yd", "registry.gitlab.com"}
	cmd(cmdName, cmdArgs)
	log.Println("dockerLogin done")
}

func cmd(cmdName string, cmdArgs []string) error {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating StdoutPipe for Cmd", err)
		// os.Exit(1)
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Println("Command out | ", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Println("Error starting Cmd", err)
		// os.Exit(1)
		return err

	}

	err = cmd.Wait()
	if err != nil {
		log.Println("Error waiting for Cmd", err)
		// os.Exit(1)
		return err

	}
	return nil
}
