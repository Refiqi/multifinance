
### Prerequisite
  - Make sure You already have Golang ,Gin, PostgresqlDB Installed
  - Create `multifinance` Database in your Postgresql

### Run Local

    git clone <url repo> | Clone this repo
    Open the folder that contains the file that have been cloned
    Set you Environment with key `ENV_MODE` to `stg`. on windows you set the env by running this in CMD `set ENV_MODE=stg`. on Linux you can simply set the env when you run the server `ENV_MODE=stg go run main.go`
    Run `go mod tidy && go run main.go`
    The server Will run at localhost:8888
    
### Application Architecture Diagram
![Application_Architecture_Diagram](https://github.com/user-attachments/assets/78cf5041-e173-4f55-a3c4-86cccc50841e)

### ERD Database
![ERD multifinance](https://github.com/user-attachments/assets/ad52a329-3110-4337-83d7-f6f00bdc087b)
