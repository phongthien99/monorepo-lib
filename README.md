# Go Workspace Monorepo Example

Đây là một ví dụ về Go workspace với nhiều libraries và một example application.

## Cấu trúc project

```
monorepo-lib/
├── go.work                    # Go workspace file
├── libs/
│   ├── greetings/            # Library 1: Greeting functions
│   │   ├── go.mod
│   │   └── greetings.go
│   └── math/                 # Library 2: Math utilities
│       ├── go.mod
│       └── math.go
└── cmd/
    └── hello/                # Example hello world app
        ├── go.mod
        └── main.go
```

## Libraries

### 1. Greetings Library (`libs/greetings`)
Cung cấp các hàm xử lý lời chào:
- `Hello(name string)` - Trả về lời chào
- `Goodbye(name string)` - Trả về lời tạm biệt
- `Welcome(names ...string)` - Chào mừng nhiều người

### 2. Math Library (`libs/math`)
Cung cấp các hàm toán học cơ bản:
- `Add(a, b int)` - Cộng
- `Subtract(a, b int)` - Trừ
- `Multiply(a, b int)` - Nhân
- `Divide(a, b int)` - Chia
- `Max(a, b int)` - Giá trị lớn nhất
- `Min(a, b int)` - Giá trị nhỏ nhất

## Cách chạy

### 1. Chạy example app

```bash
cd cmd/hello
go run main.go
```

### 2. Build example app

```bash
cd cmd/hello
go build -o hello
./hello
```

## Go Workspace

Project này sử dụng Go Workspace (Go 1.18+) để quản lý nhiều modules trong cùng một repository.

File `go.work` định nghĩa các modules trong workspace:

```go
go 1.21

use (
	./libs/greetings
	./libs/math
	./cmd/hello
)
```

## Lợi ích của Go Workspace

1. Quản lý nhiều modules trong cùng một repository
2. Dễ dàng phát triển và test các libraries cục bộ
3. Không cần sử dụng `replace` directive trong go.mod
4. Các thay đổi trong libraries được phản ánh ngay lập tức

## Thêm library mới

Để thêm một library mới vào workspace:

1. Tạo thư mục cho library
```bash
mkdir -p libs/newlib
```

2. Khởi tạo module
```bash
cd libs/newlib
go mod init github.com/yourusername/monorepo-lib/libs/newlib
```

3. Thêm vào go.work
```bash
cd ../..
go work use ./libs/newlib
```
