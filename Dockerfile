# Dockerfile optimizado para el servidor MCP Gemini-Claude
FROM golang:1.24.2-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git ca-certificates

# Crear directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download && go mod verify

# Copiar código fuente
COPY main.go ./

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server main.go

# Imagen final minimalista
FROM alpine:latest

# Instalar certificados SSL para llamadas HTTPS
RUN apk --no-cache add ca-certificates

# Crear usuario no-root para seguridad
RUN addgroup -g 1001 appgroup && \
    adduser -D -u 1001 -G appgroup appuser

# Crear directorio de trabajo
WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/server .

# Cambiar propietario
RUN chown -R appuser:appgroup /app

# Cambiar a usuario no-root
USER appuser

# Exponer puerto
EXPOSE 3000

# Variables de entorno por defecto
ENV PORT=3000
ENV GEMINI_MODEL=gemini-1.5-flash

# Comando para ejecutar la aplicación
CMD ["./server"]