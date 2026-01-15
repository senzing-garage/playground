ifeq ($(OS),Windows_NT)
    OSTYPE = windows
    ifeq ($(PROCESSOR_ARCHITECTURE), AMD64)
        OSARCH = x86_64
    endif
    ifeq ($(PROCESSOR_ARCHITECTURE), ARM64)
        OSARCH = arm64
    endif
    define activate-venv
        .venv\\Scripts\\activate
    endef
else
    OSTYPE = $(shell uname -s | tr '[:upper:]' '[:lower:]')
    OSARCH = $(shell uname -m | tr '[:upper:]' '[:lower:]')
    ifeq ($(OSARCH),aarch64)
        OSARCH = arm64
    endif
    define activate-venv
        @. .venv/bin/activate
    endef
endif
