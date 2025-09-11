#!/bin/bash
set -e

# Clean previously generated files
rm -rf ISO20022
mkdir -p ISO20022

files=($(find ./Internal/XSD -name "*.xsd" | sort -u))
for file in "${files[@]}"
do
    moovio_xsd2go convert "$file" github.com/mbanq/iso20022-go ISO20022
done

# Use our custom Time types
go run ./scripts/fix_imports.go
go run ./scripts/fix_inner_xml.go

# run go fmt and goimports for every generated file
files=($(find ./ISO20022 -name '*.go'))
for file in "${files[@]}"
do
    gofmt -w $file
    goimports -w $file
done