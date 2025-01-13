# run program
file := "./main.go"   
run: cls
	@go run $(file)

# clean termianl
cls:
	@clear