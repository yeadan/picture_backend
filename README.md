# picture_backend
__Api para una red social de imágenes, con comentarios y likes, parte de servidor__

## Utiliza:
  "[github.com/urfave/negroni](https://github.com/urfave/negroni)" - Middleware para net/http (logs, panics recovery y static en "public")  
  "[github.com/jinzhu/gorm](https://github.com/jinzhu/gorm)" - ORM para Postgres  
  "[github.com/gorilla/mux](https://github.com/gorilla/mux)" - Router  
  "[github.com/patrickmn/go-cache](https://github.com/patrickmn/go-cache)" - cache de llave/valor utilizada para almacenar usuarios y tokens HMAC  
  "[github.com/asaskevich/govalidator](https://github.com/asaskevich/govalidator)" - Validador de estructuras, utiizado solo para los users  
  "[github.com/graux/image-manager](https://github.com/graux/image-manager)" Mini librería para preparar imágenes en JPEG
   "[github.com/joho/godotenv](https://github.com/joho/godotenv)" - Load .env variable
  

__.env utilizado para curso backend PalmaActiva__
secret_key=PalmaActiva2019
DB_HOST=127.0.0.1
DB_DRIVER=postgres
DB_USER=sergio
DB_PASSWORD=bosch
DB_NAME=image_api
DB_PORT=5432 #Default postgres port

## Estructura de la api

### Users
Perfiles de usuario con roles (anonymous/user/admin)
Sus rutas/métodos son las siguientes:

"/users"              - post - Registro de usuarios, encripta los passwords con SHA256
"/users/login"        - post - método post - Login de usuario
"/users"              - get  - Listado de todos los usuarios. Solo un administrador
"/users/{id:[0-9]+}"  - get  - Detalles de un usuario en concreto. 
"/users/{id:[0-9]+"   - put  - Editar la información del usuario. No se puede cambiar ID ni username. El role solo lo puede cambiar un admin

### Images/pictures
