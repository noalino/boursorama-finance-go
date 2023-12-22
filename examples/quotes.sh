#!/bin/bash

inputFolder=$1

# Get data in CSV files
mkdir -p $inputFolder
parallel echo {} '|' quotes get --period weekly '>' $inputFolder/{}.csv

# Merge CSV files
i=0
header="date"

for file in $inputFolder/*.csv; do
  ((i++))

  filename=$(basename $file)
  header+=",${filename%.*}"

  if [ $i -eq 1 ]
  then
    xsv select date,close $file > $inputFolder/tmp$i.csv
    continue
  fi

  xsv join date $inputFolder/tmp$((i-1)).csv date $inputFolder/"$(basename "$file")" | xsv select '!date[1],open,performance,high,low' > $inputFolder/tmp$i.csv
done

# Replace header
sed -i '' "1s/.*/$header/" $inputFolder/tmp$i.csv

# Output result
cat $inputFolder/tmp$i.csv

# Remove tmp files
rm $inputFolder/tmp*.csv
