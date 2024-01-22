package main

import (
	"encoding/json"
	"ghu-check/util"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	CheckTokens()
	c := cron.New()
	_, _ = c.AddFunc("@every 15m", CheckTokens)
	c.Start()
	select {}
}

func CheckTokens() {
	// 读取.env文件

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return
	}
	url, device_code, icon_url := os.Getenv("BARK_URL"), os.Getenv("BARK_DEVICE_CODE"), os.Getenv("BARK_ICON_URL")
	ghus := os.Getenv("GHUS")
	tokens := strings.Split(ghus, ",")
	header := util.CopilotHeaders
	var exp_tokens []string
	log.Println("Tokens:", tokens)
	if len(tokens) == 0 {
		log.Println("未配置环境变量GHUS")
		return
	}
	for _, token := range tokens {
		header["Authorization"] = "token " + token
		requestResult, err := util.SendHTTPRequest("GET", "https://api.github.com/copilot_internal/v2/token", header, nil)
		if err != nil {
			util.SendBarkNotice(url, device_code, icon_url, "Cocopilot-Token", "测试请求失败Token: "+token)
			continue
		}
		var coToken CoToken
		_ = json.Unmarshal(requestResult, &coToken)
		if coToken.ExpiresAt < time.Now().Unix() {
			exp_tokens = append(exp_tokens, token)
		}
		time.Sleep(15 * time.Second)
	}
	//发送bark通知
	if len(exp_tokens) > 0 {
		util.SendBarkNotice(url, device_code, icon_url, "Cocopilot-Token", "失效Token: "+strings.Join(exp_tokens, ","))
	}
}

type CoToken struct {
	AnnotationsEnabled                 bool   `json:"annotations_enabled"`
	ChatEnabled                        bool   `json:"chat_enabled"`
	ChatJetbrainsEnabled               bool   `json:"chat_jetbrains_enabled"`
	CodeQuoteEnabled                   bool   `json:"code_quote_enabled"`
	CopilotIdeAgentChatGpt4SmallPrompt bool   `json:"copilot_ide_agent_chat_gpt4_small_prompt"`
	CopilotignoreEnabled               bool   `json:"copilotignore_enabled"`
	ExpiresAt                          int64  `json:"expires_at"`
	IntellijEditorFetcher              bool   `json:"intellij_editor_fetcher"`
	Prompt8K                           bool   `json:"prompt_8k"`
	PublicSuggestions                  string `json:"public_suggestions"`
	RefreshIn                          int    `json:"refresh_in"`
	Sku                                string `json:"sku"`
	SnippyLoadTestEnabled              bool   `json:"snippy_load_test_enabled"`
	Telemetry                          string `json:"telemetry"`
	Token                              string `json:"token"`
	TrackingID                         string `json:"tracking_id"`
	VscPanelV2                         bool   `json:"vsc_panel_v2"`
}
