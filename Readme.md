# Layer8 Authentication Server

### Project structure

```md
layer8-auth-server
├── **`cmd`: entry points**
│   └── `server`
│       └── `main.go`: entry point, wire dependencies and start server.
├── **`internal`: Core application logic; contains code that is private to this project and cannot be imported by external applications.**
│   ├── `config` (pkg): Definitions for configuration objects and logic for loading them from environment variables or files.
│   ├── `consts` (pkg): Centralized repository for global constants and fixed values used across the application.
│   ├── `dto` (pkg): Data Transfer Objects used to define the schema for data entering and leaving the API layer. 
│   │   ├── `requestdto` (pkg): Structs representing incoming request payloads and validation rules (future update).
│   │   └── `responsedto` (pkg): Structs representing outgoing response payloads to ensure a consistent API format.
│   ├── `models`: Definitions for domain entities and persistence structures.
│   │   ├── `gormModels` (pkg): database (postgres + gorm) models
│   │   ├── `email.go`: Domain models related to email structures and metadata.
│   │   ├── `jwt.go`: Models for handling JWT claims and authentication tokens.
│   │   └── `statistics.go`: Domain models representing client traffic usage, quotas, and resource consumption metrics.
│   ├── **`handlers`: API handlers (controllers)**
│   │   ├── `clientH` (pkg): Client portal APIs handlers.
│   │   ├── `oauthH` (pkg): Oauth APIs handlers.
│   │   └── `userH` (pkg): User portal APIs handlers.
│   ├── **`usecases`: Orchestration Layer - The central business logic layer, contains business logic and coordinates data flow between handlers and repositories.**
│   │   ├── `ucerror` (pkg): Defines custom error types returned by use case functions, allowing the handler to translate business failures into appropriate HTTP status codes. 
│   │   ├── `clientUC` (pkg): Handles business logic for the Client Portal.
│   │   ├── `oauthUC` (pkg): Handles oauth2 related apis logic.
│   │   ├── `userUC` (pkg): Handles business logic for the User Portal.
│   │   └── `workerUC` (pkg): Contains background task logic and asynchronous functions triggered directly by the system rather than HTTP requests.
│   └── **`repositories`: Data Access Layer**
│       ├── `codeGenRepo` (pkg): Logic for generating and managing short-lived verification codes (OTP) for phone or email.
│       ├── `emailRepo` (pkg): Interfaces with external SMTP servers or email providers to handle message delivery and verification.
│       ├── `ethRepo` (pkg): Handles interactions with the Ethereum blockchain, such as smart contract calls or transaction monitoring.
│       ├── `influxdbRepo` (pkg): Manages time-series data storage and retrieval specifically for InfluxDB.
│       ├── `phoneRepo` (pkg): Handles interactions with SMS gateways or phone verification services.
│       ├── `postgresRepo` (pkg): Implements persistence logic and relational queries using GORM and PostgreSQL.
│       ├── `tokenRepo` (pkg): Generate and verify tokens.
│       └── `zkRepo` (pkg): Interfaces with Zero-Knowledge proof providers or local ZK circuit verification logic.
├── **`pkg`: Shared Utility Libraries**
│   ├── `code` (pkg): Reusable logic for generating and formatting verification codes (OTP) for email and phone.
│   ├── `eth` (pkg): Low-level Ethereum client wrappers and blockchain primitive helpers.
│   ├── `ginUtils` (pkg): Helper functions and middleware to extend the Gin web framework (e.g., custom error responders).
│   ├── `log` (pkg): A standardized logging wrapper to ensure consistent output formats across the entire app.
│   ├── `oauth` (pkg): Shared OAuth2 primitives and helper structs for handling third-party provider data.
│   ├── `scram` (pkg): Implementation of Salted Challenge Response Authentication Mechanism (SCRAM) for secure password hashing.
│   ├── `telegram` (pkg): Client wrapper for interacting with the Telegram Bot API for notifications or commands.
│   ├── `utils` (pkg): General-purpose helper functions (string manipulation, date formatting, slice helpers, etc.).
│   └── `zk` (pkg): Reusable Zero-Knowledge cryptographic primitives and mathematical helper functions.
├── `smart-contract`: Contains Solidity source code, interfaces, and deployment scripts (e.g., Hardhat) for blockchain integration.
└── **`web`: The frontend application (e.g. Vue); handles the user interface and interacts with the Go API and smart contracts.**
```