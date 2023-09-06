#!/bin/bash

inputFolder=$1

# Get data in CSV files
mkdir -p $inputFolder
parallel echo {} '|' quotes get --period 7 '>' $inputFolder/{}.csv

# Merge CSV files
i=0

for filename in $inputFolder/*.csv; do
  ((i++))

  if [ $i -eq 1 ]
  then
    cp $filename $inputFolder/tmp$i.csv
  fi

  xsv join date $inputFolder/tmp$i.csv date $inputFolder/"$(basename "$filename")" | xsv select '!date[1]' > $inputFolder/tmp$((i+1)).csv
done

# Output result
cat $inputFolder/tmp$i.csv

# Remove tmp files
rm $inputFolder/tmp*.csv
