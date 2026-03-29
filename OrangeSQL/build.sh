#!/bin/bash
set -e

echo "=== Building frontend ==="
cd frontend && pnpm build && cd ..

echo "=== Building Go binary ==="
go build -o orangesql.exe .

echo "=== Done: orangesql.exe ==="
ls -lh orangesql.exe
