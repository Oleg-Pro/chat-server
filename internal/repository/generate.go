package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/minimock -i ChatRepository -o ./mocks/ -s "_minimock.go"
