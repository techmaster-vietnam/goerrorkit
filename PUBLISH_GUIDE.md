# Hướng dẫn Publish GoErrorKit lên GitHub

## Bước 1: Chuẩn bị GitHub Repository

### 1.1. Tạo Repository mới trên GitHub

1. Truy cập https://github.com/new
2. Điền thông tin:
   - **Repository name:** `goerrorkit`
   - **Description:** `🚀 Framework-agnostic error handling library for Go with accurate panic location tracking`
   - **Visibility:** Public
   - **Không tick:** "Add a README file", "Add .gitignore", "Choose a license" (đã có sẵn)
3. Click **Create repository**

### 1.2. Copy Repository URL

Sau khi tạo, GitHub sẽ hiển thị URL như:
```
https://github.com/your-username/goerrorkit.git
```

## Bước 2: Initialize Git và Push Code

### 2.1. Di chuyển thư mục goerrorkit ra ngoài

Hiện tại thư mục `goerrorkit` đang nằm trong `fiber_log/`. Ta cần di chuyển nó ra ngoài:

```bash
# Từ thư mục fiber_log
cd /Users/cuong/CODE/fiber_log

# Copy toàn bộ goerrorkit ra ngoài
cp -r goerrorkit /Users/cuong/CODE/goerrorkit

# Hoặc dùng mv để di chuyển
# mv goerrorkit /Users/cuong/CODE/goerrorkit
```

### 2.2. Initialize Git Repository

```bash
cd /Users/cuong/CODE/goerrorkit

# Initialize git
git init

# Add all files
git add .

# Commit
git commit -m "Initial commit: GoErrorKit v0.1.0

Features:
- Framework-agnostic error handling
- Accurate panic location tracking
- Custom error types (Business, System, Validation, Auth, External)
- Structured logging with JSON format
- File logging with rotation
- Fiber adapter support
"
```

### 2.3. Connect và Push lên GitHub

```bash
# Add remote (thay YOUR_USERNAME bằng username GitHub của bạn)
git remote add origin https://github.com/YOUR_USERNAME/goerrorkit.git

# Rename branch to main (nếu cần)
git branch -M main

# Push code
git push -u origin main
```

## Bước 3: Tạo Git Tag cho Version v0.1.0

```bash
# Tạo annotated tag
git tag -a v0.1.0 -m "Release v0.1.0

Initial release with:
- Core error handling functionality
- Panic recovery with accurate location tracking
- Fiber v2 adapter
- Logrus-based logging
- Example application
"

# Push tag lên GitHub
git push origin v0.1.0
```

## Bước 4: Cập nhật Module Path

### 4.1. Update go.mod files

Sau khi push lên GitHub, bạn cần cập nhật module paths trong các file:

**goerrorkit/go.mod:**
```go
module github.com/YOUR_USERNAME/goerrorkit
```

**goerrorkit/adapters/fiber/*.go:**
```go
import (
    "github.com/YOUR_USERNAME/goerrorkit/core"
    ...
)
```

**goerrorkit/config/logger.go:**
```go
import (
    "github.com/YOUR_USERNAME/goerrorkit/core"
    ...
)
```

**goerrorkit/examples/fiber-demo/go.mod:**
```go
require (
    github.com/YOUR_USERNAME/goerrorkit v0.1.0
    ...
)

// For local development
replace github.com/YOUR_USERNAME/goerrorkit => ../..
```

### 4.2. Commit và push changes

```bash
# Replace tất cả "github.com/cuong/goerrorkit" bằng path thực của bạn
find . -type f -name "*.go" -o -name "go.mod" | xargs sed -i '' 's|github.com/cuong/goerrorkit|github.com/YOUR_USERNAME/goerrorkit|g'

# Commit
git add .
git commit -m "Update module paths to actual GitHub repository"

# Push
git push origin main

# Update tag
git tag -d v0.1.0
git push origin :refs/tags/v0.1.0
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

## Bước 5: Verify trên pkg.go.dev

### 5.1. Trigger Indexing

Go proxy sẽ tự động index module khi có tag mới. Để verify:

```bash
# Request module từ Go proxy
go get github.com/YOUR_USERNAME/goerrorkit@v0.1.0
```

### 5.2. Check pkg.go.dev

Sau 5-10 phút, module sẽ xuất hiện tại:
```
https://pkg.go.dev/github.com/YOUR_USERNAME/goerrorkit
```

Nếu chưa thấy, có thể request manually:
```
https://pkg.go.dev/github.com/YOUR_USERNAME/goerrorkit@v0.1.0
```

## Bước 6: Tạo GitHub Release (Optional)

1. Truy cập `https://github.com/YOUR_USERNAME/goerrorkit/releases/new`
2. Chọn tag: `v0.1.0`
3. Release title: `v0.1.0 - Initial Release`
4. Description:
```markdown
## 🎉 Initial Release

GoErrorKit v0.1.0 brings accurate panic location tracking and comprehensive error handling to Go applications.

### ✨ Features
- ✅ Automatic panic recovery with exact source location
- ✅ Detailed stack trace with full call chain
- ✅ Framework-agnostic core design
- ✅ Fiber v2 adapter included
- ✅ Custom error types (Business, System, Validation, Auth, External)
- ✅ Structured JSON logging
- ✅ File logging with automatic rotation

### 📦 Installation
```bash
go get github.com/YOUR_USERNAME/goerrorkit@v0.1.0
```

### 📚 Documentation
See [README.md](README.md) for full documentation and examples.

### 🚀 Quick Start
```go
import (
    "github.com/YOUR_USERNAME/goerrorkit/adapters/fiber"
    "github.com/YOUR_USERNAME/goerrorkit/config"
    "github.com/YOUR_USERNAME/goerrorkit/core"
)

func main() {
    config.InitDefaultLogger()
    core.ConfigureForApplication("github.com/yourname/yourapp")
    
    app := fiber.New()
    app.Use(fiber.ErrorHandler())
    // ...
}
```
```

5. Click **Publish release**

## Bước 7: Update README với Badges (Optional)

Thêm badges vào đầu README.md:

```markdown
# GoErrorKit

[![Go Reference](https://pkg.go.dev/badge/github.com/YOUR_USERNAME/goerrorkit.svg)](https://pkg.go.dev/github.com/YOUR_USERNAME/goerrorkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/goerrorkit)](https://goreportcard.com/report/github.com/YOUR_USERNAME/goerrorkit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/YOUR_USERNAME/goerrorkit)](https://github.com/YOUR_USERNAME/goerrorkit/releases)

🚀 **Framework-agnostic error handling library for Go** ...
```

## Bước 8: Test Installation

Tạo project mới để test:

```bash
mkdir /tmp/test-goerrorkit
cd /tmp/test-goerrorkit
go mod init test

# Install thư viện
go get github.com/YOUR_USERNAME/goerrorkit@v0.1.0

# Tạo test file
cat > main.go << 'EOF'
package main

import (
    "github.com/YOUR_USERNAME/goerrorkit/core"
    "fmt"
)

func main() {
    err := core.NewBusinessError(404, "Test error")
    fmt.Println(err)
}
EOF

# Run
go run main.go
```

## Bước 9: Promote Library

### 9.1. Share trên các channels

- Reddit: r/golang
- Twitter/X: #golang hashtag
- Go Forum: https://forum.golangbridge.org/
- Dev.to: Write an article about your library
- HackerNews: Share your GitHub repo

### 9.2. Tạo các project examples

Tạo các project examples sử dụng thư viện và share:
- Blog posts
- YouTube tutorials
- GitHub discussions

## Checklist Summary

- [ ] Tạo GitHub repository
- [ ] Copy code và initialize git
- [ ] Push code lên GitHub
- [ ] Tạo git tag v0.1.0
- [ ] Update module paths
- [ ] Verify trên pkg.go.dev
- [ ] Tạo GitHub Release
- [ ] Add badges to README
- [ ] Test installation
- [ ] Share với community

## Troubleshooting

### Module không xuất hiện trên pkg.go.dev

1. Check tag exists: `git tag -l`
2. Verify module path trong go.mod
3. Ensure go.mod file ở root directory
4. Wait 10-15 minutes for indexing
5. Manually request: Visit `https://pkg.go.dev/github.com/YOUR_USERNAME/goerrorkit@v0.1.0`

### Import errors

1. Run `go mod tidy` in your project
2. Check module path spelling
3. Verify version tag exists on GitHub

---

🎉 Chúc mừng! Library của bạn đã sẵn sàng để chia sẻ với cộng đồng Go!

