package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// Estructuras MCP
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// Función helper para manejar IDs nulos
func safeID(id interface{}) interface{} {
	if id == nil {
		return 0
	}
	return id
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Estructuras para herramientas
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var geminiClient *genai.GenerativeModel

func main() {
	// Cargar .env
	godotenv.Load()

	// Configurar logging a stderr para no interferir con stdio
	log.SetOutput(os.Stderr)
	log.Println("🚀 Iniciando MCP Server para Gemini...")

	// Inicializar cliente Gemini
	if err := initGemini(); err != nil {
		log.Fatalf("Error inicializando Gemini: %v", err)
	}

	// Leer de stdin y escribir a stdout
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		log.Printf("Received message: %s", line)
		
		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			// Enviar error de parsing
			errorResp := MCPResponse{
				JSONRPC: "2.0",
				ID:      0,
				Error: &MCPError{
					Code:    -32700,
					Message: "Parse error: " + err.Error(),
				},
			}
			responseJSON, _ := json.Marshal(errorResp)
			fmt.Println(string(responseJSON))
			continue
		}

		response := handleRequest(req)
		
		responseJSON, _ := json.Marshal(response)
		log.Printf("Sending response: %s", string(responseJSON))
		fmt.Println(string(responseJSON))
	}
}

func initGemini() error {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY no configurada")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return err
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-1.5-flash"
	}

	geminiClient = client.GenerativeModel(model)
	
	// Configuración optimizada
	temp := float32(0.7)
	maxTokens := int32(2048)
	geminiClient.GenerationConfig = genai.GenerationConfig{
		Temperature:     &temp,
		MaxOutputTokens: &maxTokens,
	}

	log.Printf("✅ Gemini client inicializado con modelo: %s", model)
	return nil
}

func handleRequest(req MCPRequest) MCPResponse {
	log.Printf("Received request: %s with ID: %v", req.Method, req.ID)
	
	switch req.Method {
	case "initialize":
		return handleInitialize(req)
	case "tools/list":
		return handleToolsList(req)
	case "tools/call":
		return handleToolsCall(req)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      safeID(req.ID),
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func handleInitialize(req MCPRequest) MCPResponse {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "gemini-mcp-server",
			"version": "1.0.0",
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      safeID(req.ID),
		Result:  result,
	}
}

func handleToolsList(req MCPRequest) MCPResponse {
	tools := []Tool{
		{
			Name:        "ask_gemini",
			Description: "Hacer una consulta a Google Gemini AI",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"prompt": map[string]interface{}{
						"type":        "string",
						"description": "La pregunta o prompt para Gemini",
					},
					"temperature": map[string]interface{}{
						"type":        "number",
						"description": "Temperatura para la generación (0.0-1.0)",
						"minimum":     0.0,
						"maximum":     1.0,
					},
				},
				"required": []string{"prompt"},
			},
		},
		{
			Name:        "analyze_code",
			Description: "Analizar código con Gemini AI",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]interface{}{
						"type":        "string",
						"description": "El código a analizar",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "Lenguaje de programación",
					},
					"task": map[string]interface{}{
						"type":        "string",
						"description": "Tipo de análisis: review, explain, optimize, debug",
						"enum":        []string{"review", "explain", "optimize", "debug"},
					},
				},
				"required": []string{"code", "task"},
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      safeID(req.ID),
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

func handleToolsCall(req MCPRequest) MCPResponse {
	paramsBytes, _ := json.Marshal(req.Params)
	var params CallToolParams
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      safeID(req.ID),
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	switch params.Name {
	case "ask_gemini":
		return handleAskGemini(req, params)
	case "analyze_code":
		return handleAnalyzeCode(req, params)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      safeID(req.ID),
			Error: &MCPError{
				Code:    -32601,
				Message: "Tool not found",
			},
		}
	}
}

func handleAskGemini(req MCPRequest, params CallToolParams) MCPResponse {
	prompt, ok := params.Arguments["prompt"].(string)
	if !ok || prompt == "" {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      safeID(req.ID),
			Error: &MCPError{
				Code:    -32602,
				Message: "prompt is required",
			},
		}
	}

	// Configurar temperatura si se proporciona
	model := geminiClient
	if temp, ok := params.Arguments["temperature"].(float64); ok {
		tempFloat32 := float32(temp)
		tempModel := *geminiClient
		tempModel.GenerationConfig.Temperature = &tempFloat32
		model = &tempModel
	}

	ctx := context.Background()
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      safeID(req.ID),
			Error: &MCPError{
				Code:    -32603,
				Message: "Error generating content: " + err.Error(),
			},
		}
	}

	var responseText string
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			responseText = string(textPart)
		}
	}

	if responseText == "" {
		responseText = "No se pudo generar una respuesta"
	}

	result := ToolResult{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: responseText,
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      safeID(req.ID),
		Result:  result,
	}
}

func handleAnalyzeCode(req MCPRequest, params CallToolParams) MCPResponse {
	code, ok := params.Arguments["code"].(string)
	if !ok || code == "" {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      safeID(req.ID),
			Error: &MCPError{
				Code:    -32602,
				Message: "code is required",
			},
		}
	}

	task, ok := params.Arguments["task"].(string)
	if !ok || task == "" {
		task = "review"
	}

	language, _ := params.Arguments["language"].(string)

	// Construir prompt especializado
	var prompt strings.Builder
	prompt.WriteString("Actúa como un experto desarrollador senior.\n\n")

	switch task {
	case "review":
		prompt.WriteString("Realiza una revisión exhaustiva del siguiente código:\n")
	case "explain":
		prompt.WriteString("Explica de manera clara y detallada el siguiente código:\n")
	case "optimize":
		prompt.WriteString("Optimiza el siguiente código sugiriendo mejoras:\n")
	case "debug":
		prompt.WriteString("Analiza el siguiente código para encontrar y corregir errores:\n")
	}

	if language != "" {
		prompt.WriteString(fmt.Sprintf("Lenguaje: %s\n", language))
	}

	prompt.WriteString("\nCódigo:\n```\n")
	prompt.WriteString(code)
	prompt.WriteString("\n```\n\n")
	prompt.WriteString("Proporciona un análisis detallado y constructivo.")

	ctx := context.Background()
	resp, err := geminiClient.GenerateContent(ctx, genai.Text(prompt.String()))
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      safeID(req.ID),
			Error: &MCPError{
				Code:    -32603,
				Message: "Error analyzing code: " + err.Error(),
			},
		}
	}

	var responseText string
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			responseText = string(textPart)
		}
	}

	if responseText == "" {
		responseText = "No se pudo analizar el código"
	}

	result := ToolResult{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: responseText,
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      safeID(req.ID),
		Result:  result,
	}
}
