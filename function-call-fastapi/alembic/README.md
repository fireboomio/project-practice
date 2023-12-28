# Alembic

Alembic es una herramienta de migración de bases de datos para Python. Permite a los desarrolladores de software administrar y versionar esquemas de bases de datos en aplicaciones Python.

## Uso básico

Para usar Alembic, primero debes instalarlo en tu entorno de Python. Luego, debes crear una configuración de Alembic para tu proyecto. Puedes hacerlo ejecutando el siguiente comando en tu terminal:

``` alembic init alembic ```

Este comando creará una estructura de archivos para Alembic en tu proyecto. Luego, puedes crear una migración inicial ejecutando el siguiente comando:

``` alembic revision --autogenerate -m "Migración inicial" ```

Este comando creará una migración de base de datos para tu proyecto. Luego, puedes aplicar la migración ejecutando el siguiente comando:

``` alembic upgrade head ```

Este comando actualizará tu base de datos a la última versión de la migración.
