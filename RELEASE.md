# Release Management Guide

Hướng dẫn quản lý version và release cho các Go libraries trong monorepo này.

## Mục lục

- [1. Tổng quan](#1-tổng-quan)
- [2. Cấu trúc Project](#2-cấu-trúc-project)
- [3. Quy trình Release](#3-quy-trình-release)
- [4. Conventional Commits](#4-conventional-commits)
- [5. Ví dụ Release](#5-ví-dụ-release)
- [6. CI/CD Integration](#6-cicd-integration)

---

## 1. Tổng quan

### Vấn đề

Trong monorepo chứa nhiều Go libraries (greetings, math, ...), việc quản lý version và release cho từng library riêng lẻ rất phức tạp nếu làm thủ công:

- Phải cập nhật version trong go.mod
- Tạo CHANGELOG thủ công
- Tạo Git tag tương ứng
- Đảm bảo code quality (test, lint)

### Giải pháp

Sử dụng **Release It** để tự động hóa 100% quy trình release:

✅ **Tự động tăng version** - Semantic versioning (patch/minor/major)  
✅ **Tự động sinh CHANGELOG** - Dựa trên conventional commits  
✅ **Tự động tạo Git tags** - Format: `<library>/vX.Y.Z`  
✅ **Tự động chạy tests & lint** - Đảm bảo code quality  
✅ **Tách biệt releases** - Mỗi library có chu kỳ release độc lập

---

## 2. Cấu trúc Project

```
monorepo-lib/
├── libs/
│   ├── greetings/
│   │   ├── .release-it.json       # Config release cho greetings
│   │   ├── CHANGELOG.md            # Auto-generated changelog
│   │   ├── go.mod
│   │   └── greetings.go
│   └── math/
│       ├── .release-it.json       # Config release cho math
│       ├── CHANGELOG.md            # Auto-generated changelog
│       ├── go.mod
│       └── math.go
├── cmd/
│   └── hello/
│       ├── go.mod
│       └── main.go
├── scripts/
│   └── release.sh                  # Release helper script
├── package.json                    # npm dependencies & scripts
├── go.work                         # Go workspace
└── RELEASE.md                      # Tài liệu này
```

---

## 3. Quy trình Release

### Bước 1: Cài đặt dependencies (chỉ làm 1 lần)

```bash
npm install
```

Hoặc với pnpm:

```bash
pnpm install
```

### Bước 2: Commit code theo Conventional Commits

Xem phần [4. Conventional Commits](#4-conventional-commits) để biết cách commit đúng chuẩn.

### Bước 3: Release library

#### Cách 1: Sử dụng script helper

```bash
# Release greetings library (patch version)
./scripts/release.sh greetings patch

# Release math library (minor version)
./scripts/release.sh math minor

# Release với major version
./scripts/release.sh greetings major
```

#### Cách 2: Sử dụng npm scripts

```bash
# Release greetings
npm run release:greetings

# Release math
npm run release:math

# Với version type cụ thể
npm run release:greetings -- minor
```

#### Cách 3: Chạy trực tiếp release-it

```bash
# Dry run (không thực sự release)
npx release-it --config ./libs/greetings/.release-it.json --dry-run

# Release thật
npx release-it --config ./libs/greetings/.release-it.json
```

### Release It sẽ tự động:

1. ✅ Chạy `go vet` và `go fmt`
2. ✅ Sinh CHANGELOG.md dựa trên commits
3. ✅ Tăng version trong go.mod
4. ✅ Commit changes
5. ✅ Tạo Git tag: `<library>/vX.Y.Z`
6. ✅ Push tag lên GitHub
7. ✅ Tạo GitHub Release (nếu configured)

---

## 4. Conventional Commits

Release It sử dụng **Angular commit convention** để tự động sinh changelog.

### Format

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types

| Type | Mô tả | Version Bump |
|------|-------|--------------|
| `feat` | Tính năng mới | MINOR |
| `fix` | Sửa bug | PATCH |
| `perf` | Cải thiện performance | PATCH |
| `refactor` | Refactor code | PATCH |
| `docs` | Cập nhật docs | - |
| `style` | Format code | - |
| `test` | Thêm tests | - |
| `chore` | Công việc khác | - |
| `BREAKING CHANGE` | Breaking change | MAJOR |

### Scope

Scope nên là tên library: `greetings`, `math`

### Ví dụ

```bash
# Feature mới cho greetings library
git commit -m "feat(greetings): add WelcomeMultiple function"

# Sửa bug cho math library
git commit -m "fix(math): handle division by zero correctly"

# Breaking change
git commit -m "feat(greetings): redesign Hello API

BREAKING CHANGE: Hello() now requires a Language parameter"

# Refactor
git commit -m "refactor(math): optimize Max/Min functions"

# Documentation
git commit -m "docs(greetings): add usage examples"
```

---

## 5. Ví dụ Release

### Scenario 1: Release greetings library v0.1.0

```bash
# 1. Commit các thay đổi theo conventional commits
git add libs/greetings/greetings.go
git commit -m "feat(greetings): add Hello, Goodbye, Welcome functions"

# 2. Chạy release
./scripts/release.sh greetings minor

# 3. Kết quả:
# - CHANGELOG.md được cập nhật
# - Tag greetings/v0.1.0 được tạo
# - Commit "chore(greetings): release v0.1.0" được push
```

### Scenario 2: Sửa bug cho math library

```bash
# 1. Fix bug
git add libs/math/math.go
git commit -m "fix(math): prevent division by zero panic"

# 2. Release patch version
./scripts/release.sh math patch

# 3. Kết quả:
# - Version tăng từ v0.1.0 → v0.1.1
# - Tag math/v0.1.1 được tạo
```

### Scenario 3: Breaking change

```bash
# 1. Commit với BREAKING CHANGE
git commit -m "feat(greetings)!: change Hello signature

BREAKING CHANGE: Hello() now returns (string, error)"

# 2. Release sẽ tự động bump major version
./scripts/release.sh greetings major

# 3. Kết quả:
# - Version tăng từ v0.1.0 → v1.0.0
# - CHANGELOG ghi rõ breaking change
```

---

## 6. CI/CD Integration

### GitHub Actions

Tạo file `.github/workflows/release.yml`:

```yaml
name: Release

on:
  workflow_dispatch:
    inputs:
      library:
        description: 'Library to release (greetings, math)'
        required: true
        type: choice
        options:
          - greetings
          - math
      version:
        description: 'Version bump type'
        required: true
        type: choice
        options:
          - patch
          - minor
          - major

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
          token: ${{ secrets.GITHUB_TOKEN }}

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install dependencies
        run: npm install

      - name: Configure Git
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"

      - name: Run tests
        run: cd libs/${{ inputs.library }} && go test ./...

      - name: Release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: npm run release:${{ inputs.library }} -- ${{ inputs.version }}
```

### GitLab CI

Tạo file `.gitlab-ci.yml`:

```yaml
release:
  stage: release
  image: golang:1.21
  before_script:
    - apt-get update && apt-get install -y nodejs npm
    - npm install
  script:
    - |
      if [ "$CI_COMMIT_MESSAGE" =~ "feat(greetings)" ]; then
        npm run release:greetings
      elif [ "$CI_COMMIT_MESSAGE" =~ "feat(math)" ]; then
        npm run release:math
      fi
  only:
    - main
```

---

## 7. Best Practices

### ✅ DO

- Luôn commit theo conventional commits format
- Chạy tests trước khi release: `go test ./...`
- Sử dụng semantic versioning đúng cách
- Review CHANGELOG trước khi push tag
- Sử dụng `--dry-run` để kiểm tra trước khi release thật

### ❌ DON'T

- Không sửa version trong go.mod thủ công
- Không tạo tag thủ công
- Không commit message theo kiểu "update code", "fix"
- Không release nhiều libraries cùng lúc nếu không cần thiết

---

## 8. Troubleshooting

### Lỗi: "Working directory is not clean"

```bash
# Stash hoặc commit changes trước
git stash
# hoặc
git add . && git commit -m "chore: prepare for release"
```

### Lỗi: "No commits found"

```bash
# Đảm bảo đã commit code với conventional format
git log --oneline | grep "feat\|fix"
```

### Xem dry-run trước khi release

```bash
npx release-it --config ./libs/greetings/.release-it.json --dry-run
```

### Rollback một release

```bash
# Xóa tag local
git tag -d greetings/v0.1.0

# Xóa tag remote
git push origin :refs/tags/greetings/v0.1.0

# Revert commit
git revert HEAD
```

---

## 9. Thêm Library Mới

Khi thêm library mới vào monorepo:

1. Tạo thư mục library
```bash
mkdir -p libs/newlib
cd libs/newlib
go mod init github.com/phongthien99/monorepo-lib/libs/newlib
```

2. Copy config từ library khác
```bash
cp ../greetings/.release-it.json .
```

3. Sửa tên library trong config
```json
{
  "git": {
    "tagName": "newlib/v${version}",
    "commitMessage": "chore(newlib): release v${version}"
  }
}
```

4. Thêm scripts vào package.json
```json
{
  "scripts": {
    "release:newlib": "release-it --config ./libs/newlib/.release-it.json"
  }
}
```

5. Tạo CHANGELOG.md
```bash
touch libs/newlib/CHANGELOG.md
```

---

## 10. Tài liệu tham khảo

- [Release It Documentation](https://github.com/release-it/release-it)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
