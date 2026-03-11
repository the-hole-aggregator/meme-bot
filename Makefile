# Defining variables for all scripts
SCRIPTS_DIR := scripts
GIT_HOOKS_INIT := $(SCRIPTS_DIR)/git_hooks_init.sh


# Tasks to run each script
hooks_init:
	sh $(GIT_HOOKS_INIT)

	
help:
	@echo " - hooks_init: init git hooks"
