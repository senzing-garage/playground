# Makefile extensions for linux-x86_64.

# -----------------------------------------------------------------------------
# OS-ARCH specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: linux/amd64
	@mkdir -p $(TARGET_DIRECTORY)/linux
	@cp $(TARGET_DIRECTORY)/linux-amd64/$(PROGRAM_NAME) $(TARGET_DIRECTORY)/linux/$(PROGRAM_NAME)


.PHONY: build-with-libsqlite3-osarch-specific
build-with-libsqlite3-osarch-specific: linux/amd64/libsqlite3
	@mkdir -p $(TARGET_DIRECTORY)/linux
	@cp $(TARGET_DIRECTORY)/linux-amd64/$(PROGRAM_NAME) $(TARGET_DIRECTORY)/linux/$(PROGRAM_NAME)


.PHONY: docker-build-osarch-specific
docker-build-osarch-specific:
	@docker build \
		--platform linux/amd64 \
		--tag $(DOCKER_IMAGE_NAME) \
		--tag $(DOCKER_IMAGE_NAME):$(BUILD_VERSION) \
		.

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-linux-x86_64
only-linux-x86_64:
	$(info Only linux-x86_64 has this Makefile target.)
