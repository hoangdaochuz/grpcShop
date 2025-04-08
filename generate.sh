#!/bin/bash

#!/bin/bash

# Đường dẫn tới các thư mục
PROTO_DIR="./proto"
BE_BASE_DIR="./backend/apis"
FE_GEN_DIR="./frontend/src/apis"

# Tạo thư mục frontend nếu chưa tồn tại
mkdir -p $FE_GEN_DIR

# Tìm tất cả file .proto trong thư mục proto/
PROTO_FILES=$(find $PROTO_DIR -name "*.proto")

# Kiểm tra xem có file .proto nào không
if [ -z "$PROTO_FILES" ]; then
  echo "No .proto files found in $PROTO_DIR!"
  exit 1
fi

# Sinh code cho Backend (Go + gRPC + gRPC-Gateway)
echo "Generating Backend code..."
for PROTO_FILE in $PROTO_FILES; do
  # Lấy tên file không có đuôi .proto (ví dụ: "ecommerce" từ "ecommerce.proto")
  BASE_NAME=$(basename "$PROTO_FILE" .proto)
  
  # Thư mục đích cho service (ví dụ: backend/product/)
  BE_GEN_DIR="$BE_BASE_DIR/$BASE_NAME"
  
  # Tạo thư mục nếu chưa tồn tại
  mkdir -p "$BE_GEN_DIR"
  
  echo "Processing $PROTO_FILE -> $BE_GEN_DIR..."
  protoc -I $PROTO_DIR \
    --go_out="$BE_GEN_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$BE_GEN_DIR" --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out="$BE_GEN_DIR" --grpc-gateway_opt=paths=source_relative \
    "$PROTO_FILE"
done

# Sinh code cho Frontend (TypeScript)
echo "Generating Frontend code..."
for PROTO_FILE in $PROTO_FILES; do
  echo "Processing $PROTO_FILE..."
  protoc -I $PROTO_DIR \
    --ts_out=$FE_GEN_DIR \
    $PROTO_FILE
done

echo "Code generation completed!"
echo "Backend code: $BE_GEN_DIR"
echo "Frontend code: $FE_GEN_DIR"
