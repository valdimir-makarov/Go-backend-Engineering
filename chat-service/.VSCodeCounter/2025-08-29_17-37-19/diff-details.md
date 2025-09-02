# Diff Details

Date : 2025-08-29 17:37:19

Directory /home/bubun/Go-backend-Engineering/chat-service/chat-service

Total : 40 files,  -848 codes, 176 comments, 255 blanks, all -417 lines

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [Text-Editor-Service/.mocharc.json](/Text-Editor-Service/.mocharc.json) | JSON | -5 | 0 | 0 | -5 |
| [Text-Editor-Service/dist/server.d.ts](/Text-Editor-Service/dist/server.d.ts) | TypeScript | -1 | -1 | 0 | -2 |
| [Text-Editor-Service/dist/server.js](/Text-Editor-Service/dist/server.js) | JavaScript | -25 | -4 | 0 | -29 |
| [Text-Editor-Service/dist/type.d.ts](/Text-Editor-Service/dist/type.d.ts) | TypeScript | -9 | -1 | 0 | -10 |
| [Text-Editor-Service/dist/type.js](/Text-Editor-Service/dist/type.js) | JavaScript | -1 | -1 | 0 | -2 |
| [Text-Editor-Service/package-lock.json](/Text-Editor-Service/package-lock.json) | JSON | -2,165 | 0 | -1 | -2,166 |
| [Text-Editor-Service/package.json](/Text-Editor-Service/package.json) | JSON | -31 | 0 | -1 | -32 |
| [Text-Editor-Service/src/CRDT.ts](/Text-Editor-Service/src/CRDT.ts) | TypeScript | -54 | -20 | -17 | -91 |
| [Text-Editor-Service/src/server.ts](/Text-Editor-Service/src/server.ts) | TypeScript | -71 | -15 | -35 | -121 |
| [Text-Editor-Service/src/type.ts](/Text-Editor-Service/src/type.ts) | TypeScript | -11 | 0 | -3 | -14 |
| [Text-Editor-Service/tsconfig.json](/Text-Editor-Service/tsconfig.json) | JSON with Comments | -22 | 0 | -1 | -23 |
| [chat-service/Dockerfile](/chat-service/Dockerfile) | Docker | 13 | 3 | 9 | 25 |
| [chat-service/README.md](/chat-service/README.md) | Markdown | 36 | 0 | 8 | 44 |
| [chat-service/cmd/main.go](/chat-service/cmd/main.go) | Go | 66 | 5 | 17 | 88 |
| [chat-service/config/config.go](/chat-service/config/config.go) | Go | 18 | 2 | 5 | 25 |
| [chat-service/database/migration/000001\_init\_schema.down.sql](/chat-service/database/migration/000001_init_schema.down.sql) | MS SQL | 3 | 1 | 0 | 4 |
| [chat-service/database/migration/000001\_init\_schema.up.sql](/chat-service/database/migration/000001_init_schema.up.sql) | MS SQL | 49 | 5 | 7 | 61 |
| [chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/chat.pb.go](/chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/chat.pb.go) | Go | 174 | 9 | 28 | 211 |
| [chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/chat\_grpc.pb.go](/chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/chat_grpc.pb.go) | Go | 77 | 29 | 16 | 122 |
| [chat-service/go.mod](/chat-service/go.mod) | Go Module File | 27 | 0 | 5 | 32 |
| [chat-service/go.sum](/chat-service/go.sum) | Go Checksum File | 104 | 0 | 1 | 105 |
| [chat-service/init-db.sql](/chat-service/init-db.sql) | MS SQL | 8 | 0 | 1 | 9 |
| [chat-service/internal/delivery/FIleUplaod.handler.go](/chat-service/internal/delivery/FIleUplaod.handler.go) | Go | 46 | 1 | 12 | 59 |
| [chat-service/internal/delivery/chat\_handler.go](/chat-service/internal/delivery/chat_handler.go) | Go | 289 | 26 | 53 | 368 |
| [chat-service/internal/delivery/grpc\_chat\_handler.go](/chat-service/internal/delivery/grpc_chat_handler.go) | Go | 14 | 10 | 8 | 32 |
| [chat-service/internal/kafka/Auth\_kafka/Kafka\_prod.go](/chat-service/internal/kafka/Auth_kafka/Kafka_prod.go) | Go | 41 | 1 | 9 | 51 |
| [chat-service/internal/kafka/kafka\_con.go](/chat-service/internal/kafka/kafka_con.go) | Go | 34 | 0 | 7 | 41 |
| [chat-service/internal/kafka/kafka\_prod.go](/chat-service/internal/kafka/kafka_prod.go) | Go | 63 | 0 | 11 | 74 |
| [chat-service/internal/models/chat.go](/chat-service/internal/models/chat.go) | Go | 26 | 8 | 5 | 39 |
| [chat-service/internal/models/user.go](/chat-service/internal/models/user.go) | Go | 13 | 0 | 2 | 15 |
| [chat-service/internal/repository/chat\_repository.go](/chat-service/internal/repository/chat_repository.go) | Go | 255 | 15 | 37 | 307 |
| [chat-service/internal/repository/message\_repository.go](/chat-service/internal/repository/message_repository.go) | Go | 1 | 0 | 1 | 2 |
| [chat-service/internal/service/chat\_service.go](/chat-service/internal/service/chat_service.go) | Go | 61 | 10 | 20 | 91 |
| [chat-service/internal/usecase/chat\_usecase.go](/chat-service/internal/usecase/chat_usecase.go) | Go | 1 | 58 | 12 | 71 |
| [chat-service/test/k6.js](/chat-service/test/k6.js) | JavaScript | 35 | 2 | 13 | 50 |
| [chat-service/test/test.go](/chat-service/test/test.go) | Go | 1 | 23 | 7 | 31 |
| [chat-service/test/test2.go](/chat-service/test/test2.go) | Go | 59 | 4 | 11 | 74 |
| [chat-service/tmp/build-errors.log](/chat-service/tmp/build-errors.log) | Log | 1 | 0 | 0 | 1 |
| [chat-service/utils/jwt\_utils.go](/chat-service/utils/jwt_utils.go) | Go | 1 | 0 | 1 | 2 |
| [chat-service/utils/logger.go](/chat-service/utils/logger.go) | Go | 31 | 6 | 7 | 44 |

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details