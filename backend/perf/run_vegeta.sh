#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://localhost:30099"
GET_PATH="/v1/get?key=54a34060-4e99-4a41-afd4-a3b32d88f650"
#BULK_PATH="/v1/bulk?key=749510f7-c84b-4994-88ac-530a7606382c&data_keys=beffdd7d-8bfa-49c6-b802-a9ad7ecd1456,afa7cb5a-ac48-4fd4-b1f7-13fd12824d5b,64fdf6e6-e2d0-49b5-8e50-11970cdb9207"

RATE=5000
DURATION=15m
CONNECTIONS=100
CUR_TIME="`date +%Y.%m.%d.%H.%M.%S`";
OUT_DIR="./result/${CUR_TIME}"
mkdir -p "${OUT_DIR}"

BIN_FILE="${OUT_DIR}/results_${RATE}rps_${DURATION}.bin"
TXT_FILE="${OUT_DIR}/results_${RATE}rps_${DURATION}.txt"

cat <<EOF | vegeta attack \
  -rate=${RATE} \
  -duration=${DURATION} \
  -connections=${CONNECTIONS} \
  | tee "${BIN_FILE}" \
  | vegeta report > "${TXT_FILE}"
GET ${BASE_URL}${GET_PATH}
EOF

echo "Done."
echo "Binary results: ${BIN_FILE}"
echo "TXT report:   ${TXT_FILE}"
