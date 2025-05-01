# playground examples

## Using external Sqlite database

1. Identify the location of the external Sqlite database.

   Example:

    ```console
    export DEMO_DATABASE_DIR=~/demo-databases/with-config
    ```

1. Optionally, download a Sqlite database (e.g. `sz.db`) containing a Senzing database schema
   and an initial Senzing Configuration.

   Example:

    ```console
    mkdir -p ${DEMO_DATABASE_DIR}
    curl --output ${DEMO_DATABASE_DIR}/sz.db https://raw.githubusercontent.com/senzing-garage/demo-databases/refs/heads/main/sqlite/with-config/sz.db
    ```

1. Create a `SENZING_ENGINE_CONFIGURATION_JSON` that will reference the external database.
   **Remember:**  The path in `SQL.CONNECTION` must be relative to the *inside* of the Docker container.

   Example:

    ```console
    export SENZING_ENGINE_CONFIGURATION_JSON='
    {
        "PIPELINE": {
            "CONFIGPATH": "/etc/opt/senzing",
            "LICENSESTRINGBASE64": "",
            "RESOURCEPATH": "/opt/senzing/er/resources",
            "SUPPORTPATH": "/opt/senzing/data"
        },
        "SQL": {
            "CONNECTION": "sqlite3://na:na@/data/sz.db"
        }
    }
    '
    ```

1. Run the Docker container.

   Example:

    ```console
    docker run \
    --env SENZING_ENGINE_CONFIGURATION_JSON \
    --interactive \
    --name senzing-playground \
    --publish 8260:8260 \
    --publish 8261:8261 \
    --pull always \
    --rm \
    --tty \
    --volume ${DEMO_DATABASE_DIR}:/data \
    senzing/playground
    ```

1. Visit the web application at [localhost:8260].
1. Optionally, enter the docker container.

   Example:

    ```console
    docker exec -it senzing-playground /bin/bash
    ```

[localhost:8260]: http://localhost:8260
