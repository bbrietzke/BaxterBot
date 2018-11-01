/*
Package web is a pretty basic web server at this point.

Manual Testing with curl:

	curl --header "Content-Type: application/json" --request POST --data '{"username":"xyz","password":"xyz"}' http://localhost:8080/api/fred
	curl -X GET -H "Content-type: application/json" -H "Accept: application/json"  "http://localhost:8080/api/fred"

Load Testing with vegeta
	cat body.txt | vegeta attack -duration 10s  | tee /tmp/report.bin | vegeta report -type=text && cat /tmp/report.bin | vegeta plot > /tmp/page.html && open /tmp/page.html
*/
package web
