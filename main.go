package main

import (
	"encoding/json"
	"ghu-check/util"
	"os"
	"strings"
	"time"
)

func main() {
	envVar := os.Getenv("GHUS")
	strArray := strings.Split(envVar, ",")
	CheckTokens(strArray)
}

func CheckTokens(tokens []string) {
	header := util.CopilotHeaders
	var exp_tokens []string
	for _, token := range tokens {
		header["Authorization"] = "token " + token
		requestResult, err := util.SendHTTPRequest("GET", "https://api.github.com/copilot_internal/v2/token", header, nil)
		if err != nil {
			util.SendBarkNotice("Token测试请求失败", "Token: "+token)
			continue
		}
		var coToken CoToken
		_ = json.Unmarshal(requestResult, &coToken)
		if coToken.ExpiresAt < time.Now().Unix() {
			exp_tokens = append(exp_tokens, token)
		}
	}
	//发送bark通知
	util.SendBarkNotice("Token测试结果", "失效Token: "+strings.Join(exp_tokens, ","))
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
