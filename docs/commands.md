Migrate Commands:

To Create the migration file

    migrate create -ext sql -dir db/migrations -seq create_users_table

To execute migrations:

    export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/postgres?sslmode=disable'
    
    migrate -database ${POSTGRESQL_URL} -path db/migrations up
    