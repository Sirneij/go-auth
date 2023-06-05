# go-auth-backend

This powers the backend of the application.

## Running locally

To run this application locally, kindly follow the instructions below:

### Step 1: Install some external requirements

This app uses PostgreSQL and Redis. You need to install them. Then, create a database for the application in PostgreSQL. After that, you need to put your database's URL in `.env` at the root of the application. Your database URL should be like:

```shell
DATABASE_URL=postgres://<database_username>:<database_password>@localhost:<database_port>/<database_name>?sslmode=disable
```

Replace each variable with a real value.

You also need to set `REDIS_URL`.

Generally your `.env` file should look like:

```shell
DATABASE_URL=
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_MAX_IDLE_TIME=15m
REDIS_URL=
TOKEN_EXPIRATION=15m
SESSION_EXPIRATION=60m
FRONTEND_URL=http://localhost:3000
HMC_SECRET_KEY=13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b
PORT=8080
DEBUG=true
EMAIL_HOST_SERVER=
EMAIL_SERVER_PORT=2525
EMAIL_USERNAME=
EMAIL_PASSWORD=
AWS_S3_BUCKET_NAME=
AWS_SECRET_ACCESS_KEY=
AWS_ACCESS_KEY_ID=
AWS_REGION=
```

You are at liberty to change the supplied values.

### Step 2: Migrate the database

Have set those up, you need to migrate the database using [golang-migrate][1].

```shell
migrate -path=./migrations -database=<DATABASE_URL> up
```

### Step 3: Run the app

To run, just issue the following command:

```shell
make run/api
```

[1]: https://github.com/golang-migrate/migrate "Golang migrate"
