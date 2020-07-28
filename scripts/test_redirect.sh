#!/bin/bash
# simply tests if we get redirected by the server
if [ -z "$2" ]; then
	address=http://localhost:8080
else
	address=http://$2
fi
curl -s -L -D - $address/$1 -o /dev/null -w '%{url_effective}' | grep -E "404 | 500" || echo "Redirect Success"

