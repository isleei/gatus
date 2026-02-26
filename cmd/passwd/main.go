package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "config.yaml"

func main() {
	configPath := defaultConfigPath
	if v := os.Getenv("GATUS_CONFIG_PATH"); v != "" {
		configPath = v
	}
	// Allow override via -c flag
	for i, arg := range os.Args[1:] {
		if arg == "-c" && i+1 < len(os.Args[1:])-1+1 {
			if i+2 < len(os.Args) {
				configPath = os.Args[i+2]
			}
		}
	}

	password, err := readPassword()
	if err != nil {
		fatal("读取密码失败: %v", err)
	}
	if len(password) < 6 {
		fatal("密码长度不能少于 6 个字符")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fatal("生成 bcrypt 哈希失败: %v", err)
	}
	encoded := base64.StdEncoding.EncodeToString(hash)

	data, err := os.ReadFile(configPath)
	if err != nil {
		fatal("读取配置文件 %s 失败: %v", configPath, err)
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fatal("解析配置文件失败: %v", err)
	}

	// Check if security.basic already exists
	sec, _ := cfg["security"].(map[string]interface{})
	if sec == nil {
		sec = map[string]interface{}{}
		cfg["security"] = sec
	}
	basic, _ := sec["basic"].(map[string]interface{})
	if basic == nil {
		basic = map[string]interface{}{}
		sec["basic"] = basic
	}

	oldUser, _ := basic["username"].(string)
	if oldUser == "" {
		oldUser = "admin"
	}

	// Ask for username
	fmt.Printf("用户名 [%s]: ", oldUser)
	var username string
	fmt.Scanln(&username)
	username = strings.TrimSpace(username)
	if username == "" {
		username = oldUser
	}

	basic["username"] = username
	basic["password-bcrypt-base64"] = encoded

	// Write back using line-level replacement to preserve formatting
	updated := updateYAMLSecurity(string(data), username, encoded)
	if updated == "" {
		// Fallback: full marshal
		out, err := yaml.Marshal(cfg)
		if err != nil {
			fatal("序列化配置失败: %v", err)
		}
		updated = string(out)
	}

	if err := os.WriteFile(configPath, []byte(updated), 0644); err != nil {
		fatal("写入配置文件失败: %v", err)
	}

	fmt.Println()
	fmt.Printf("已更新 %s:\n", configPath)
	fmt.Printf("  用户名: %s\n", username)
	fmt.Printf("  密码哈希: %s...\n", encoded[:20])
	fmt.Println("重启 Gatus 后生效。")
}

func readPassword() (string, error) {
	// If stdin is a terminal, read interactively (hidden)
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("请输入新密码: ")
		p1, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}
		fmt.Print("请确认新密码: ")
		p2, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}
		if string(p1) != string(p2) {
			fatal("两次输入的密码不一致")
		}
		return string(p1), nil
	}
	// Non-interactive: read from stdin
	var pw string
	fmt.Scanln(&pw)
	return strings.TrimSpace(pw), nil
}

// updateYAMLSecurity does line-level replacement to preserve the original YAML formatting.
func updateYAMLSecurity(content, username, hash string) string {
	lines := strings.Split(content, "\n")
	hasSecuritySection := false
	hasBasicSection := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "security:" {
			hasSecuritySection = true
		}
		if hasSecuritySection && trimmed == "basic:" {
			hasBasicSection = true
		}
		if hasBasicSection {
			if strings.Contains(trimmed, "username:") {
				indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
				lines[i] = indent + `username: "` + username + `"`
			}
			if strings.Contains(trimmed, "password-bcrypt-base64:") {
				indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
				lines[i] = indent + `password-bcrypt-base64: "` + hash + `"`
				return strings.Join(lines, "\n")
			}
		}
	}

	if !hasSecuritySection {
		// Append security section
		section := fmt.Sprintf("\nsecurity:\n  basic:\n    username: \"%s\"\n    password-bcrypt-base64: \"%s\"\n", username, hash)
		return content + section
	}
	return ""
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
