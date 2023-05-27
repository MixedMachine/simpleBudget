icon = icon.png
appID = com.mixedmachine.simple-budget

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
	@go build -o bin/budget.exe ./...

build.lin: init
	@echo "Building linux executable..."
	@go build -o bin/budget ./...

run: build.win build.lin
	@echo "Building for all platforms complete."
	@echo "Running..."
	@./bin

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

package.mobile.and:
	@echo "Packaging for android..."
	@mkdir -p bin/mobile
	fyne package -os android -appID $(appID) -icon $(icon)

package.mobile.ios:
	@echo "Packaging for ios..."
	@mkdir -p bin/mobile
	fyne package -os ios -appID $(appID) -icon $(icon)

package.mobile.all: package.mobile.and package.mobile.ios
	@echo "Packaging for all mobile platforms complete."

package.desktop.win:
	@echo "Packaging for windows..."
	@mkdir -p bin/desktop
	fyne package -os windows -icon $(icon)

package.desktop.lin:
	@echo "Packaging for linux..."
	@mkdir -p bin/desktop
	fyne package -os linux -icon $(icon)

package.desktop.mac:
	@echo "Packaging for mac..."
	@mkdir -p bin/desktop
	fyne package -os darwin -icon $(icon)

package.desktop.all: package.desktop.linux package.desktop.mac package.desktop.win
	@echo "Packaging for all desktop platforms complete."

clean:
	@echo "Cleaning..."
	@rm -rf bin/

.PHONY: init \
build.win build.lin run dev \
test scan secure \
package.mobile.and package.mobile.ios package.mobile.all \
package.desktop.win package.desktop.lin package.desktop.mac package.desktop.all \
clean