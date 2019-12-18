## Explicación de pruebas hechas con el runner de [POSTMAN](https://www.getpostman.com/):  
__1-test role user creacion-edicion de todo__: Registro, login y creación/edición de todo el contenido que puede hacer un usuario con el role de "user". Sin fallos  

__2-test usuario anonymous__: Registro, login y pruebas sobre el distinto contenido. Falla en todo lo que no tiene permiso. Solo pasan los test en las 3 cosas que puede hacer un usuario de role "anonymous"  

__3-user admin__: Registro, login y pruebas de usuario con role "admin". Añadidos dos fallos para comprobar los 404 en crear/borrar likes y añadir/quitar fotos que no existen en una galería. También falla la prueba de borrar una galería que no es suya (Única cosa que no debe ni puede borrar un admin)  

__4-user con fallos__: Prueba de fallos con usuario "user". Los cuatro fallos consisten en:  
1. Intentar editar su password con menos letras de las permitidas -> __BadRequest__  
2. Crear un comentario con un JSON mal construido -> __BadRequest__  
3. Crear una imagen con un archivo que no es una tipo imagen -> __BadRequest__  
4. Borrar una galería que no es suya -> __Forbidden__  

Por último, en la carpeta "__Collection empleada__" está el export con todas las pruebas utilizadas para testear la api  

__Actualización 18-12-2019__: Modificado el código del signup para que devuelva un __201__ (StatusCreated) en vez de un __200__ (StatusOK), así que debería fallar ese test al crear un usuario.
