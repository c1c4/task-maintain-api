# Task Maintain API
This api has the objective to create tasks, and list and notify the Managers about them.

#
## The Api
The api is built with: 

### API
1.  [*Golang*](https://go.dev/)
1.  [*Gin*](https://github.com/gin-gonic/gin)
1.  [*Google Pub/Sub*](https://cloud.google.com/pubsub/docs/overview)
1.  [*Checkmail*](https://github.com/badoux/checkmail)
1.  [*JWT*](https://github.com/dgrijalva/jwt-go) is recommended to see the new [*repotisory*](https://github.com/golang-jwt/jwt)
1.  [*Godotenv*](https://github.com/joho/godotenv)

### DB
1.  [*Gorm*](https://gorm.io/index.html)

### Test
1.  [*Testify*](https://github.com/stretchr/testify)
1.  [*Go-sqlmock*](https://github.com/DATA-DOG/go-sqlmock)


#

## Running the API
You can run the task maintain api with two ways:

1.  Docker for this you will need docker installed in your machine [Docker](https://www.docker.com/)

        make build.application
    
    or if you already build

        make start.application

    This will put the api and database online.
    
    You can use:


        docker ps

    You should see an output that starts with something that looks like the following:
    
    CONTAINER ID | IMAGE
    ------------ | -----
    a123bc007edf | task-api
    867g5309hijk | task-db


    with this and I believe you should be fine and start to using the api in this url **localhost:8080**

1.  Open the project in you preferred IDE Goland or VSCode then run in your terminal:
    
        go mod download

    You can run:
        
        docker-compose up -d task-maintain-db

    To instantiate the Postgre or download and install in your local machine.

#

## Test the API
Well you can test the mutant api with two ways as well no surprises I hope:

1.  These are the commands you need and you see above
    
    For unit tests:

        make test.unit
    
    For integration tests:
        
        make test.integration
    
    The only thing new here is the **pytest -v** this will run all the test are in the test folder.

1.  You can open the project in PyCharm or VSCode or the IDE your like but I know these two has support for test and configure their tests
    1. [Golang on Goland](https://www.jetbrains.com/go/)
    1. [Golang on VSCode](https://code.visualstudio.com/docs/languages/go)

    Or you can open you terminal go to the project folder and run

        go test ./... -v <- this will run all tests

    You only need to make sure if you install the all the dependencies and has your .env setup.
