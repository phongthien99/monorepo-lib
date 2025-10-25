# Go Workspace Monorepo Example

Đây là một ví dụ về Go workspace với nhiều libraries và automated release management.

[![Release](https://img.shields.io/badge/release-automated-brightgreen.svg)](RELEASE.md)
[![Go Version](https://img.shields.io/badge/go-1.21-blue.svg)](https://golang.org/)

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

## Release Management

Project này sử dụng **Release It** để tự động hóa quy trình release cho từng library.

### Quick Start

1. Cài đặt dependencies:
```bash
npm install
```

2. Commit code theo [Conventional Commits](https://www.conventionalcommits.org/):
```bash
git commit -m "feat(greetings): add new WelcomeMultiple function"
```

3. Release library:
```bash
# Sử dụng script helper
./scripts/release.sh greetings patch

# Hoặc sử dụng npm scripts
npm run release:greetings
```

### Tính năng

✅ **Automated versioning** - Tự động tăng version theo semantic versioning  
✅ **Auto-generated changelog** - CHANGELOG.md được sinh tự động  
✅ **Git tags** - Tạo tags theo format `<library>/vX.Y.Z`  
✅ **Quality checks** - Tự động chạy tests và lint  
✅ **Independent releases** - Mỗi library có chu kỳ release riêng

### Hướng dẫn chi tiết

Xem [RELEASE.md](RELEASE.md) để biết thêm chi tiết về:
- Quy trình release đầy đủ
- Conventional commits guide
- CI/CD integration
- Troubleshooting

## Thêm library mới

Để thêm một library mới vào workspace:

1. Tạo thư mục cho library
```bash
mkdir -p libs/newlib
```

2. Khởi tạo module
```bash
cd libs/newlib
go mod init github.com/phongthien99/monorepo-lib/libs/newlib
```

3. Thêm vào go.work
```bash
cd ../..
go work use ./libs/newlib
```

4. Thiết lập release config (xem [RELEASE.md](RELEASE.md#9-thêm-library-mới))

## Contributing

Xem [RELEASE.md](RELEASE.md#4-conventional-commits) để biết cách commit code đúng chuẩn.

## License

MIT
