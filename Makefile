# run program
file := "./example/main.go"   
run: cls
	@go run $(file)

# clean termianl
cls:
	@clear