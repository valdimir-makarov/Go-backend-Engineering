# Diff Details

Date : 2025-08-29 17:37:36

Directory /home/bubun/Go-backend-Engineering/chat-service/Auth-Service

Total : 55 files,  -763 codes, -55 comments, -163 blanks, all -981 lines

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [Auth-Service/.env](/Auth-Service/.env) | Properties | 5 | 0 | 0 | 5 |
| [Auth-Service/Dockerfile](/Auth-Service/Dockerfile) | Docker | 11 | 9 | 11 | 31 |
| [Auth-Service/config/config.go](/Auth-Service/config/config.go) | Go | 54 | 0 | 9 | 63 |
| [Auth-Service/controllers/auth\_controller.go](/Auth-Service/controllers/auth_controller.go) | Go | 121 | 6 | 23 | 150 |
| [Auth-Service/controllers/user\_controller.go](/Auth-Service/controllers/user_controller.go) | Go | 1 | 0 | 1 | 2 |
| [Auth-Service/db/migrations/20250716210945\_create\_users\_table.down.sql](/Auth-Service/db/migrations/20250716210945_create_users_table.down.sql) | MS SQL | 1 | 0 | 0 | 1 |
| [Auth-Service/db/migrations/20250716210945\_create\_users\_table.up.sql](/Auth-Service/db/migrations/20250716210945_create_users_table.up.sql) | MS SQL | 12 | 0 | 0 | 12 |
| [Auth-Service/db/migrations/20250716211556\_adding\_Role\_col.down.sql](/Auth-Service/db/migrations/20250716211556_adding_Role_col.down.sql) | MS SQL | 1 | 0 | 0 | 1 |
| [Auth-Service/db/migrations/20250716211556\_adding\_Role\_col.up.sql](/Auth-Service/db/migrations/20250716211556_adding_Role_col.up.sql) | MS SQL | 1 | 0 | 0 | 1 |
| [Auth-Service/db/migrations/20250716211741\_drop\_users\_table.down.sql](/Auth-Service/db/migrations/20250716211741_drop_users_table.down.sql) | MS SQL | 0 | 0 | 1 | 1 |
| [Auth-Service/db/migrations/20250716211741\_drop\_users\_table.up.sql](/Auth-Service/db/migrations/20250716211741_drop_users_table.up.sql) | MS SQL | 0 | 0 | 1 | 1 |
| [Auth-Service/docker-compose.yml](/Auth-Service/docker-compose.yml) | YAML | 20 | 44 | 8 | 72 |
| [Auth-Service/go.mod](/Auth-Service/go.mod) | Go Module File | 56 | 0 | 6 | 62 |
| [Auth-Service/go.sum](/Auth-Service/go.sum) | Go Checksum File | 177 | 0 | 1 | 178 |
| [Auth-Service/kafka/Kafka\_prod.go](/Auth-Service/kafka/Kafka_prod.go) | Go | 37 | 0 | 11 | 48 |
| [Auth-Service/krakend.json](/Auth-Service/krakend.json) | JSON | 0 | 69 | 1 | 70 |
| [Auth-Service/main.go](/Auth-Service/main.go) | Go | 49 | 10 | 14 | 73 |
| [Auth-Service/makefile](/Auth-Service/makefile) | Makefile | 29 | 5 | 9 | 43 |
| [Auth-Service/middleware/auth\_middleware.go](/Auth-Service/middleware/auth_middleware.go) | Go | 31 | 2 | 8 | 41 |
| [Auth-Service/middleware/jwt.go](/Auth-Service/middleware/jwt.go) | Go | 35 | 1 | 8 | 44 |
| [Auth-Service/models/role.go](/Auth-Service/models/role.go) | Go | 1 | 0 | 1 | 2 |
| [Auth-Service/models/user.go](/Auth-Service/models/user.go) | Go | 16 | 4 | 5 | 25 |
| [Auth-Service/repository/repo.go](/Auth-Service/repository/repo.go) | Go | 70 | 10 | 13 | 93 |
| [Auth-Service/routes/routes.go](/Auth-Service/routes/routes.go) | Go | 11 | 1 | 3 | 15 |
| [Auth-Service/utils/jwt.go](/Auth-Service/utils/jwt.go) | Go | 35 | 2 | 12 | 49 |
| [Auth-Service/utils/password.go](/Auth-Service/utils/password.go) | Go | 10 | 0 | 4 | 14 |
| [chat-service/Dockerfile](/chat-service/Dockerfile) | Docker | -13 | -3 | -9 | -25 |
| [chat-service/README.md](/chat-service/README.md) | Markdown | -36 | 0 | -8 | -44 |
| [chat-service/cmd/main.go](/chat-service/cmd/main.go) | Go | -66 | -5 | -17 | -88 |
| [chat-service/config/config.go](/chat-service/config/config.go) | Go | -18 | -2 | -5 | -25 |
| [chat-service/database/migration/000001\_init\_schema.down.sql](/chat-service/database/migration/000001_init_schema.down.sql) | MS SQL | -3 | -1 | 0 | -4 |
| [chat-service/database/migration/000001\_init\_schema.up.sql](/chat-service/database/migration/000001_init_schema.up.sql) | MS SQL | -49 | -5 | -7 | -61 |
| [chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/chat.pb.go](/chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/chat.pb.go) | Go | -174 | -9 | -28 | -211 |
| [chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/chat\_grpc.pb.go](/chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/chat_grpc.pb.go) | Go | -77 | -29 | -16 | -122 |
| [chat-service/go.mod](/chat-service/go.mod) | Go Module File | -27 | 0 | -5 | -32 |
| [chat-service/go.sum](/chat-service/go.sum) | Go Checksum File | -104 | 0 | -1 | -105 |
| [chat-service/init-db.sql](/chat-service/init-db.sql) | MS SQL | -8 | 0 | -1 | -9 |
| [chat-service/internal/delivery/FIleUplaod.handler.go](/chat-service/internal/delivery/FIleUplaod.handler.go) | Go | -46 | -1 | -12 | -59 |
| [chat-service/internal/delivery/chat\_handler.go](/chat-service/internal/delivery/chat_handler.go) | Go | -289 | -26 | -53 | -368 |
| [chat-service/internal/delivery/grpc\_chat\_handler.go](/chat-service/internal/delivery/grpc_chat_handler.go) | Go | -14 | -10 | -8 | -32 |
| [chat-service/internal/kafka/Auth\_kafka/Kafka\_prod.go](/chat-service/internal/kafka/Auth_kafka/Kafka_prod.go) | Go | -41 | -1 | -9 | -51 |
| [chat-service/internal/kafka/kafka\_con.go](/chat-service/internal/kafka/kafka_con.go) | Go | -34 | 0 | -7 | -41 |
| [chat-service/internal/kafka/kafka\_prod.go](/chat-service/internal/kafka/kafka_prod.go) | Go | -63 | 0 | -11 | -74 |
| [chat-service/internal/models/chat.go](/chat-service/internal/models/chat.go) | Go | -26 | -8 | -5 | -39 |
| [chat-service/internal/models/user.go](/chat-service/internal/models/user.go) | Go | -13 | 0 | -2 | -15 |
| [chat-service/internal/repository/chat\_repository.go](/chat-service/internal/repository/chat_repository.go) | Go | -255 | -15 | -37 | -307 |
| [chat-service/internal/repository/message\_repository.go](/chat-service/internal/repository/message_repository.go) | Go | -1 | 0 | -1 | -2 |
| [chat-service/internal/service/chat\_service.go](/chat-service/internal/service/chat_service.go) | Go | -61 | -10 | -20 | -91 |
| [chat-service/internal/usecase/chat\_usecase.go](/chat-service/internal/usecase/chat_usecase.go) | Go | -1 | -58 | -12 | -71 |
| [chat-service/test/k6.js](/chat-service/test/k6.js) | JavaScript | -35 | -2 | -13 | -50 |
| [chat-service/test/test.go](/chat-service/test/test.go) | Go | -1 | -23 | -7 | -31 |
| [chat-service/test/test2.go](/chat-service/test/test2.go) | Go | -59 | -4 | -11 | -74 |
| [chat-service/tmp/build-errors.log](/chat-service/tmp/build-errors.log) | Log | -1 | 0 | 0 | -1 |
| [chat-service/utils/jwt\_utils.go](/chat-service/utils/jwt_utils.go) | Go | -1 | 0 | -1 | -2 |
| [chat-service/utils/logger.go](/chat-service/utils/logger.go) | Go | -31 | -6 | -7 | -44 |

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details