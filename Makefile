swaggo:
	echo "Starting swagger generating"
	swag init -g **/**/*.go

load-test-create-tags:
	jq -ncM 'while(true; '.'+1) | tostring | {method: "POST", url: "https://iiqbal2000-bareknews-45597j96cjvxp-3333.githubpreview.dev/api/tags", body: {name: .} | @base64 }' |   vegeta attack -lazy -format=json -duration=1s |  tee results.bin |  vegeta report