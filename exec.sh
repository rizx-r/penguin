# ===apps/social===
goctl rpc protoc ./apps/social/rpc/social.proto --go_out=./apps/social/rpc/ --go-grpc_out=./apps/social/rpc/ --zrpc_out=./apps/social/rpc/

goctl model mysql ddl -src="./deploy/sql/social.sql" -dir="./apps/social/social_models" -c

goctl api go -api apps/social/api/social.api -dir apps/social/api -style gozero

# ===apps/user===
goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc/ --go-grpc_out=./apps/user/rpc/ --zrpc_out=./apps/user/rpc/

# 根据sql文件生成models
goctl model mysql ddl -src="./deploy/sql/user.sql" -dir="./apps/user/models" -c

goctl api go -api apps/user/api/user.api -dir apps/user/api -style gozero

# ===apps/im===
goctl model mongo --type chatLog --dir ./apps/im/im_models