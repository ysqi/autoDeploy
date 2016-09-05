package controllers

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// GitHookController work for github hook.
type GitHookController struct {
	beego.Controller
}

// Payload @title Github Payload
// @router /payload [post]
func (g *GitHookController) Payload() {

	event := g.Ctx.Request.Header.Get("X-Github-Event")
	delivery := g.Ctx.Request.Header.Get("X-GitHub-Delivery")
	if event == "" {
		g.Ctx.Output.SetStatus(406)
		g.Ctx.WriteString("HTTP head 'Missing X-Github-Event' is missing")
		return
	}
	if delivery == "" {
		g.Ctx.Output.SetStatus(406)
		g.Ctx.WriteString("HTTP head 'Missing X-Github-Event' is missing")
		return
	}

	data := make(map[string]interface{})
	err := json.Unmarshal(g.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.Error(err)
		g.Abort("500")
	}

	repository := data["repository"].(map[string]interface{})
	if repository == nil {
		g.Abort("500")
	}
	gitURL := repository["html_url"].(string)
	gitInfo, err := beego.AppConfig.GetSection(gitURL)
	if err != nil {
		logs.Error(err)
		g.Abort("500")
	}
	if gitInfo == nil || len(gitInfo) == 0 {
		g.Ctx.Output.SetStatus(406)
		return
	}

	if secret, ok := gitInfo["secret"]; ok {
		signature := g.Ctx.Request.Header.Get("X-Hub-Signature")
		if false == verifySignature([]byte(secret), signature, g.Ctx.Input.RequestBody) {
			g.Ctx.Output.SetStatus(401)
			g.Ctx.WriteString("Invalid signature")
			return
		}
	}

	switch event {
	default:
		g.Ctx.Output.SetStatus(406)
		return
	case "ping":
		g.Ctx.Output.SetStatus(202)
		return
	case "push":
		go func() {
			if err := execShell(gitInfo); err != nil {
				logs.Error(err)
			}
		}()
		g.Ctx.Output.SetStatus(202)
		return

	}

}

func execShell(cfg map[string]string) error {
	shFile := cfg["sh"]
	cmd := exec.Command("sh", shFile)
	cmd.Env = append([]string{"work=" + cfg["work"]}, os.Environ()...)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			logs.Info(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}
