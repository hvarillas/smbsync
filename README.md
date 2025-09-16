# Go SMB Sync Tool

Esta es una herramienta de línea de comandos (CLI) para sincronizar archivos con recursos compartidos SMB. Su principal objetivo es copiar archivos de forma segura, verificar su integridad y, opcionalmente, comprimirlos y eliminarlos del origen.

## Características

- **Verificación de Integridad:** Garantiza que los archivos no se corrompan durante la transferencia calculando y comparando hashes SHA256 del archivo de origen y destino.
- **Compresión:** Permite comprimir archivos en formato `.zip` antes de transferirlos para ahorrar ancho de banda y espacio.
- **Borrado Seguro:** Opción para eliminar el archivo local solo después de una copia y verificación exitosas.
- **Logging Avanzado:** Utiliza `zap` para logs estructurados y `lipgloss` para una salida en color, configurable mediante la variable de entorno `LOG_LEVEL`.
- **Configuración Flexible:** Admite configuración mediante flags, variables de entorno o un archivo `.env`.
- **Encriptación de Contraseñas:** Soporte para contraseñas encriptadas con AES-GCM.
- **Encriptación de Strings:** Permite encriptar cualquier texto usando AES-GCM con clave personalizable.
- **Notificaciones Telegram:** Alertas automáticas por errores críticos.

## Estructura del Proyecto

```
smbsync/
├── cmd/smbsync/           # Punto de entrada de la aplicación
├── internal/              # Paquetes internos
│   ├── config/           # Configuración y flags
│   ├── crypto/           # Encriptación/desencriptación
│   ├── logger/           # Sistema de logging
│   ├── notification/     # Notificaciones Telegram
│   └── smb/             # Cliente SMB y operaciones
├── pkg/banner/           # Banner de la aplicación
├── testdata/            # Datos de prueba
├── .env.example         # Ejemplo de configuración
└── Makefile            # Comandos de build
```

## Instalación

1. Clona el repositorio:
   ```bash
   git clone https://github.com/hvarillas/smbsync
   cd smbsync
   ```

2. Instala las dependencias:
   ```bash
   go mod tidy
   ```

3. Compila la aplicación:
   ```bash
   make build
   ```

## Configuración

1. Copia el archivo de ejemplo:
   ```bash
   cp .env.example .env
   ```

2. Edita `.env` con tus credenciales:
   ```bash
   SMB_HOST=192.168.1.100
   SMB_USER=miusuario
   SMB_PASS=micontraseña
   SMB_SHARED=backups
   # ... más configuraciones
   ```

## Uso

### Flags Obligatorios

- `--user` o `-u`: El nombre de usuario para la autenticación SMB.
- `--pass` o `-p`: La contraseña para el usuario (o `--encrypted-pass`).
- `--host`: La dirección IP o el nombre de host del servidor SMB.
- `--shared` o `-s`: El nombre del recurso compartido que se va a montar.

### Flags Opcionales

- `--path`: Directorio local donde se encuentran los archivos a copiar. Por defecto, es el directorio actual (`.`).
- `--sharedPath`: La ruta relativa dentro del recurso compartido donde se copiarán los archivos. Por defecto, es la raíz (`.`).
- `--regex` o `-r`: Una expresión regular para filtrar los archivos a copiar.
- `--delete` o `-d`: Elimina el archivo local después de una copia y verificación exitosas.
- `--zip` o `-z`: Comprime cada archivo en un `.zip` individual antes de transferirlo.
- `--encrypted-pass`: Usar contraseña encriptada en lugar de texto plano.
- `--generate-encrypted`: Generar contraseña encriptada desde `--pass`.
- `--encrypt-text`: Encriptar cualquier string usando AES-GCM.
- `--encryption-key`: Clave de encriptación de 16 bytes (sobrescribe variable de entorno).

### Ejemplos de Ejecución

1. **Copiar todos los archivos `.bak`**:
   ```bash
   ./smbsync -u miusuario -p micontraseña --host 192.168.1.100 -s backups -r "\.bak$"
   ```

2. **Usar contraseña encriptada**:
   ```bash
   # Generar contraseña encriptada
   ./smbsync --generate-encrypted --pass "micontraseña"
   
   # Usar contraseña encriptada
   ./smbsync -u user --encrypted-pass "base64_encrypted_pass" --host host -s share
   ```

3. **Encriptar cualquier texto**:
   ```bash
   # Con clave personalizada
   ./smbsync --encrypt-text "mi texto secreto" --encryption-key "1234567890123456"
   
   # Con clave por defecto o variable de entorno
   ./smbsync --encrypt-text "datos sensibles"
   ```

4. **Copiar y comprimir archivos con eliminación**:
   ```bash
   ./smbsync -u user -p pass --host host -s share -r "\.log$" --zip --delete
   ```

## Variables de Entorno

Configura valores por defecto vía variables de entorno:

```bash
export SMB_HOST=192.168.1.100
export SMB_USER=miusuario
export SMB_PASS=micontraseña
export SMB_SHARED=backups
export LOG_LEVEL=info
export TELEGRAM_BOT_TOKEN=tu_token
export TELEGRAM_CHAT_ID=tu_chat_id
export ENCRYPTION_KEY=tu_clave_16_bytes
```

## Comandos de Build

```bash
# Compilar para la plataforma actual
make build

# Compilar para Windows
make build-windows

# Compilar para Linux
make build-linux

# Ejecutar tests
make test

# Limpiar artefactos
make clean

# Formatear código
make fmt
```

## Seguridad

- **Nunca** subas el archivo `.env` con credenciales reales al control de versiones.
- Usa contraseñas encriptadas en producción con `--generate-encrypted`.
- Configura las variables de entorno `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID` y `ENCRYPTION_KEY` para mayor seguridad.
- La clave de encriptación debe tener exactamente 16 bytes para AES-128.
- La herramienta verifica automáticamente la integridad de cada archivo con SHA256.

## Notas

- La herramienta crea automáticamente los directorios remotos si no existen.
- Todos los logs se escriben tanto a archivo como a consola.
- Las notificaciones Telegram se envían solo para errores críticos.
- La verificación de integridad es obligatoria para todos los archivos transferidos.
- La funcionalidad de encriptación permite proteger cualquier string sensible, no solo contraseñas.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT. Consulta el archivo [LICENSE](LICENSE) para más detalles.
