CURRENT := $(abspath .)

DBUSER := "growerlab"
DBPWD := "growerlab"
DBNAME := "growerlab"


dbschema:
	echo $(CURRENT)
	mysqldump -u$(DBUSER) -p$(DBPWD) -d $(DBNAME) > db/growerlab.sql;

dbschema-docker:
	echo $(CURRENT)
	docker exec -it 6da bash -c 'mysqldump -u$(DBUSER) -p$(DBPWD) -d $(DBNAME)' > db/growerlab.sql