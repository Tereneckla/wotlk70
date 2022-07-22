# Recursive wildcard function
rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

OUT_DIR=dist/wotlk
ASSETS_INPUT := $(shell find assets/ -type f)
ASSETS := $(patsubst assets/%,$(OUT_DIR)/assets/%,$(ASSETS_INPUT))
GOROOT := $(shell go env GOROOT)

ifeq ($(shell go env GOOS),darwin)
    SED:=sed -i "" -E -e
else
    SED:=sed -i -E -e
endif

# Make everything. Keep this first so it's the default rule.
$(OUT_DIR): ui_shared \
 $(OUT_DIR)/balance_druid/index.js $(OUT_DIR)/balance_druid/index.css $(OUT_DIR)/balance_druid/index.html \
 $(OUT_DIR)/feral_druid/index.js $(OUT_DIR)/feral_druid/index.css $(OUT_DIR)/feral_druid/index.html \
 $(OUT_DIR)/feral_tank_druid/index.js $(OUT_DIR)/feral_tank_druid/index.css $(OUT_DIR)/feral_tank_druid/index.html \
 $(OUT_DIR)/elemental_shaman/index.js $(OUT_DIR)/elemental_shaman/index.css $(OUT_DIR)/elemental_shaman/index.html \
 $(OUT_DIR)/enhancement_shaman/index.js $(OUT_DIR)/enhancement_shaman/index.css $(OUT_DIR)/enhancement_shaman/index.html \
 $(OUT_DIR)/hunter/index.js $(OUT_DIR)/hunter/index.css $(OUT_DIR)/hunter/index.html \
 $(OUT_DIR)/mage/index.js $(OUT_DIR)/mage/index.css $(OUT_DIR)/mage/index.html \
 $(OUT_DIR)/rogue/index.js $(OUT_DIR)/rogue/index.css $(OUT_DIR)/rogue/index.html \
 $(OUT_DIR)/retribution_paladin/index.js $(OUT_DIR)/retribution_paladin/index.css $(OUT_DIR)/retribution_paladin/index.html \
 $(OUT_DIR)/protection_paladin/index.js $(OUT_DIR)/protection_paladin/index.css $(OUT_DIR)/protection_paladin/index.html \
 $(OUT_DIR)/shadow_priest/index.js $(OUT_DIR)/shadow_priest/index.css $(OUT_DIR)/shadow_priest/index.html \
 $(OUT_DIR)/smite_priest/index.js $(OUT_DIR)/smite_priest/index.css $(OUT_DIR)/smite_priest/index.html \
 $(OUT_DIR)/warlock/index.js $(OUT_DIR)/warlock/index.css $(OUT_DIR)/warlock/index.html \
 $(OUT_DIR)/warrior/index.js $(OUT_DIR)/warrior/index.css $(OUT_DIR)/warrior/index.html \
 $(OUT_DIR)/protection_warrior/index.js $(OUT_DIR)/protection_warrior/index.css $(OUT_DIR)/protection_warrior/index.html \
 $(OUT_DIR)/deathknight/index.js $(OUT_DIR)/deathknight/index.css $(OUT_DIR)/deathknight/index.html \
 $(OUT_DIR)/raid/index.js $(OUT_DIR)/raid/index.css $(OUT_DIR)/raid/index.html

ui_shared: $(OUT_DIR)/lib.wasm \
 $(OUT_DIR)/core/tsconfig.tsbuildinfo \
 $(OUT_DIR)/index.html \
 $(OUT_DIR)/sim_worker.js \
 $(OUT_DIR)/net_worker.js \
 $(OUT_DIR)/detailed_results/index.js \
 $(OUT_DIR)/detailed_results/index.css \
 $(OUT_DIR)/detailed_results/index.html

$(OUT_DIR)/index.html:
	cp ui/index.html $(OUT_DIR)

.PHONY: clean
clean:
	rm -rf ui/core/proto/*.ts \
	  sim/core/proto/*.pb.go \
	  wowsimwotlk \
	  wowsimwotlk-windows.exe \
	  wowsimwotlk-amd64-darwin \
	  wowsimwotlk-amd64-linux \
	  dist \
	  binary_dist
	find . -name "*.results.tmp" -type f -delete

# Host a local server, for dev testing
.PHONY: host
host: $(OUT_DIR)
	# Intentionally serve one level up, so the local site has 'wotlk' as the first
	# directory just like github pages.
	npx http-server $(OUT_DIR)/..

ui/core/proto/api.ts: proto/*.proto node_modules
	mkdir -p $(OUT_DIR)/protobuf-ts
	cp -r node_modules/@protobuf-ts/runtime/build/es2015/* $(OUT_DIR)/protobuf-ts
	$(SED) "s/from '(.*)';/from '\1.js';/g" $(OUT_DIR)/protobuf-ts/*.js
	$(SED) "s/from \"(.*)\";/from '\1.js';/g" $(OUT_DIR)/protobuf-ts/*.js
	npx protoc --ts_opt generate_dependencies --ts_out ui/core/proto --proto_path proto proto/api.proto
	npx protoc --ts_out ui/core/proto --proto_path proto proto/test.proto
	npx protoc --ts_out ui/core/proto --proto_path proto proto/ui.proto

package-lock.json:
	npm install

node_modules: package-lock.json
	npm ci

$(OUT_DIR)/core/tsconfig.tsbuildinfo: $(call rwildcard,ui/core,*.ts) ui/core/proto/api.ts
	npx tsc -p ui/core
	$(SED) 's#@protobuf-ts/runtime#/wotlk/protobuf-ts/index#g' $(OUT_DIR)/core/proto/*.js
	$(SED) "s/from \"(.*)\";/from '\1.js';/g" $(OUT_DIR)/core/proto/*.js

# Generic rule for hosting any class directory
.PHONY: host_%
host_%: ui_shared %
	npx http-server $(OUT_DIR)/..

# Generic rule for building index.js for any class directory
$(OUT_DIR)/%/index.js: ui/%/index.ts ui/%/*.ts $(OUT_DIR)/core/tsconfig.tsbuildinfo
	npx tsc -p $(<D) 
	touch $@ # TSC does not guarantee a file touch.

# Generic rule for building index.css for any class directory
$(OUT_DIR)/%/index.css: ui/%/index.scss ui/%/*.scss $(call rwildcard,ui/core,*.scss)
	mkdir -p $(@D)
	npx sass $< $@

# Generic rule for building index.html for any class directory
$(OUT_DIR)/%/index.html: ui/index_template.html $(ASSETS)
	$(eval title := $(shell echo $(shell basename $(@D)) | sed -r 's/(^|_)([a-z])/\U \2/g' | cut -c 2-))
	echo $(title)
	mkdir -p $(@D)
	cat ui/index_template.html | sed 's/@@TITLE@@/WOTLK $(title) Simulator/g' > $@

.PHONY: wasm
wasm: $(OUT_DIR)/lib.wasm

# Builds the generic .wasm, with all items included.
$(OUT_DIR)/lib.wasm: sim/wasm/* sim/core/proto/api.pb.go $(filter-out sim/core/items/all_items.go, $(call rwildcard,sim,*.go))
	@echo "Starting webassembly compile now..."
	@if GOOS=js GOARCH=wasm go build -o ./$(OUT_DIR)/lib.wasm ./sim/wasm/; then \
		echo "\033[1;32mWASM compile successful.\033[0m"; \
	else \
		echo "\033[1;31mWASM COMPILE FAILED\033[0m"; \
		exit 1; \
	fi
	

# Generic sim_worker that uses the generic lib.wasm
$(OUT_DIR)/sim_worker.js: ui/worker/sim_worker.js
	cat $(GOROOT)/misc/wasm/wasm_exec.js > $(OUT_DIR)/sim_worker.js
	cat ui/worker/sim_worker.js >> $(OUT_DIR)/sim_worker.js

$(OUT_DIR)/net_worker.js: ui/worker/net_worker.js
	cp ui/worker/net_worker.js $(OUT_DIR)

$(OUT_DIR)/assets/%: assets/%
	mkdir -p $(@D)
	cp $< $@

binary_dist/dist.go: sim/web/dist.go.tmpl
	mkdir -p binary_dist/wotlk
	touch binary_dist/wotlk/embedded
	cp sim/web/dist.go.tmpl binary_dist/dist.go

binary_dist: $(OUT_DIR)
	rm -rf binary_dist
	mkdir -p binary_dist
	cp -r $(OUT_DIR) binary_dist/
	rm binary_dist/wotlk/lib.wasm
	rm -rf binary_dist/wotlk/assets/item_data

# Builds the web server with the compiled client.
.PHONY: wowsimwotlk
wowsimwotlk: binary_dist devserver

.PHONY: devserver
devserver: sim/core/proto/api.pb.go sim/web/main.go binary_dist/dist.go
	@echo "Starting server compile now..."
	@if go build -o wowsimwotlk ./sim/web/main.go; then \
		echo "\033[1;32mBuild Completed Succeessfully\033[0m"; \
	else \
		echo "\033[1;31mBUILD FAILED\033[0m"; \
		exit 1; \
	fi

rundevserver: devserver
	./wowsimwotlk --usefs=true --launch=false

release: wowsimwotlk
	GOOS=windows GOARCH=amd64 go build -o wowsimwotlk-windows.exe -ldflags="-X 'main.Version=$(VERSION)'" ./sim/web/main.go
	GOOS=darwin GOARCH=amd64 go build -o wowsimwotlk-amd64-darwin -ldflags="-X 'main.Version=$(VERSION)'" ./sim/web/main.go
	GOOS=linux GOARCH=amd64 go build -o wowsimwotlk-amd64-linux   -ldflags="-X 'main.Version=$(VERSION)'" ./sim/web/main.go

sim/core/proto/api.pb.go: proto/*.proto
	protoc -I=./proto --go_out=./sim/core ./proto/*.proto

.PHONY: items
items: sim/core/items/all_items.go sim/core/proto/api.pb.go

sim/core/items/all_items.go: generate_items/*.go $(call rwildcard,sim/core/proto,*.go)
	go run generate_items/*.go -outDir=sim/core/items
	gofmt -w ./sim/core/items

.PHONY: test
test: $(OUT_DIR)/lib.wasm binary_dist/dist.go
	go test ./...

.PHONY: update-test
update-tests:
	find . -name "*.results" -type f -delete
	find . -name "*.results.tmp" -exec bash -c 'cp "$$1" "$${1%.results.tmp}".results' _ {} \;

.PHONY: fmt
fmt: tsfmt
	gofmt -w ./sim
	gofmt -w ./generate_items

.PHONY: tsfmt
tsfmt:
	for dir in $$(find ./ui -maxdepth 1 -type d -not -path "./ui" -not -path "./ui/worker"); do \
		echo $$dir; \
		npx tsfmt -r --useTsfmt ./tsfmt.json --baseDir $$dir; \
	done

# one time setup to install pre-commit hook for gofmt and npm install needed packages
setup:
	cp pre-commit .git/hooks
	chmod +x .git/hooks/pre-commit
