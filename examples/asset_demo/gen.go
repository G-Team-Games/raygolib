package main

//go:generate go run github.com/G-Team-Games/raygolib/cmd/assetgen -root ./assets -pkg main -out ./generated_assets.go -kinds model -naming pascal -v -dry-run

// Use `go generate ./...` to generate asset paths in separate file.