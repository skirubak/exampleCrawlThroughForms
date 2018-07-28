#!/bin/bash
rm -f ./tempFile.txt
rm -f ./outFile.txt
# first awk is printing below from Scott list.
#Last Name,First Name,Middle Name,Address Line 1,City,State,Zip,Gender,Ethnicity,County,APN Category,APN Approval Date
awk -F , '{
	for(i=1;i<NF;i++){
        	gsub(/[ \t]+$/,"",$i)
		}
	printf $2","$3","$4","$5","$8","$9","$10","$12","$13","$14","$15","$20","
	if (split($5,a," ",seps) > 3) {
		print a[2]"|"a[3]
		} else {
		print a[2]
		}
        }' uniqList.csv > ./tempFile.txt

sed 1d ./tempFile.txt | while read line
do
input=`echo $line |awk -F , '{for(i=1;i<NF;i++){printf $i","}}'`
street=`echo $line |awk -F , '{print $NF}'`
modifyStreet=`echo $street | awk -F "|" '{if (NF > 1 ) {
						if ( $1 == "NORTH" ) {
							print "N "$2"|"$1" "$2
							} else {
							print $1" "$2
							}
						} else {
							print $0
						}
					}'`
lastName=`echo $line |awk -F , '{print $1}'`
county=`echo $line |awk -F , '{print $10}'`
if [ "$county" = "COLLIN" ]; then
	echo -n "$input"
	grep -E "$modifyStreet" collincad.txt |grep "$lastName" | awk 'BEGIN{FPAT="([^,]*)|(\"[^\"]+\")";}{printf $3$5","$6","$73","}'
else
	echo -n "$input"
	grep -E "$modifyStreet" dentonShortOutput.txt | grep "$lastName" |awk '{printf substr($0,609,70)","substr($0,1916,15)",";}'
fi
echo ""
done
