# 🤖 Gemini-Claude MCP Server

Un servidor MCP (Model Context Protocol) optimizado que conecta Claude Desktop con Google Gemini AI, proporcionando capacidades avanzadas de análisis de código, generación de contenido creativo y consultas generales.

## ✨ Características

- **🔗 Integración nativa con Claude Desktop** - Protocolo MCP optimizado
- **🧠 Modelos Gemini actualizados** - Soporte para Gemini 1.5 Flash/Pro y 2.0 Flash
- **💻 Análisis especializado de código** - Review, debugging, optimización y explicación
- **🎨 Generación de contenido creativo** - Historias, poemas, artículos y emails
- **⚡ Alto rendimiento** - Configuración optimizada para respuestas rápidas
- **🛡️ Seguridad mejorada** - Filtros de contenido configurables y CORS habilitado
- **📊 Métricas detalladas** - Tracking de tokens, tiempo de respuesta y metadatos

## 🚀 Instalación Rápida

### Prerrequisitos

- Go 1.24.2 o superior
- API Key de Google Gemini ([Obtener aquí](https://makersuite.google.com/app/apikey))
- Claude Desktop instalado

### 1. Clonar y configurar

```bash
git clone https://github.com/scopweb/gemini-go-server-mcp.git
cd gemini-claude-mcp-server

# Configurar variables de entorno
export GEMINI_API_KEY="tu_api_key_aqui"
export GEMINI_MODEL="gemini-1.5-flash"  # opcional
export PORT="8080"  # opcional
```

### 2. Instalar dependencias

```bash
go mod tidy
```

### 3. Ejecutar el servidor

```bash
go run main.go
```

### 4. Configurar Claude Desktop

Edita tu archivo de configuración de Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json` en macOS):

```json
{
  "mcpServers": {
    "gemini-mcp-server": {
      "command": "go",
      "args": ["run", "main.go"],
      "cwd": "/ruta/completa/a/tu/proyecto",
      "env": {
        "GEMINI_API_KEY": "tu_api_key_aqui",
        "GEMINI_MODEL": "gemini-1.5-flash",
        "PORT": "8080"
      }
    }
  }
}
```

## 📡 Endpoints API

### 1. Consultas Generales
**POST** `/ask-gemini`

```json
{
  "prompt": "¿Cómo funciona el aprendizaje automático?",
  "temperature": 0.7,
  "max_tokens": 1000,
  "system_role": "Eres un experto en IA que explica conceptos de forma clara"
}
```

### 2. Análisis de Código
**POST** `/analyze-code`

```json
{
  "prompt": "def factorial(n):\n    if n <= 1:\n        return 1\n    return n * factorial(n-1)",
  "language": "python",
  "task": "review",
  "code_context": "Función recursiva para cálculo de factorial",
  "temperature": 0.3
}
```

**Tipos de análisis disponibles:**
- `review` - Revisión exhaustiva del código
- `explain` - Explicación detallada del funcionamiento
- `optimize` - Sugerencias de optimización y mejores prácticas
- `debug` - Detección y corrección de errores

### 3. Contenido Creativo
**POST** `/create-content`

```json
{
  "prompt": "Una historia sobre un desarrollador que descubre un bug que cambia la realidad",
  "content_type": "story",
  "style": "ciencia ficción",
  "length": "medium",
  "temperature": 0.9
}
```

**Tipos de contenido:**
- `story` - Narrativas y cuentos
- `poem` - Poesía y versos
- `article` - Artículos informativos
- `email` - Comunicaciones profesionales

### 4. Estado del Servidor
**GET** `/health`

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "model": "gemini-1.5-flash",
  "version": "2.0.0"
}
```

### 5. Información de Modelos
**GET** `/models`

```json
{
  "current_model": "gemini-1.5-flash",
  "available_models": ["gemini-1.5-flash", "gemini-1.5-pro", "gemini-2.0-flash-exp"],
  "capabilities": {
    "text_generation": true,
    "code_analysis": true,
    "creative_writing": true,
    "multilingual": true,
    "reasoning": true
  }
}
```

## 🏗️ Estructura del Proyecto

```
gemini-claude-mcp-server/
├── main.go              # Código principal del servidor
├── go.mod               # Dependencias de Go
├── go.sum               # Checksums de dependencias
├── Dockerfile           # Configuración para Docker
├── README.md            # Esta documentación
└── examples/            # Ejemplos de uso
    ├── basic_query.json
    ├── code_analysis.json
    └── creative_content.json
```

## ⚙️ Configuración Avanzada

### Variables de Entorno

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `GEMINI_API_KEY` | **Requerido** - Tu API key de Gemini | - |
| `GEMINI_MODEL` | Modelo de Gemini a utilizar | `gemini-1.5-flash` |
| `PORT` | Puerto del servidor | `8080` |

### Modelos Disponibles

- **`gemini-1.5-flash`** - Rápido y eficiente para la mayoría de tareas
- **`gemini-1.5-pro`** - Más potente para tareas complejas
- **`gemini-2.0-flash-exp`** - Modelo experimental más reciente

### Configuración de Parámetros

```go
// Configuración por defecto
config = ServerConfig{
    Temperature: 0.7,    // Creatividad (0.0-1.0)
    MaxTokens:   2048,   // Máximo de tokens de salida
    TopP:        0.95,   // Nucleus sampling
    TopK:        40,     // Top-K sampling
}
```

## 🐳 Despliegue con Docker

### Construir imagen

```bash
docker build -t gemini-claude-mcp:latest .
```

### Ejecutar contenedor

```bash
docker run -p 8080:8080 \
  -e GEMINI_API_KEY="tu_api_key" \
  -e GEMINI_MODEL="gemini-1.5-flash" \
  gemini-claude-mcp:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  gemini-mcp:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - GEMINI_MODEL=gemini-1.5-flash
      - PORT=8080
    restart: unless-stopped
```

## 🧪 Ejemplos de Uso

### Análisis de Código Python

```bash
curl -X POST http://localhost:8080/analyze-code \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "async def fetch_data(url):\n    async with aiohttp.ClientSession() as session:\n        async with session.get(url) as response:\n            return await response.json()",
    "language": "python",
    "task": "review",
    "code_context": "Función asíncrona para obtener datos de API"
  }'
```

### Generación de Historia

```bash
curl -X POST http://localhost:8080/create-content \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Un robot que aprende a sentir emociones",
    "content_type": "story",
    "style": "emotivo y reflexivo",
    "length": "short",
    "temperature": 0.8
  }'
```

### Consulta Técnica

```bash
curl -X POST http://localhost:8080/ask-gemini \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Explica las diferencias entre REST y GraphQL",
    "system_role": "Eres un arquitecto de software senior",
    "temperature": 0.5
  }'
```

## 🔧 Solución de Problemas

### Error: "GEMINI_API_KEY no está configurada"

```bash
# Verificar que la variable esté configurada
echo $GEMINI_API_KEY

# Si está vacía, configurarla
export GEMINI_API_KEY="tu_api_key_real"
```

### Error: "Puerto ya en uso"

```bash
# Cambiar el puerto
export PORT="8081"

# O matar el proceso que usa el puerto
lsof -ti:8080 | xargs kill -9
```

### Error de conexión con Gemini

1. Verificar que tu API key sea válida
2. Confirmar que tienes cuota disponible en Gemini
3. Revisar logs del servidor para más detalles

### Claude Desktop no reconoce el servidor

1. Verificar la ruta en `claude_desktop_config.json`
2. Reiniciar Claude Desktop después de cambios
3. Revisar logs en la consola de Claude Desktop

## 📈 Optimizaciones de Rendimiento

### Para Alto Volumen

```go
// Aumentar límites si es necesario
config.MaxTokens = 4096
config.Temperature = 0.5  // Menos aleatorio = más rápido
```

### Para Respuestas Creativas

```go
config.Temperature = 0.9  // Más creativo
config.TopP = 0.95        // Mayor diversidad
```

### Para Análisis Técnico

```go
config.Temperature = 0.3  // Más determinista
config.MaxTokens = 2048   // Respuestas focalizadas
```

## 🛡️ Seguridad

- **API Key**: Nunca hardcodear en el código, usar variables de entorno
- **CORS**: Configurado para permitir origen desde Claude Desktop
- **Filtros de Seguridad**: Configurados para bloquear contenido inapropiado
- **Usuario no-root**: El Docker usa usuario sin privilegios

## 🤝 Contribuir

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Changelog

### v2.0.0 (2024-01-15)
- ✨ Soporte para múltiples tipos de consultas especializadas
- 🚀 Optimización para Claude Desktop MCP
- 🧠 Actualización a modelos Gemini más recientes
- 📊 Métricas y logging mejorados
- 🛡️ Mejoras de seguridad y manejo de errores

### v1.0.0 (2024-01-01)
- 🎉 Release inicial
- 💬 Consultas básicas a Gemini
- 🔌 Integración HTTP simple

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 📞 Soporte

- **Issues**: [GitHub Issues](https://github.com/scopweb/gemini-go-server-mcp/issues)
- **Documentación**: Este README y comentarios en el código
- **API de Gemini**: [Documentación oficial](https://ai.google.dev/docs)

---

**¡Disfruta conectando Claude Desktop con Gemini AI! 🚀🤖**