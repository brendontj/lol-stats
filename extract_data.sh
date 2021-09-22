#!/bin/bash
# Export PATH in case SENDMAIL is not found by the script
export PATH=/usr/sbin:/usr/bin
# EXPORT POSTGRES SQL COMMAND AS CSV
# Adjust the Select SQL command to your liking
# Note: You may need to define the path to export.csv file
psql -d lolstats -h localhost -t -A -F"," -f /pkg/db/sql/queries/extract_data.sql > export.csv
# SEND AN EMAIL
# Note: You may need to define the path to export.csv file on the bottom line
( echo "to: brendonhps@gmail.com"
echo "from: POSTGRESQL SERVER <server@domain.com>"
echo "Subject: PostgreSQL Export as of $(date +%F)"
echo "mime-version: 1.0"
echo "content-type: multipart/related; boundary=messageBoundary"
echo
echo "--messageBoundary"
echo "content-type: text/plain"
echo
echo "Please find the export in CSV format attached to this email."
echo "Created: $(date)."
echo
echo "--messageBoundary"
echo "content-type: text; name=export.csv"
echo "content-transfer-encoding: base64"
echo
openssl base64 < export.csv) | sendmail -t -i