server "local" {
    connection "localhost" {
        host = "127.0.0.1"
        port = 1883
        #username = ""
        #password = ""
        #clientID = ""
        ssl = false
    }

    dataFile = "./rocket-view.dat"
}
