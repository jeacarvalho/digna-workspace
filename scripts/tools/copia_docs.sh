#!/usr/bin/env bash

SRC_DIR="$1"
OUTPUT_FILE="$2"

if [ -z "$SRC_DIR" ] || [ -z "$OUTPUT_FILE" ]; then
  echo "Uso: $0 <diretorio_origem> <arquivo_saida>"
  exit 1
fi

SRC_DIR=$(realpath "$SRC_DIR")

echo "Diretório origem: $SRC_DIR"
echo "Arquivo saída: $OUTPUT_FILE"
echo ""

TMP_LIST=$(mktemp)

find "$SRC_DIR" -type f -name "*.md" | sort > "$TMP_LIST"

# limpa saída
> "$OUTPUT_FILE"

echo "# MERGED DOCUMENTATION" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "## FILE INDEX" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# índice
while read -r file; do
    rel="${file#$SRC_DIR/}"
    echo "- $rel" >> "$OUTPUT_FILE"
done < "$TMP_LIST"

echo "" >> "$OUTPUT_FILE"
echo "---" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# concatenação
while read -r file; do

    rel="${file#$SRC_DIR/}"
    size=$(stat -c%s "$file")

    echo "Processando: $rel (${size} bytes)"

    echo "===== BEGIN FILE: $rel =====" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"

    tr -d '\000' < "$file" >> "$OUTPUT_FILE"

    echo "" >> "$OUTPUT_FILE"
    echo "===== END FILE: $rel =====" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"

done < "$TMP_LIST"

rm "$TMP_LIST"

echo ""
echo "Merge concluído."
