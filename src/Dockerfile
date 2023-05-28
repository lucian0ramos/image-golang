# Imagen base
FROM golang:1.16-alpine

# Establecer el directorio de trabajo
WORKDIR /src

# Copiar los archivos necesarios
COPY go.mod go.sum ./
COPY . .

# Descargar las dependencias del módulo Go
RUN go mod download

# Compilar la aplicación
RUN go build -o main .

# Comando de inicio de la aplicación
CMD ["./main"]
