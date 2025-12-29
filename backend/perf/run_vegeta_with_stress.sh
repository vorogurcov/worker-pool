#!/usr/bin/env bash
set -euo pipefail

RATE=5000
DURATION=15m
CONNECTIONS=100

CUR_TIME="$(date +%Y.%m.%d.%H.%M.%S)"
OUT_DIR="./result/${CUR_TIME}"
mkdir -p "${OUT_DIR}"

BIN_FILE="${OUT_DIR}/results_${RATE}rps_${DURATION}.bin"
TXT_FILE="${OUT_DIR}/results_${RATE}rps_${DURATION}.txt"

# Читаем готовые vegeta targets из stdout Go-генератора
go run stress.go | vegeta attack \
  -rate="${RATE}" \
  -duration="${DURATION}" \
  -connections="${CONNECTIONS}" \
  -output="${BIN_FILE}"

vegeta report -inputs="${BIN_FILE}" > "${TXT_FILE}"

echo "Done."
echo "Binary results: ${BIN_FILE}"
echo "TXT report:   ${TXT_FILE}"
