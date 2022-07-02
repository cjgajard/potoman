# potoman

Herramienta para hacer _requests_ HTTP de manera simplificada, sin tener que
cambiar las costumbres que se adquiren por usar `bash`.

## Instalación

Requiere golang >= 1.16

```
go install github.com/cjgajard/potoman@latest
```

Alternativamente, puedes compilar localmente:

```
cd /path/to/potoman
go install .
```

## Ejemplo de uso

1. Escribe un archivo con el _request_ HTTP que desees, donde quieras y con el
   nombre que más te guste:

`./get-momma-joke.txt`
```
GET api.yomomma.info
```

2. Ejecuta potoman

```
potoman get-momma-joke.txt
```

3. Recibe tu chiste en `stdout`

```
{
    "joke": "Yo mamma so fat she leaves footsteps in concrete"
}
```

Me dirás ¡pErO wN, eSo eN cUrL eS mUcHo mÁs fÁcIl! y te respondo con otro
ejemplo:

---

1. Escribe un archivo con el _request_ HTTP que desees, donde quieras y con el
   nombre que más te guste:

`./fake-endpoint.txt`
```
POST https://${API_DOMAIN}/register HTTP/3.0
X-API-Key: $API_KEY
Content-Type: application/json

{
  "email":${EMAIL:-cringey_username_created@tenyears.old},
  "password":$PASSWORD
}
```

2. Escribe un archivo `.env` con los datos que quieras cargar a través de
   el entorno de tu terminal.

```
API_DOMAIN=successful.startup.com
API_KEY=e4ea3a6694b74b5b8e0034bc98c12ec3
```

2. Ejecuta potoman

```
potoman fake-endpoint.txt
```

¡MIERDA, Olvidé agregar `EMAIL` y `PASSWORD` al `.env`!

Tranquilo, `potoman` es a prueba de weones y lo rellenará porque pusiste un
`${VALOR:-pordefecto}` usando la misma sintáxis en `bash`. Mientras que te
preguntará por el resto de variables sin definir :

>
> Please enter value of $PASSWORD, parameter is empty: 
>

Tú respondes escribiendo un valor, confirmándolo con ENTER, tu _request_ se
ejecutará correctamente y todos seremos felices :D

3. Recibe tu respuesta en `stdout`

```
{
  "status": 410,
  "error": "unable to create a new user",
  "reason": "We're sorry our eco-friendly company went bankrupt"
}
```
