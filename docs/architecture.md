# 📊 Flutter Project Architecture Documentation (Clean Architecture + BLoC)

The architecture in this project strictly follows the principles of **Clean Architecture**:
## 📁 Project Structure

```plaintext
project/  
├── api/                 # API-related files, like OpenAPI, protocol definitions  
│ ├── openapi/  
│ │ ├── healtcheck.yaml  
│ │ └── ...  
│ └── ...  
├── cmd/                 # Main applications for this project.  
│ ├── server/  
│ │ ├── main.go          # Application entry point  
│ │ └── ...  
│ └── cron/  
│ ├── main.go            # Another application entry point  
│ └── ...  
├── internal/            # Private application and package code  
│ ├── handler/  
│ │── service/  
│ ├── database/  
│ └── ...   
├── scripts/             # Build, deployment, and maintenance scripts  
│ ├── build.sh  
│ ├── deploy.sh  
│ └── ...  
├── configs/             # Configuration files for different environments  
│ ├── development.yaml  
│ ├── production.yaml  
│ └── ...  
├── tests/               # Some additional tests and test data  
│ ├── integration/  
│ │ ├── ...  
│ └── testdata/  
│ └── ...  
├── docs/                # Project documentation  
├── .gitignore           # Gitignore file  
├── go.mod               # Go module file  
├── go.sum               # Go module dependencies file  
└── README.md            # Project README
```