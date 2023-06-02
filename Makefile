icon = ./assets/icon.png
appID = com.mixedmachine.simplebudget

init:
	@echo "Initializing..."
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Verifying dependencies..."
	@go mod verify
	@echo prettifying code...
	@go fmt ./...

build.win: init
	@echo "Building windows executable..."
	@go build -o bin/budget.exe

build.lin: init
	@echo "Building linux executable..."
	@go build -o bin/budget ./...

build.all: build.win build.lin
	@echo "Building for all platforms complete."

run.win: build.win
	@echo "Running..."
	@./bin/budget.exe

run.lin: build.lin
	@echo "Running..."
	@./bin/budget

dev: init
	@echo "Running in dev mode..."
	@go run ./main.go

debug: init
	@echo "Debugging..."
	@dlv debug ./main.go

test:
	@echo "Testing..."
	@go test -v ./...

scan:
	@echo "Scanning..."
	@golangci-lint run ./...

secure:
	@echo "Security..."
	@gosec ./...

pkg.init:
	@go install fyne.io/fyne/v2/cmd/fyne@latest

pkg.mobile.and: pkg.init
	@echo "Packaging for android..."
	@mkdir -p bin/mobile
	fyne package -os android -appID $(appID) -icon $(icon)
	@mv simple_budget_app.apk bin/mobile/simple_budget_app.apk

pkg.mobile.ios: pkg.init
	@echo "Packaging for ios..."
	@mkdir -p bin/mobile
	fyne package -os ios -appID $(appID) -icon $(icon)
	@mv simple_budget_app.ipa bin/mobile/simple_budget_app.ipa

pkg.mobile.all: pkg.mobile.and pkg.mobile.ios
	@echo "Packaging for all mobile platforms complete."

pkg.desktop.win: pkg.init
	@echo "Packaging for windows..."
	@mkdir -p bin/desktop
	fyne package --release --appID $(appID) --exe bin/desktop --os windows --icon $(icon)

pkg.desktop.lin: pkg.init
	@echo "Packaging for linux..."
	@mkdir -p bin/desktop
	fyne package --exe bin/desktop --os linux --icon $(icon)

pkg.desktop.mac: pkg.init
	@echo "Packaging for mac..."
	@mkdir -p bin/desktop
	fyne package --exe bin/desktop --os darwin --icon $(icon)

pkg.desktop.all: pkg.desktop.linux pkg.desktop.mac pkg.desktop.win
	@echo "Packaging for all desktop platforms complete."

clean:
	@echo "Cleaning..."
	@rm -rf bin/

.PHONY: init \
build.win build.lin build.all \
run.win run.lin dev \
test scan secure \
pkg.init \
pkg.mobile.and pkg.mobile.ios \
pkg.mobile.all \
pkg.desktop.win pkg.desktop.lin pkg.desktop.mac \
pkg.desktop.all \
clean