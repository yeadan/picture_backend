# picture_backend
__Api creada en [GO](https://golang.org/) para una red social de imágenes, con comentarios y likes. Parte de servidor__  
Directorio public con la __primera__ versión del Front no utilizado.

## Librerías externas utilizadas:
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

"/users"              - __post__ - Registro de usuarios, encripta los passwords con SHA256  
"/users/login"        - __post__  - Login de usuario  
"/users"              - __get__  - Listado de todos los usuarios. Exclusivo admin  
"/users/{id:[0-9]+}"  - __get__  - Detalles de un usuario en concreto.  
"/users/{id:[0-9]+"   - __put__  - Editar un usuario. No se puede cambiar ID ni username. El role solo lo puede cambiar un admin  

#### Roles  
__admin__: Tiene control sobre todo, menos para crear contenido como si fuera otro usuario. Puede borrar y editar lo que suben los demás.  
__user__: No puede listar todos los usuarios ni borrar o editar el contenido de los otros usuarios.  
__anonymous(u otros)__: Usuario de prueba. Solamente puede editar su propia información, ver todas las fotos y ver los detalles de una foto concreta.  

### Images/pictures  
Estructura de las imágenes, con sus metadatos (picture) y donde están almacenadas (image)

"/images"  - __post__ - Al pasarle un JPEG crea una imagen y una picture. Se almacenan en la carpeta ./images  
"/images" - __get__ - Listado ordenado de todas las imágenes. Las más recientes primero  
"/images/{id:[0-9]+}" - __put__ - Añade un título y una descripción a la imagen  
"/images/{id:[0-9]+}" - __get__ - Listado de las imágenes de un usuario  
"/image/{id:[0-9]+}" - __get__ - Ver los detalles de una imagen en concreto  
"/images/{id:[0-9]+}" - __delete__ - Borra una imagen (image y picture)  
"/images/avatar"  - __post__ - Crea una imagen para el avatar de un usuario. Se almacena en la carpeta ./images/avatars  

### Comment  
Comentarios de las imágenes. Una foto puede tener todos los comentarios que se quieran, pero si se borra la foto, se borran todos

"/comment" - __post__ - Crea un comentario en una imagen  
"/comment/{id:[0-9]+} - __put__ - Edita un comentario  
"/comment/{id:[0-9]+}" - __delete__ - Borra un comentario  

### Like
Likes de las imágenes. Solo un like por foto (del mismo usuario).  

"/like" - __post__ - Crea un like en una imagen  
"/like/{id:[0-9]+}" - __delete__ - Borra un like   

### Gallery/PicGallery  
Galerías de fotos de un usuario (gallery) y relación entre pictures/galleries (picgallery)  

"/gallery" - __post__ - Crea una galería de fotos vacía  
"/gallery/{id:[0-9]+}" - __delete__ - Borra una galería de fotos, junto a todos sus picgalleries  
"/gallery/{gal:[0-9]+}/{pic:[0-9]+}" - __delete__ - Quita una foto de dentro de una galería (no borra la foto)  
"/gallery/{gal:[0-9]+}/{pic:[0-9]+}" - __post__ - Mete una foto dentro de una galería  
"/gallery/{id:[0-9]+}"  __put__ - Edita el título de una galería  
"/gallery/{id:[0-9]+}"  __get__ - Lista todas las fotos de una galería  
