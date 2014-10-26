rm -f ~/gothings/bin/revelgen
rm -f template.go
go-bindata -o template.go template/
go install
rm -f template.go