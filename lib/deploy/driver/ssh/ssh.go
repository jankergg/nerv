package ssh

import (
	"golang.org/x/crypto/ssh"
	"bytes"
	"fmt"
	"log"
	"strings"
	"github.com/ChaosXu/nerv/lib/db"
	"github.com/ChaosXu/nerv/lib/env"
	"github.com/ChaosXu/nerv/lib/deploy/driver/util"
)

//RemoteExecute a script on the host of addr
func Execute(addr string, scriptUri string, args map[string]string, credentialRef string) error {
	rep := env.Config().GetMapString("scripts", "repository")
	if rep == "" {
		return fmt.Errorf("scripts.repository isn't setted")
	}
	scriptUrl := rep + scriptUri
	log.Printf("url:%s\n", scriptUrl)
	script, err := util.LoadScript(scriptUrl)
	if err != nil {
		return err
	}

	credential := Credential{}
	pairs := strings.Split(credentialRef, ",")
	if len(pairs) < 2 {
		return fmt.Errorf("error credential: %s", credentialRef)
	}
	if err := db.DB.Where("type=? and name=?", pairs[0], pairs[1]).First(&credential).Error; err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User:credential.User,
		Auth:[]ssh.AuthMethod{
			ssh.Password(credential.Password),
		},
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	export := ""
	for k, v := range args {
		export = export + fmt.Sprintf(" %s=%s", k, v)
	}
	script = "export " + export + " && " + script
	log.Println(script)

	stdoutContent := ""
	stderrContent := ""
	if err := session.Run(script); err != nil {
		stdoutContent = stdout.String()
		if stdoutContent != "" {
			log.Println("stdout\n" + stdoutContent)
		}
		stderrContent = stderr.String()
		if stderrContent != "" {
			log.Println("stderr\n" + stderrContent)
			return fmt.Errorf("%s\n%s", err.Error(), stderrContent)
		} else {
			return err
		}
	} else {
		stdoutContent = stdout.String()
		if stdoutContent != "" {
			log.Println("stdout\n" + stdoutContent)
		}
		stderrContent = stderr.String()
		if stderrContent != "" {
			log.Println("stderr\n" + stderrContent)
			return fmt.Errorf("%s", stderrContent)
		}
	}

	return nil
}