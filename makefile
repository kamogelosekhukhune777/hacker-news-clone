#Variables
GOCMD=go
GORUN=$(GOCMD) run

#Targets
run : $(GORUN) cmd/web/*.go