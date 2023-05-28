# image-golang
- ✅ Imagen Docker funcionando
- ✅ Test unitarios
- ✅ Validación de parámetros
- ✅ Middlewares
- ✅ MySQL
- ⏳ Mejorando README
- ⏳ Arreglando MySQL Image

# Instrucciones para levantar la imagen de src de api
- Ejecutar cmd:
```ENV=dev IMAGE_TAG=l REGISTRY=l docker-compose up --build src```

# Instrucciones para probar test unitarios
- Situarse en la carpeta test:
```cd src/test```
- Ejecutar cmd:
```go test```

# Instrucciones para arrancar cualquier fichero golang en local
- Ejecutar cmd:
``` go run <nombredelfichero> ```

# Instrucciones para llamar al endpoint
- Una vez arrancada la imagen de docker en local, ejecutar cmd para llamar al endpoint:
``` curl --request POST \ --url http://localhost:8080/manage-errors \ --header 'Authorization: Basic d459vd29zOnc9MHcwasdxPXo=' \ --header 'Content-Type: application/json' \ --data '{  "deploy_id": 1, "app_name": "test", "type": "ci/release", "version": "1.0.0", "values": {"key": "value"}}' ```

