package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i AuthService -o mocks -s _minimock.go
//go:generate minimock -i PVZService -o mocks -s _minimock.go
//go:generate minimock -i ReceptionService -o mocks -s _minimock.go
