#!/bin/bash
colt="init"
csv=""
while IFS=, read -r col1 col2 #IFS= Initial Field Seperator
do
    [ "$col1" != "$colt" ] && colt="$col1" && fname="$colt.csv" && touch "$fname" && echo "$csv" > "$fname" && csv=""
    echo "Current ID:$col1|$col2" && csv+="$col1,$col2 \n"
done < $1 || echo "Missing Input CSV-File"