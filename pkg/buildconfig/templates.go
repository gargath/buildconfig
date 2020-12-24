package buildconfig

const scriptTmpl string = `#!/usr/bin/env bash

if (( ${BASH_VERSION:0:1} < 4 )); then
  echo "This configure script requires at least Bash 4"
  if [[ "$OSTYPE" == darwin* ]]; then
    echo "On macOS, you can upgrade it using Homebrew <https://brew.sh/> by typing: brew install bash"
  fi
  exit 1
fi

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

declare -A tools=()
declare -A desired=()

vercomp () {
    if [[ $1 == $2 ]]
    then
        return 0
    fi
    local IFS=.
    local i ver1=($1) ver2=($2)
    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    return 0
}

check_for() {
  echo -n "Checking for $1... "
  if ! [ -z "${desired[$1]}" ]; then
    TOOL_PATH="${desired[$1]}"
  else
    TOOL_PATH=$(command -v $1)
  fi
  if ! [ -x "$TOOL_PATH" -a -f "$TOOL_PATH" ]; then
    printf "${RED}not found${NC}\n"
    cd - > /dev/null
    exit 1
  else
    printf "${GREEN}found${NC}\n"
    tools[$1]=$TOOL_PATH
  fi
}

{{- range .Dependencies -}}
{{ if .VersionCheck.MinVersion }}
check_{{ .Name }}_version() {
  echo -n "Checking {{ .Name }} version... "
  {{ .Name | ToUpper }}_VERSION=$(${tools[{{ .Name }}]} {{ .VersionCheck.Command }})
  vercomp ${{.Name | ToUpper }}_VERSION {{ .VersionCheck.MinVersion}}
  case $? in
    0) ;&
    1)
      printf "${GREEN}"
      echo ${{ .Name | ToUpper }}_VERSION
      printf "${NC}"
      ;;
    2)
      printf "${RED}"
      echo "${{ .Name | ToUpper }}_VERSION < {{ .VersionCheck.MinVersion }}"
      exit 1
      ;;
  esac
}
{{- end }}
{{ end }}

for arg in "$@"; do
  case ${arg%%=*} in

{{- range .Dependencies}}
    "--with-{{.Name}}")
      desired[{{.Name}}]="${arg##*=}"
      ;;
{{- end}}
    "--help")
      printf "${GREEN}$0${NC}\n"
      printf "  available options:\n"
      printf "  --with-go=${BLUE}<path_to_go_binary>${NC}\n"
      exit 0
      ;;
    *)
      echo "Unknown option: $arg"
      exit 2
      ;;
  esac
done


cd ${0%/*}

{{range .Dependencies -}}
check_for {{.Name}}
{{if .VersionCheck.MinVersion -}}
check_{{.Name}}_version
{{end -}}
{{- end}}

cat <<- EOF > .env
{{- range .Dependencies}}
{{.Name | ToUpper}} := ${tools[{{.Name}}]}
{{- end}}
EOF

echo "Environment configuration written to .env"

cd - > /dev/null`

const makefileTmpl string = `include .env

BINARY := {{ .Binary }}
VERSION := $(shell git describe --always --dirty --tags 2>/dev/null || echo "undefined")
ECHO := echo

.NOTPARALLEL:

.PHONY: all
all: test build

.PHONY: build
build: clean $(BINARY)

.PHONY: clean
clean:
	rm -f $(BINARY)

.PHONY: distclean
distclean: clean
	rm -f .env

.PHONY: fmt
fmt:
	$(GO) fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	$(GO) vet -tags dev -composites=false ./pkg/... ./cmd/...

lint:
	@ $(ECHO) "\033[36mLinting code\033[0m"
	$(LINTER) run --disable-all \
                --exclude-use-default=false \
                --enable=govet \
                --enable=ineffassign \
                --enable=deadcode \
                --enable=golint \
                --enable=goconst \
                --enable=gofmt \
                --enable=goimports \
                --skip-dirs=pkg/client/ \
                --deadline=120s \
                --tests ./...
	@ $(ECHO)

.PHONY: check
check: fmt lint vet test

.PHONY: test
test:
	@ $(ECHO) "\033[36mRunning test suite in Ginkgo\033[0m"
	$(GINKGO) -v -p -race -randomizeAllSpecs ./pkg/... ./cmd/...
	@ $(ECHO)

# Build sis
$(BINARY): fmt vet
	GO111MODULE=on CGO_ENABLED=0 $(GO) build -o $(BINARY) -ldflags="-X main.VERSION=${VERSION}" ./cmd/...
`
